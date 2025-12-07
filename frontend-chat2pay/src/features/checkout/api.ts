import { apiClient } from "@/shared/api/client";
import type { BackendResponse } from "@/shared/api/types";
import { endpoints } from "@/shared/api/endpoints";
import type { Province, City, CourierResult, CreateOrderRequest, Order } from "./types";

export async function getProvinces() {
  const response = await apiClient.get<BackendResponse<Province[]>>(endpoints.shipping.provinces);
  if (!response.data.status) {
    throw new Error(response.data.error || "Failed to get provinces");
  }
  return response.data.data || [];
}

export async function getCities(provinceId?: string) {
  const response = await apiClient.get<BackendResponse<City[]>>(endpoints.shipping.cities(provinceId));
  if (!response.data.status) {
    throw new Error(response.data.error || "Failed to get cities");
  }
  return response.data.data || [];
}

export async function getShippingCost(origin: string, destination: string, weight: number, courier?: string) {
  const response = await apiClient.get<BackendResponse<CourierResult[]>>(
    endpoints.shipping.cost(origin, destination, weight, courier)
  );
  if (!response.data.status) {
    throw new Error(response.data.error || "Failed to get shipping cost");
  }
  return response.data.data || [];
}

export async function createOrder(token: string, data: CreateOrderRequest) {
  const response = await apiClient.post<BackendResponse<Order>>(endpoints.orders.create, {
    json: data,
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });
  if (!response.data.status) {
    throw new Error(response.data.error || "Failed to create order");
  }
  return response.data.data;
}

export async function getCustomerOrders(token: string) {
  const response = await apiClient.get<BackendResponse<{ orders: Order[]; total: number }>>(
    endpoints.orders.customer,
    {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    }
  );
  if (!response.data.status) {
    throw new Error(response.data.error || "Failed to get orders");
  }
  return response.data.data?.orders || [];
}

export async function getMerchantOrders(token: string) {
  const response = await apiClient.get<BackendResponse<{ orders: Order[]; total: number }>>(
    endpoints.orders.merchant,
    {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    }
  );
  if (!response.data.status) {
    throw new Error(response.data.error || "Failed to get orders");
  }
  return response.data.data?.orders || [];
}

export interface TrackingResult {
  delivered: boolean;
  summary: {
    courier_code: string;
    courier_name: string;
    waybill_number: string;
    service_code: string;
    waybill_date: string;
    shipper_name: string;
    receiver_name: string;
    origin: string;
    destination: string;
    status: string;
  };
  manifest: {
    manifest_description: string;
    manifest_date: string;
    manifest_time: string;
    city_name: string;
  }[];
  delivery_status: {
    status: string;
    pod_receiver: string;
    pod_date: string;
    pod_time: string;
  };
}

export async function trackShipment(waybill: string, courier: string): Promise<TrackingResult> {
  const response = await apiClient.get<BackendResponse<TrackingResult>>(
    endpoints.shipping.track(waybill, courier)
  );
  if (!response.data.status) {
    throw new Error(response.data.error || "Failed to track shipment");
  }
  return response.data.data!;
}
