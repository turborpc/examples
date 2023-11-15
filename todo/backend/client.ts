const datePrefix = "__turborpc.Date";

class RPCError extends Error {
	readonly service: string;
	readonly method: string;

	constructor(message: string, service: string, method: string) {
		super(message);

		this.name = "RPCError";
		this.service = service;
		this.method = method;
	}
}

function reviver(key: string, value: any): Date | any {
	if (typeof value !== "string" || !value.startsWith(datePrefix)) {
		return value;
	}

	const str = value.slice(datePrefix.length + 1, -1);
	const timestamp = Date.parse(str);

	if (isNaN(timestamp)) {
		return value;
	}

	return new Date(timestamp);
}

async function call(url: string, service: string, method: string, input: any, headers?: HeadersInit): Promise<unknown> {
	const res = await fetch(url + "?service=" + service + "&method=" + method, {
		method: "POST",
		headers: headers,
		body: JSON.stringify(input)
	});

	const text = await res.text();
	const data = JSON.parse(text, reviver);

	if (res.status !== 200) {
		if (typeof data.message === "string") {
			throw new RPCError(data.message, service, method);
		} else {
			throw new RPCError("unknown error", service, method);
		}
	}

	return data.output;
}




export class Todo {
	name: string;
	url: string;
	headers?: HeadersInit;

	constructor(url: string, headers?: HeadersInit) {
		this.name = "Todo";
		this.url = url;
		this.headers = headers;
	}

	async add(input: string) {
		await call(this.url, this.name, "Add", input, this.headers);
	}
	async all(): Promise<{ "Content": string; "Created": Date; "Finished": (Date | null); }[]> {
		return call(this.url, this.name, "All", null, this.headers) as Promise<{ "Content": string; "Created": Date; "Finished": (Date | null); }[]>;
	}
	async toggle(input: number): Promise<{ "Content": string; "Created": Date; "Finished": (Date | null); }> {
		return call(this.url, this.name, "Toggle", input, this.headers) as Promise<{ "Content": string; "Created": Date; "Finished": (Date | null); }>;
	}
	
}


export class RPC {
	url: string;
	headers?: HeadersInit;

	todo: Todo;

	constructor(url: string, headers?: HeadersInit) {
		this.url = url;
		this.headers = headers;

		this.todo = new Todo(url, headers);
	}
}