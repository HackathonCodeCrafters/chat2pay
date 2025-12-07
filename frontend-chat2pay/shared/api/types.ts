export type HttpMethod = "GET" | "POST" | "PUT" | "PATCH" | "DELETE";

export type QueryValue = string | number | boolean | null | undefined;
export type QueryParams = Record<string, QueryValue | QueryValue[]>;

export type ApiResponse<T> = {
  data: T;
  status: number;
  headers: Headers;
};

export type ApiRequestOptions = Omit<RequestInit, "body" | "headers" | "method"> & {
  method?: HttpMethod;
  headers?: HeadersInit;
  query?: QueryParams;
  body?: BodyInit | null;
  json?: unknown;
};

export class ApiError extends Error {
  status: number;
  payload?: unknown;

  constructor(init: { message: string; status: number; payload?: unknown }) {
    super(init.message);
    this.name = "ApiError";
    this.status = init.status;
    this.payload = init.payload;
  }
}

// Backend response wrapper
export interface BackendResponse<T> {
  status: boolean;
  data: T;
  error: string | null;
}

// Pagination
export interface PaginatedResponse<T> {
  items: T[];
  total: number;
  page: number;
  limit: number;
}
