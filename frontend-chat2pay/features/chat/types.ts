export interface Product {
  id: string;
  merchant_id: string;
  outlet_id?: string;
  category_id?: string;
  category?: {
    id: string;
    name: string;
  };
  name: string;
  description?: string;
  sku?: string;
  price: number;
  stock: number;
  status: string;
  image?: string;
  images?: {
    id: string;
    image_url: string;
    is_primary: boolean;
  }[];
  created_at: string;
  updated_at: string;
}

export interface LLMResponse {
  products: Product[] | null;
  message: string;
}

export interface ChatMessage {
  id: string;
  role: "user" | "assistant";
  content: string;
  products?: Product[];
  timestamp: Date;
}

export interface AskProductRequest {
  prompt: string;
}
