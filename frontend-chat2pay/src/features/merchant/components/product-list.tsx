"use client";

import * as React from "react";
import { Package, Plus } from "lucide-react";
import { Button, Badge } from "@/shared/components/atoms";
import { Card, EmptyState } from "@/shared/components/molecules";
import { formatCurrency } from "@/shared/lib/utils";
import { Product } from "../types";

interface ProductListProps {
  products: Product[];
  isLoading: boolean;
  onAddProduct: () => void;
}

export function ProductList({ products, isLoading, onAddProduct }: ProductListProps) {
  if (isLoading) {
    return (
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {[1, 2, 3].map((i) => (
          <Card key={i} className="p-4 animate-pulse">
            <div className="h-32 bg-muted rounded mb-3" />
            <div className="h-4 bg-muted rounded w-3/4 mb-2" />
            <div className="h-4 bg-muted rounded w-1/2" />
          </Card>
        ))}
      </div>
    );
  }

  if (products.length === 0) {
    return (
      <div className="py-12">
        <EmptyState
          icon={Package}
          title="No Products Yet"
          description="Start adding products to your store"
          action={
            <Button onClick={onAddProduct}>
              <Plus className="h-4 w-4 mr-2" />
              Add Product
            </Button>
          }
        />
      </div>
    );
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h2 className="text-lg font-semibold">Products ({products.length})</h2>
        <Button onClick={onAddProduct} size="sm">
          <Plus className="h-4 w-4 mr-2" />
          Add Product
        </Button>
      </div>

      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {products.map((product) => (
          <Card key={product.id} className="p-4">
            <div className="flex items-start justify-between mb-3">
              {product.image ? (
                <img
                  src={product.image}
                  alt={product.name}
                  className="h-12 w-12 object-cover rounded-lg"
                />
              ) : (
                <div className="flex h-12 w-12 items-center justify-center rounded-lg bg-muted">
                  <Package className="h-6 w-6 text-muted-foreground" />
                </div>
              )}
              <Badge variant={product.status === "active" ? "default" : "secondary"}>
                {product.status}
              </Badge>
            </div>
            <h3 className="font-medium line-clamp-1">{product.name}</h3>
            {product.description && (
              <p className="text-sm text-muted-foreground line-clamp-2 mt-1">
                {product.description}
              </p>
            )}
            <div className="flex items-center justify-between mt-3 pt-3 border-t">
              <span className="font-bold text-primary">
                {formatCurrency(product.price)}
              </span>
              <span className="text-sm text-muted-foreground">
                Stock: {product.stock}
              </span>
            </div>
          </Card>
        ))}
      </div>
    </div>
  );
}
