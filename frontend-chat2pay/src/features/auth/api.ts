import { apiClient, endpoints, BackendResponse } from "@/shared/api";
import {
  CustomerAuthResponse,
  LoginCredentials,
  CustomerRegisterData,
} from "./types";

export async function loginCustomer(credentials: LoginCredentials) {
  const response = await apiClient.post<BackendResponse<CustomerAuthResponse>>(
    endpoints.auth.customer.login,
    { json: credentials }
  );

  if (!response.data.status) {
    throw new Error(response.data.error || "Login failed");
  }

  return response.data.data;
}

export async function registerCustomer(data: CustomerRegisterData) {
  const response = await apiClient.post<BackendResponse<CustomerAuthResponse>>(
    endpoints.auth.customer.register,
    { json: data }
  );

  if (!response.data.status) {
    throw new Error(response.data.error || "Registration failed");
  }

  return response.data.data;
}
