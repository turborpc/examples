# Todo Example

Monorepo NextJS+TurboRPC Todo example.

## Relevant files

- [backend.go](./backend/backend.go) the TurboRPC todo service using [Chi](https://github.com/go-chi/chi).
- [index.tsx](./frontend/pages/index.tsx) using the generated client in the frontend.

## Run

```
npm i
npm start --workspace=@app/backend
npm start --workspace=@app/frontend
```
