export interface User {
  id: string;
  name: string;
  email: string;
  phone?: string;
  role: "customer" | "merchant";
}

export interface CustomerAuthResponse {
  id: string;
  name: string;
  email?: string;
  phone?: string;
  role: string;
  access_token: string;
}

export interface MerchantAuthResponse {
  id: string;
  merchant_id: string;
  merchant?: {
    id: string;
    name: string;
    email: string;
    status: string;
  };
  name: string;
  email: string;
  role: string;
  status: string;
  access_token: string;
}

export interface LoginCredentials {
  email: string;
  password: string;
}

export interface CustomerRegisterData {
  name: string;
  email: string;
  phone?: string;
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

export interface AuthState {
  user: User | null;
  token: string | null;
  isLoading: boolean;
  isAuthenticated: boolean;
}
