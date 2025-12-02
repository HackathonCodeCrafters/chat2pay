import { getRuntimeConfig } from "@/lib/config/runtime";

import {
  ApiClientOptions,
  ApiError,
  ApiRequestOptions,
  ApiResponse,
  HttpMethod,
  QueryParams,
} from "./types";

const isAbsoluteUrl = (value: string) => /^https?:\/\//i.test(value);
const normalizePath = (path: string) => {
  const trimmed = path.trim();
  if (!trimmed) return "/";
  if (isAbsoluteUrl(trimmed)) return trimmed;
  return trimmed.startsWith("/") ? trimmed : `/${trimmed}`;
};

const appendQuery = (url: string, query?: QueryParams) => {
  if (!query) return url;

  const params = new URLSearchParams();

  Object.entries(query).forEach(([key, value]) => {
    if (value === null || value === undefined) return;
    if (Array.isArray(value)) {
      value.forEach((item) => params.append(key, String(item)));
      return;
    }
    params.append(key, String(value));
  });

  const serialized = params.toString();
  if (!serialized) return url;

  const separator = url.includes("?") ? "&" : "?";
  return `${url}${separator}${serialized}`;
};

const resolveUrl = (baseUrl: string, path: string, query?: QueryParams) => {
  const normalizedPath = normalizePath(path);
  const base = baseUrl.replace(/\/+$/, "");

  const fullPath = isAbsoluteUrl(normalizedPath)
    ? normalizedPath
    : `${base}${normalizedPath}`;

  return appendQuery(fullPath, query);
};

const toErrorMessage = (fallback: string, payload?: unknown) => {
  if (payload && typeof payload === "object" && "message" in payload) {
    const message = payload.message;
    if (typeof message === "string" && message.length > 0) return message;
  }
  return fallback;
};

export class ApiClient {
  private readonly baseUrl: string;
  private readonly fetcher: typeof fetch;
  private readonly defaultHeaders?: HeadersInit;
  private readonly getAuthToken?: ApiClientOptions["getAuthToken"];

  constructor({
    baseUrl,
    defaultHeaders,
    fetcher,
    getAuthToken,
  }: ApiClientOptions) {
    const runtimeBase = getRuntimeConfig().apiBaseUrl;
    this.baseUrl = baseUrl ?? runtimeBase;
    this.defaultHeaders = defaultHeaders;
    this.fetcher = fetcher ?? fetch;
    this.getAuthToken = getAuthToken;
  }

  async request<T>(
    path: string,
    options: ApiRequestOptions = {}
  ): Promise<ApiResponse<T>> {
    const { method = "GET", headers, json, body, query, ...rest } = options;
    const url = resolveUrl(this.baseUrl, path, query);

    const mergedHeaders = new Headers(this.defaultHeaders);

    if (!mergedHeaders.has("Accept")) {
      mergedHeaders.set("Accept", "application/json");
    }

    if (headers) {
      new Headers(headers).forEach((value, key) =>
        mergedHeaders.set(key, value)
      );
    }

    let requestBody: BodyInit | null = body ?? null;
    if (json !== undefined) {
      if (!mergedHeaders.has("Content-Type")) {
        mergedHeaders.set("Content-Type", "application/json");
      }
      requestBody = JSON.stringify(json);
    }

    if (this.getAuthToken) {
      const token = await this.getAuthToken();
      if (token) mergedHeaders.set("Authorization", `Bearer ${token}`);
    }

    const response = await this.fetcher(url, {
      method,
      headers: mergedHeaders,
      body: requestBody,
      ...rest,
    });

    const responseType = response.headers.get("content-type") ?? "";
    const isJson = responseType.includes("json");
    const hasBody = ![204, 205].includes(response.status);
    let parsedBody: unknown = null;

    if (hasBody) {
      try {
        parsedBody = isJson ? await response.json() : await response.text();
      } catch {
        parsedBody = null;
      }
    }

    if (!response.ok) {
      const payload = isJson ? parsedBody : undefined;
      const message = toErrorMessage(
        `Request failed with status ${response.status}`,
        payload
      );

      throw new ApiError({
        message,
        status: response.status,
        payload,
        requestId: response.headers.get("x-request-id"),
      });
    }

    return {
      data: parsedBody as T,
      status: response.status,
      headers: response.headers,
    };
  }

  get<T>(
    path: string,
    options?: Omit<ApiRequestOptions, "method" | "body" | "json">
  ) {
    return this.request<T>(path, { ...options, method: "GET" });
  }

  post<T>(path: string, options?: Omit<ApiRequestOptions, "method">) {
    return this.request<T>(path, { ...options, method: "POST" });
  }

  put<T>(path: string, options?: Omit<ApiRequestOptions, "method">) {
    return this.request<T>(path, { ...options, method: "PUT" });
  }

  patch<T>(path: string, options?: Omit<ApiRequestOptions, "method">) {
    return this.request<T>(path, { ...options, method: "PATCH" });
  }

  delete<T>(
    path: string,
    options?: Omit<ApiRequestOptions, "method" | "json">
  ) {
    return this.request<T>(path, { ...options, method: "DELETE" });
  }
}


export const apiClient = new ApiClient({});

export const createApiClient = (options?: ApiClientOptions) =>
  new ApiClient(options ?? {});
