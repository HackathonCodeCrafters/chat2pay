export { ApiClient, apiClient, createApiClient } from "@/shared/api/client";
export { ApiError } from "@/shared/api/types";
export type { ApiResponse, ApiRequestOptions, QueryParams, BackendResponse, PaginatedResponse } from "@/shared/api/types";

// Endpoints compatible with app/ pages that use .root() pattern
export const endpoints = {
  health: () => "/health",
  merchants: {
    root: () => "/merchants",
    byId: (merchantId: string) => `/merchants/${merchantId}`,
  },
  products: {
    root: () => "/products",
    byId: (productId: string) => `/products/${productId}`,
  },
  payments: {
    root: () => "/payments",
    byId: (paymentId: string) => `/payments/${paymentId}`,
    capture: (paymentId: string) => `/payments/${paymentId}/capture`,
  },
  customers: {
    root: () => "/customers",
    byId: (customerId: string) => `/customers/${customerId}`,
  },
  orders: {
    root: () => "/orders",
    byId: (orderId: string) => `/orders/${orderId}`,
  },
};
