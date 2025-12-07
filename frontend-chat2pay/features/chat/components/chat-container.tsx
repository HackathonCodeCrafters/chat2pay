"use client";

import * as React from "react";
import { useRouter } from "next/navigation";
import { MessageSquare, Sparkles, Trash2, ShoppingCart, Send, Zap, Bot, Cpu } from "lucide-react";
import { Button, Spinner } from "@/shared/components/atoms";
import { useAuth } from "@/features/auth";
import { useCart } from "@/features/cart";
import { useChat } from "../hooks/use-chat";
import { ChatInput } from "./chat-input";
import { ChatMessage } from "./chat-message";
import type { Product } from "../types";

const SUGGESTIONS = [
  { text: "Laptop gaming budget 15 juta", icon: "💻" },
  { text: "Headphone wireless untuk kerja", icon: "🎧" },
  { text: "Mouse ergonomis untuk coding", icon: "🖱️" },
  { text: "Keyboard mechanical murah", icon: "⌨️" },
];

export function ChatContainer() {
  const router = useRouter();
  const { user } = useAuth();
  const { messages, isLoading, isConnected, sendMessage, clearMessages } = useChat();
  const { addItem, totalItems } = useCart();
  const messagesEndRef = React.useRef<HTMLDivElement>(null);
  const [buyAlert, setBuyAlert] = React.useState<string | null>(null);
  const [showClearConfirm, setShowClearConfirm] = React.useState(false);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  };

  React.useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const handleBuyProduct = (product: Product) => {
    addItem({
      product: {
        id: product.id,
        name: product.name,
        price: product.price,
        image: product.image,
        weight: 1000,
        merchant_id: product.merchant_id,
      },
      quantity: 1,
    });
    setBuyAlert(`"${product.name}" ditambahkan ke keranjang!`);
    setTimeout(() => setBuyAlert(null), 3000);
  };

  const handleClearChat = () => {
    clearMessages();
    setShowClearConfirm(false);
  };

  return (
    <div className="flex h-[calc(100vh-4rem)] flex-col relative">
      {/* Animated background elements */}
      <div className="absolute inset-0 overflow-hidden pointer-events-none">
        <div className="absolute top-20 left-10 w-72 h-72 bg-violet-500/10 rounded-full blur-3xl animate-pulse" />
        <div className="absolute bottom-20 right-10 w-96 h-96 bg-cyan-500/10 rounded-full blur-3xl animate-pulse" style={{ animationDelay: '1s' }} />
        <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[600px] h-[600px] bg-violet-500/5 rounded-full blur-3xl" />
      </div>

      {/* Confirmation Dialog */}
      {showClearConfirm && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm">
          <div className="glass-card rounded-2xl p-6 max-w-sm mx-4 animate-in fade-in zoom-in duration-200">
            <div className="flex items-center gap-3 mb-4">
              <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-red-500/20">
                <Trash2 className="h-5 w-5 text-red-400" />
              </div>
              <h3 className="font-semibold text-lg">Hapus Riwayat Chat?</h3>
            </div>
            <p className="text-sm text-muted-foreground mb-6">
              Semua percakapan akan dihapus permanen dan tidak dapat dikembalikan.
            </p>
            <div className="flex gap-3 justify-end">
              <button 
                onClick={() => setShowClearConfirm(false)}
                className="px-4 py-2 rounded-lg text-sm text-muted-foreground hover:text-foreground hover:bg-white/5 transition-all"
              >
                Batal
              </button>
              <button 
                onClick={handleClearChat}
                className="px-4 py-2 rounded-lg text-sm font-medium bg-red-500/20 text-red-400 hover:bg-red-500/30 transition-all flex items-center gap-2"
              >
                <Trash2 className="h-4 w-4" />
                Hapus
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Success Alert */}
      {buyAlert && (
        <div className="fixed top-20 left-1/2 -translate-x-1/2 z-50 animate-in slide-in-from-top fade-in duration-300">
          <div className="glass-card rounded-xl px-4 py-3 flex items-center gap-3 border-emerald-500/30">
            <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-emerald-500/20">
              <ShoppingCart className="h-4 w-4 text-emerald-400" />
            </div>
            <span className="text-sm text-emerald-400">{buyAlert}</span>
          </div>
        </div>
      )}

      {/* Header */}
      <div className="relative z-10 glass border-b border-white/5">
        <div className="container py-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-4">
              <div className="relative">
                <div className="flex h-12 w-12 items-center justify-center rounded-2xl bg-gradient-to-br from-violet-500 to-cyan-500 shadow-lg shadow-violet-500/25">
                  <Bot className="h-6 w-6 text-white" />
                </div>
                <div className="absolute -bottom-1 -right-1 flex h-5 w-5 items-center justify-center rounded-full bg-emerald-500 border-2 border-background">
                  <Zap className="h-3 w-3 text-white" />
                </div>
              </div>
              <div>
                <h1 className="font-semibold text-lg flex items-center gap-2">
                  AI Shopping Assistant
                  <span className="px-2 py-0.5 rounded-full text-[10px] font-medium bg-gradient-to-r from-violet-500/20 to-cyan-500/20 text-cyan-400 border border-cyan-500/30">
                    GPT-4
                  </span>
                </h1>
                <p className="text-xs text-muted-foreground flex items-center gap-1">
                  <span className={`inline-block h-1.5 w-1.5 rounded-full ${isConnected ? 'bg-emerald-500 animate-pulse' : 'bg-yellow-500'}`} />
                  {isConnected ? 'Online - Siap membantu Anda' : 'Menghubungkan...'}
                </p>
              </div>
            </div>
            <div className="flex items-center gap-2">
              <button 
                onClick={() => router.push("/cart")} 
                className="relative flex items-center gap-2 px-3 py-2 rounded-xl text-sm text-muted-foreground hover:text-foreground glass-card-hover"
              >
                <ShoppingCart className="h-4 w-4" />
                <span className="hidden sm:inline">Keranjang</span>
                {totalItems > 0 && (
                  <span className="absolute -top-1 -right-1 h-5 w-5 rounded-full bg-gradient-to-r from-violet-500 to-cyan-500 text-[10px] font-bold text-white flex items-center justify-center shadow-lg">
                    {totalItems}
                  </span>
                )}
              </button>
              {messages.length > 0 && (
                <button 
                  onClick={() => setShowClearConfirm(true)}
                  className="flex items-center gap-2 px-3 py-2 rounded-xl text-sm text-muted-foreground hover:text-red-400 hover:bg-red-500/10 transition-all"
                >
                  <Trash2 className="h-4 w-4" />
                  <span className="hidden sm:inline">Hapus</span>
                </button>
              )}
            </div>
          </div>
        </div>
      </div>

      {/* Messages */}
      <div className="flex-1 overflow-y-auto relative z-10">
        {messages.length === 0 ? (
          <div className="flex h-full flex-col items-center justify-center p-6">
            <div className="text-center max-w-lg">
              {/* Hero Icon */}
              <div className="relative mx-auto mb-8">
                <div className="flex h-24 w-24 items-center justify-center rounded-3xl bg-gradient-to-br from-violet-500/20 to-cyan-500/20 border border-white/10 mx-auto">
                  <Cpu className="h-12 w-12 text-violet-400" />
                </div>
                <div className="absolute -top-2 -right-2 flex h-8 w-8 items-center justify-center rounded-xl bg-gradient-to-br from-violet-500 to-cyan-500 shadow-lg shadow-violet-500/30 animate-bounce">
                  <Sparkles className="h-4 w-4 text-white" />
                </div>
              </div>

              <h2 className="text-2xl font-bold mb-3 gradient-text">
                Hai, {user?.name || 'Selamat Datang'}! 👋
              </h2>
              <p className="text-muted-foreground mb-8">
                Saya adalah AI Shopping Assistant. Ceritakan produk apa yang Anda cari, 
                dan saya akan membantu menemukan yang terbaik untuk Anda.
              </p>

              {/* Suggestion Cards */}
              <div className="grid grid-cols-2 gap-3">
                {SUGGESTIONS.map((suggestion) => (
                  <button
                    key={suggestion.text}
                    onClick={() => sendMessage(suggestion.text)}
                    className="glass-card glass-card-hover rounded-xl p-4 text-left group"
                  >
                    <span className="text-2xl mb-2 block">{suggestion.icon}</span>
                    <span className="text-sm text-muted-foreground group-hover:text-foreground transition-colors">
                      {suggestion.text}
                    </span>
                  </button>
                ))}
              </div>
            </div>
          </div>
        ) : (
          <div className="container py-6">
            <div className="max-w-3xl mx-auto space-y-6">
              {messages.map((message) => (
                <ChatMessage
                  key={message.id}
                  message={message}
                  userName={user?.name}
                  onBuyProduct={handleBuyProduct}
                />
              ))}
              {isLoading && (
                <div className="flex items-start gap-4">
                  <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-gradient-to-br from-violet-500/20 to-cyan-500/20 border border-white/10 shrink-0">
                    <Bot className="h-5 w-5 text-violet-400 animate-pulse" />
                  </div>
                  <div className="glass-card rounded-2xl rounded-tl-md px-4 py-3 flex items-center gap-3">
                    <div className="flex gap-1">
                      <span className="h-2 w-2 rounded-full bg-violet-500 animate-bounce" style={{ animationDelay: '0ms' }} />
                      <span className="h-2 w-2 rounded-full bg-cyan-500 animate-bounce" style={{ animationDelay: '150ms' }} />
                      <span className="h-2 w-2 rounded-full bg-emerald-500 animate-bounce" style={{ animationDelay: '300ms' }} />
                    </div>
                    <span className="text-sm text-muted-foreground">
                      Sedang mencari produk terbaik untuk Anda...
                    </span>
                  </div>
                </div>
              )}
              <div ref={messagesEndRef} />
            </div>
          </div>
        )}
      </div>

      {/* Input */}
      <div className="relative z-10 glass border-t border-white/5">
        <div className="container py-4">
          <div className="max-w-3xl mx-auto">
            <ChatInput
              onSend={sendMessage}
              disabled={isLoading}
              placeholder="Ketik pesan atau tanyakan produk yang Anda cari..."
            />
            <p className="text-[10px] text-center text-muted-foreground mt-2">
              Chat2Pay AI dapat membuat kesalahan. Periksa informasi penting.
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}
