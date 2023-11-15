package main

import (
	"context"
	"errors"
	"net/http"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/turborpc/turborpc"
)

var ErrInvalidIndex = errors.New("invalid index")

type TodoItem struct {
	Content  string
	Created  turborpc.Date
	Finished *turborpc.Date
}

type Todo struct {
	mu    sync.Mutex
	todos []TodoItem
}

func (t *Todo) All(ctx context.Context) (turborpc.NonNullSlice[TodoItem], error) {
	return t.todos, nil
}

func (t *Todo) Add(ctx context.Context, content string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.todos = append([]TodoItem{{
		Content: content,
		Created: turborpc.Date(time.Now()),
	}}, t.todos...)

	return nil
}

func (t *Todo) Toggle(ctx context.Context, index int) (TodoItem, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if index < 0 || index >= len(t.todos) {
		return TodoItem{}, ErrInvalidIndex
	}

	var finished turborpc.Date
	if t.todos[index].Finished == nil {
		finished = turborpc.Date(time.Now())
		t.todos[index].Finished = &finished
	} else {
		t.todos[index].Finished = nil
	}

	return t.todos[index], nil
}

// dirname returns the directory name of the file that calls it.
func dirname() string {
	_, file, _, _ := runtime.Caller(1)

	return filepath.Dir(file)
}

func main() {
	rpc := turborpc.NewServer()

	_ = rpc.Register(&Todo{})

	_ = rpc.WriteTypeScriptClient(filepath.Join(dirname(), "client.ts"))

	r := chi.NewRouter()
	r.Use(middleware.Logger, cors.AllowAll().Handler)

	r.Handle("/rpc", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "turborpc" {
			turborpc.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		rpc.ServeHTTP(w, r)
	}))

	http.ListenAndServe(":5000", r)
}
