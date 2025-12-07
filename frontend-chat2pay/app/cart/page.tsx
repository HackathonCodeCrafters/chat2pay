"use client";

import * as React from "react";
import { useRouter } from "next/navigation";
import { ShoppingCart, Trash2, Plus, Minus, ArrowLeft } from "lucide-react";
import { Button, Spinner } from "@/shared/components/atoms";
import { EmptyState } from "@/shared/components/molecules";
import { useCart } from "@/features/cart";
import { formatCurrency } from "@/shared/lib/utils";

export default function CartPage() {
  const router = useRouter();
  const { items, isLoaded, totalItems, totalPrice, removeItem, updateQuantity, clearCart } = useCart();

  if (!isLoaded) {
    return (
      <div className="container py-8 flex justify-center">
        <Spinner size="lg" />
      </div>
    );
  }

  if (items.length === 0) {
    return (
      <div className="container py-8">
        <Button variant="ghost" onClick={() => router.push("/chat")} className="mb-6">
          <ArrowLeft className="h-4 w-4 mr-2" />
          Back to Chat
        </Button>
        <EmptyState
          icon={ShoppingCart}
          title="Keranjang Kosong"
          description="Belum ada produk di keranjang. Cari produk di chat untuk menambahkan ke keranjang."
        />
      </div>
    );
  }

  const handleCheckout = () => {
    router.push("/checkout");
  };

  return (
    <div className="container py-8 max-w-2xl">
      <Button variant="ghost" onClick={() => router.push("/chat")} className="mb-6">
        <ArrowLeft className="h-4 w-4 mr-2" />
        Back to Chat
      </Button>

      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold">Keranjang ({totalItems})</h1>
        <Button variant="ghost" size="sm" onClick={clearCart}>
          <Trash2 className="h-4 w-4 mr-1" />
          Kosongkan
        </Button>
      </div>

      <div className="space-y-4">
        {items.map((item) => (
          <div key={item.product.id} className="border rounded-lg p-4">
            <div className="flex gap-4">
              {item.product.image && (
                <img
                  src={item.product.image}
                  alt={item.product.name}
                  className="w-20 h-20 object-cover rounded"
                />
              )}
              <div className="flex-1">
                <h3 className="font-medium">{item.product.name}</h3>
                <p className="text-primary font-semibold">
                  {formatCurrency(item.product.price)}
                </p>
                <div className="flex items-center gap-2 mt-2">
                  <Button
                    variant="outline"
                    size="icon-sm"
                    onClick={() => updateQuantity(item.product.id, item.quantity - 1)}
                  >
                    <Minus className="h-3 w-3" />
                  </Button>
                  <span className="w-8 text-center">{item.quantity}</span>
                  <Button
                    variant="outline"
                    size="icon-sm"
                    onClick={() => updateQuantity(item.product.id, item.quantity + 1)}
                  >
                    <Plus className="h-3 w-3" />
                  </Button>
                  <Button
                    variant="ghost"
                    size="sm"
                    className="ml-auto text-destructive"
                    onClick={() => removeItem(item.product.id)}
                  >
                    <Trash2 className="h-4 w-4" />
                  </Button>
                </div>
              </div>
              <div className="text-right">
                <p className="font-semibold">
                  {formatCurrency(item.product.price * item.quantity)}
                </p>
              </div>
            </div>
          </div>
        ))}
      </div>

      {/* Summary */}
      <div className="mt-6 border rounded-lg p-4">
        <div className="flex justify-between mb-4">
          <span>Total ({totalItems} item)</span>
          <span className="font-bold text-lg text-primary">{formatCurrency(totalPrice)}</span>
        </div>
        <Button size="lg" className="w-full" onClick={handleCheckout}>
          Checkout
        </Button>
      </div>
    </div>
  );
}
