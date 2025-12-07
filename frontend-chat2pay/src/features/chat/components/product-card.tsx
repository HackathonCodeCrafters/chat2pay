"use client";

import * as React from "react";
import { Package, ShoppingCart, Plus, Star } from "lucide-react";
import { formatCurrency } from "@/shared/lib/utils";
import { Product } from "../types";

interface ProductCardProps {
  product: Product;
  onAddToCart?: (product: Product) => void;
}

export function ProductCard({ product, onAddToCart }: ProductCardProps) {
  const productImage = product.image || product.images?.find((img) => img.is_primary)?.image_url;
  const inStock = product.stock > 0;

  return (
    <div className="glass-card glass-card-hover rounded-xl overflow-hidden group">
      {/* Image */}
      <div className="relative aspect-[4/3] bg-gradient-to-br from-white/5 to-white/0 overflow-hidden">
        {productImage ? (
          <img
            src={productImage}
            alt={product.name}
            className="h-full w-full object-cover transition-transform duration-500 group-hover:scale-110"
          />
        ) : (
          <div className="flex h-full w-full items-center justify-center">
            <Package className="h-12 w-12 text-muted-foreground/30" />
          </div>
        )}
        
        {/* Overlay gradient */}
        <div className="absolute inset-0 bg-gradient-to-t from-black/60 via-transparent to-transparent" />
        
        {/* Stock badge */}
        {inStock ? (
          <div className="absolute top-2 left-2">
            <span className="px-2 py-1 rounded-lg text-[10px] font-medium bg-emerald-500/20 text-emerald-400 border border-emerald-500/30 backdrop-blur-sm">
              Tersedia
            </span>
          </div>
        ) : (
          <div className="absolute inset-0 flex items-center justify-center bg-black/60 backdrop-blur-sm">
            <span className="px-3 py-1.5 rounded-lg text-xs font-medium bg-red-500/20 text-red-400 border border-red-500/30">
              Stok Habis
            </span>
          </div>
        )}

        {/* Price on image */}
        <div className="absolute bottom-2 left-2 right-2">
          <span className="text-lg font-bold text-white drop-shadow-lg">
            {formatCurrency(product.price)}
          </span>
        </div>
      </div>

      {/* Content */}
      <div className="p-3">
        <h4 className="font-medium text-sm line-clamp-2 mb-1 group-hover:text-violet-400 transition-colors">
          {product.name}
        </h4>
        {product.description && (
          <p className="text-xs text-muted-foreground line-clamp-2 mb-3">
            {product.description}
          </p>
        )}
        
        <div className="flex items-center justify-between gap-2">
          {inStock && (
            <span className="text-[10px] text-muted-foreground">
              Stok: {product.stock}
            </span>
          )}
          {inStock && (
            <button
              onClick={() => onAddToCart?.(product)}
              className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-medium bg-gradient-to-r from-violet-500 to-cyan-500 text-white hover:shadow-lg hover:shadow-violet-500/25 transition-all duration-300"
            >
              <Plus className="h-3 w-3" />
              Keranjang
            </button>
          )}
        </div>
      </div>
    </div>
  );
}
