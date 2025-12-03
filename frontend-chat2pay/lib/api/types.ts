export type HttpMethod = "GET" | "POST" | "PUT" | "PATCH" | "DELETE";

export type QueryValue = string | number | boolean | null | undefined;
export type QueryParams = Record<string, QueryValue | QueryValue[]>;

export type ApiResponse<T> = {
  data: T;
  status: number;
  headers: Headers;
};

export type ApiClientOptions = {
  baseUrl?: string;
  defaultHeaders?: HeadersInit;
  fetcher?: typeof fetch;

  getAuthToken?: () => string | null | Promise<string | null>;
};

export type ApiRequestOptions = Omit<
  RequestInit,
  "body" | "headers" | "method"
> & {
  method?: HttpMethod;
  headers?: HeadersInit;
  query?: QueryParams;
  body?: BodyInit | null;
 
  json?: unknown;
};

export type ApiErrorPayload = {
  message?: string;
  code?: string;
  details?: unknown;
};

export class ApiError extends Error {
  status: number;
  payload?: ApiErrorPayload | unknown;
  requestId?: string | null;

  constructor(init: {
    message: string;
    status: number;
    payload?: ApiErrorPayload | unknown;
    requestId?: string | null;
  }) {
    super(init.message);
    this.name = "ApiError";
    this.status = init.status;
    this.payload = init.payload;
    this.requestId = init.requestId ?? null;
  }
}
