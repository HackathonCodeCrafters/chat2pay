import { apiClient, endpoints, BackendResponse } from "@/shared/api";
import {
  MerchantAuthResponse,
  MerchantLoginCredentials,
  MerchantRegisterData,
  Product,
  CreateProductData,
} from "./types";

export async function loginMerchant(credentials: MerchantLoginCredentials) {
  const response = await apiClient.post<BackendResponse<MerchantAuthResponse>>(
    endpoints.auth.merchant.login,
    { json: credentials }
  );

  if (!response.data.status) {
    throw new Error(response.data.error || "Login failed");
  }

  return response.data.data;
}

export async function registerMerchant(data: MerchantRegisterData) {
  const response = await apiClient.post<BackendResponse<MerchantAuthResponse>>(
    endpoints.auth.merchant.register,
    { json: data }
  );

  if (!response.data.status) {
    throw new Error(response.data.error || "Registration failed");
  }

  return response.data.data;
}

interface ProductListResponse {
  products: Product[];
  total: number;
  page: number;
  limit: number;
}

export async function getProducts(token: string, merchantId: string) {
  const response = await apiClient.get<BackendResponse<ProductListResponse>>(
    endpoints.products.listByMerchant(merchantId),
    {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    }
  );

  if (!response.data.status) {
    throw new Error(response.data.error || "Failed to get products");
  }

  return response.data.data?.products || [];
}

export async function createProduct(token: string, merchantId: string, data: CreateProductData) {
  const response = await apiClient.post<BackendResponse<Product>>(
    endpoints.products.create,
    {
      json: { ...data, merchant_id: merchantId },
      headers: {
        Authorization: `Bearer ${token}`,
      },
    }
  );

  if (!response.data.status) {
    throw new Error(response.data.error || "Failed to create product");
  }

  return response.data.data;
}
