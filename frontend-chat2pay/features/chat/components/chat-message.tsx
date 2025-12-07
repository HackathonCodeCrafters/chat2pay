"use client";

import * as React from "react";
import { Bot, User, Sparkles } from "lucide-react";
import { Avatar } from "@/shared/components/atoms";
import { cn, formatCurrency } from "@/shared/lib/utils";
import { ChatMessage as ChatMessageType, Product } from "../types";
import { ProductCard } from "./product-card";

interface ChatMessageProps {
  message: ChatMessageType;
  userName?: string;
  onBuyProduct?: (product: Product) => void;
}

export function ChatMessage({ message, userName, onBuyProduct }: ChatMessageProps) {
  const isUser = message.role === "user";

  return (
    <div
      className={cn(
        "flex gap-4",
        isUser ? "flex-row-reverse" : "flex-row"
      )}
    >
      {/* Avatar */}
      <div className="shrink-0">
        {isUser ? (
          <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-gradient-to-br from-violet-500 to-purple-600 shadow-lg shadow-violet-500/20">
            <User className="h-5 w-5 text-white" />
          </div>
        ) : (
          <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-gradient-to-br from-violet-500/20 to-cyan-500/20 border border-white/10">
            <Bot className="h-5 w-5 text-violet-400" />
          </div>
        )}
      </div>

      {/* Content */}
      <div
        className={cn(
          "flex flex-col gap-2 max-w-[85%]",
          isUser ? "items-end" : "items-start"
        )}
      >
        {/* Name */}
        <span className="text-xs text-muted-foreground px-1">
          {isUser ? (userName || "Anda") : "AI Assistant"}
        </span>

        <div
          className={cn(
            "rounded-2xl px-4 py-3",
            isUser
              ? "chat-bubble-user rounded-tr-md"
              : "chat-bubble-assistant rounded-tl-md"
          )}
        >
          <p className="text-sm whitespace-pre-wrap leading-relaxed">{message.content}</p>
        </div>

        {/* Products */}
        {message.products && message.products.length > 0 && (
          <div className="w-full mt-3">
            <div className="flex items-center gap-2 mb-3">
              <Sparkles className="h-4 w-4 text-violet-400" />
              <span className="text-xs text-muted-foreground">Produk yang ditemukan</span>
            </div>
            <div className="grid gap-3 sm:grid-cols-2">
              {message.products.map((product) => (
                <ProductCard 
                  key={product.id} 
                  product={product} 
                  onAddToCart={onBuyProduct}
                />
              ))}
            </div>
          </div>
        )}

        {/* Timestamp */}
        <span className="text-[10px] text-muted-foreground/60 px-1">
          {message.timestamp.toLocaleTimeString("id-ID", {
            hour: "2-digit",
            minute: "2-digit",
          })}
        </span>
      </div>
    </div>
  );
}
