export interface Merchant {
  id: string;
  name: string;
  legal_name?: string;
  email: string;
  phone?: string;
  status: string;
}

export interface MerchantUser {
  id: string;
  merchant_id: string;
  merchant?: Merchant;
  name: string;
  email: string;
  role: string;
  status: string;
}

export interface MerchantAuthResponse {
  id: string;
  merchant_id: string;
  merchant?: Merchant;
  name: string;
  email: string;
  role: string;
  status: string;
  access_token: string;
}

export interface MerchantLoginCredentials {
  email: string;
  password: string;
}

export interface MerchantRegisterData {
  merchant_name: string;
  legal_name?: string;
  email: string;
  phone?: string;
  name: string;
  password: string;
}

export interface Product {
  id: string;
  merchant_id: string;
  name: string;
  description?: string;
  sku?: string;
  price: number;
  stock: number;
  status: string;
  image?: string;
  created_at: string;
  updated_at: string;
}

export interface CreateProductData {
  name: string;
  description?: string;
  sku?: string;
  price: number;
  stock: number;
  image?: string;
  weight?: number;
  length?: number;
  width?: number;
  height?: number;
}
