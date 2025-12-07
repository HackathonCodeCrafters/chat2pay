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

export interface Cart {
  items: CartItem[];
  updatedAt: string;
}
