"use client";

import * as React from "react";
import Link from "next/link";
import Image from "next/image";
import { LogOut, MessageSquare, ShoppingBag, ShoppingCart, Sparkles } from "lucide-react";
import { Button, Avatar } from "@/shared/components/atoms";
import { useCart } from "@/features/cart";

interface NavbarProps {
  user?: {
    name: string;
    email?: string;
    avatar?: string;
  } | null;
  onLogout?: () => void;
}

export function Navbar({ user, onLogout }: NavbarProps) {
  const { totalItems } = useCart();

  return (
    <header className="sticky top-0 z-50 w-full glass border-b border-white/5">
      <div className="container flex h-16 items-center justify-between">
        <Link href="/" className="flex items-center">
          <div className="bg-white rounded-xl px-3 py-1.5">
            <Image 
              src="/logo-chat2pay.png" 
              alt="Chat2Pay" 
              width={140} 
              height={45}
              className="h-10 w-auto"
            />
          </div>
        </Link>

        <nav className="flex items-center gap-2">
          {user ? (
            <>
              {/* Orders link */}
              <Link 
                href="/orders" 
                className="flex items-center gap-2 px-3 py-2 rounded-lg text-sm text-muted-foreground hover:text-foreground hover:bg-white/5 transition-all duration-200"
              >
                <ShoppingBag className="h-4 w-4" />
                <span className="hidden sm:inline">Pesanan</span>
              </Link>

              {/* Cart link */}
              <Link 
                href="/cart" 
                className="relative flex items-center gap-2 px-3 py-2 rounded-lg text-sm text-muted-foreground hover:text-foreground hover:bg-white/5 transition-all duration-200"
              >
                <ShoppingCart className="h-4 w-4" />
                <span className="hidden sm:inline">Keranjang</span>
                {totalItems > 0 && (
                  <span className="absolute -top-1 -right-1 h-5 w-5 rounded-full bg-gradient-to-r from-violet-500 to-cyan-500 text-[10px] font-bold text-white flex items-center justify-center shadow-lg shadow-violet-500/30">
                    {totalItems > 99 ? "99+" : totalItems}
                  </span>
                )}
              </Link>

              {/* Divider */}
              <div className="h-8 w-px bg-white/10 mx-2 hidden sm:block" />

              {/* User info */}
              <div className="hidden items-center gap-3 sm:flex">
                <div className="relative">
                  <Avatar
                    src={user.avatar}
                    alt={user.name}
                    size="sm"
                  />
                  <div className="absolute bottom-0 right-0 h-2.5 w-2.5 rounded-full bg-emerald-500 border-2 border-background" />
                </div>
                <div className="flex flex-col">
                  <span className="text-sm font-medium">{user.name}</span>
                  {user.email && (
                    <span className="text-xs text-muted-foreground">{user.email}</span>
                  )}
                </div>
              </div>
              
              <button 
                onClick={onLogout}
                className="flex items-center justify-center h-9 w-9 rounded-lg text-muted-foreground hover:text-red-400 hover:bg-red-500/10 transition-all duration-200"
              >
                <LogOut className="h-4 w-4" />
              </button>
            </>
          ) : (
            <>
              <Link 
                href="/login"
                className="px-4 py-2 rounded-lg text-sm text-muted-foreground hover:text-foreground hover:bg-white/5 transition-all duration-200"
              >
                Login
              </Link>
              <Link 
                href="/register"
                className="px-4 py-2 rounded-lg text-sm font-medium bg-gradient-to-r from-violet-500 to-cyan-500 text-white hover:shadow-lg hover:shadow-violet-500/25 transition-all duration-300"
              >
                Register
              </Link>
            </>
          )}
        </nav>
      </div>
    </header>
  );
}
