"use client";

import * as React from "react";
import type { CartItem } from "./types";

const CART_STORAGE_KEY = "chat2pay_cart";

function loadCart(): CartItem[] {
  if (typeof window === "undefined") return [];
  try {
    const saved = localStorage.getItem(CART_STORAGE_KEY);
    if (saved) {
      return JSON.parse(saved);
    }
  } catch {}
  return [];
}

function saveCart(items: CartItem[]) {
  if (typeof window === "undefined") return;
  try {
    localStorage.setItem(CART_STORAGE_KEY, JSON.stringify(items));
  } catch {}
}

export function useCart() {
  const [items, setItems] = React.useState<CartItem[]>([]);
  const [isLoaded, setIsLoaded] = React.useState(false);

  // Load cart from localStorage on mount
  React.useEffect(() => {
    const saved = loadCart();
    setItems(saved);
    setIsLoaded(true);
  }, []);

  // Save cart to localStorage when changed
  React.useEffect(() => {
    if (isLoaded) {
      saveCart(items);
    }
  }, [items, isLoaded]);

  const addItem = React.useCallback((item: CartItem) => {
    setItems((prev) => {
      const existing = prev.find((i) => i.product.id === item.product.id);
      if (existing) {
        return prev.map((i) =>
          i.product.id === item.product.id
            ? { ...i, quantity: i.quantity + item.quantity }
            : i
        );
      }
      return [...prev, item];
    });
  }, []);

  const removeItem = React.useCallback((productId: string) => {
    setItems((prev) => prev.filter((i) => i.product.id !== productId));
  }, []);

  const updateQuantity = React.useCallback((productId: string, quantity: number) => {
    if (quantity <= 0) {
      removeItem(productId);
      return;
    }
    setItems((prev) =>
      prev.map((i) =>
        i.product.id === productId ? { ...i, quantity } : i
      )
    );
  }, [removeItem]);

  const clearCart = React.useCallback(() => {
    setItems([]);
    if (typeof window !== "undefined") {
      localStorage.removeItem(CART_STORAGE_KEY);
    }
  }, []);

  const totalItems = items.reduce((sum, i) => sum + i.quantity, 0);
  const totalPrice = items.reduce((sum, i) => sum + i.product.price * i.quantity, 0);
  const totalWeight = items.reduce((sum, i) => sum + (i.product.weight || 1000) * i.quantity, 0);

  return {
    items,
    isLoaded,
    totalItems,
    totalPrice,
    totalWeight,
    addItem,
    removeItem,
    updateQuantity,
    clearCart,
  };
}
