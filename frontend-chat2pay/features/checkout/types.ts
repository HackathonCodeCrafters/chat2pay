export interface OrderItem {
  product_id: string;
  quantity: number;
}

export interface CreateOrderRequest {
  items: OrderItem[];
  shipping_address: string;
  shipping_city: string;
  shipping_city_id: string;
  shipping_province: string;
  shipping_postal_code: string;
  courier: string;
  courier_service: string;
  shipping_cost: number;
  shipping_etd: string;
  notes?: string;
}

export interface OrderItemResponse {
  id: string;
  product_id: string;
  product_name: string;
  product_price: number;
  quantity: number;
  subtotal: number;
}

export interface Order {
  id: string;
  customer_id: string;
  merchant_id: string;
  status: string;
  subtotal: number;
  shipping_cost: number;
  total: number;
  courier?: string;
  courier_service?: string;
  shipping_etd?: string;
  tracking_number?: string;
  shipping_address?: string;
  shipping_city?: string;
  shipping_province?: string;
  shipping_postal_code?: string;
  payment_method?: string;
  payment_status: string;
  payment_url?: string;
  paid_at?: string;
  notes?: string;
  items?: OrderItemResponse[];
  created_at: string;
  updated_at: string;
}

export interface Province {
  province_id: string;
  province: string;
}

export interface City {
  city_id: string;
  province_id: string;
  province: string;
  type: string;
  city_name: string;
  postal_code: string;
}

export interface ShippingCost {
  service: string;
  description: string;
  cost: {
    value: number;
    etd: string;
    note: string;
  }[];
}

export interface CourierResult {
  code: string;
  name: string;
  costs: ShippingCost[];
}

export interface CartItem {
  product: {
    id: string;
    name: string;
    price: number;
    image?: string;
    weight: number;
    merchant_id: string;
  };
  quantity: number;
}
