"use client";

import Link from "next/link";
import Image from "next/image";
import { Sparkles, ArrowLeft } from "lucide-react";
import { LoginForm } from "@/features/auth/components";

export default function LoginPage() {
  return (
    <div className="min-h-screen flex relative overflow-hidden">
      {/* Background effects */}
      <div className="absolute inset-0 overflow-hidden pointer-events-none">
        <div className="absolute top-1/4 left-1/4 w-96 h-96 bg-violet-500/20 rounded-full blur-3xl animate-pulse" />
        <div className="absolute bottom-1/4 right-1/4 w-96 h-96 bg-cyan-500/20 rounded-full blur-3xl animate-pulse" style={{ animationDelay: '1s' }} />
        <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[800px] h-[800px] bg-violet-500/5 rounded-full blur-3xl" />
      </div>

      {/* Left side - Branding */}
      <div className="hidden lg:flex flex-1 flex-col justify-between p-12 relative">
        <Link href="/" className="flex items-center gap-3 group w-fit">
          <ArrowLeft className="h-4 w-4 text-muted-foreground group-hover:text-foreground transition-colors" />
          <span className="text-sm text-muted-foreground group-hover:text-foreground transition-colors">Kembali</span>
        </Link>

        <div>
          <div className="mb-8">
            <div className="bg-white rounded-2xl px-6 py-4 inline-block">
              <Image 
                src="/logo-chat2pay.png" 
                alt="Chat2Pay" 
                width={200} 
                height={65}
                className="h-16 w-auto"
              />
            </div>
          </div>
          <p className="text-xl text-muted-foreground max-w-md leading-relaxed">
            Temukan produk impian Anda dengan bantuan AI. Cukup ceritakan apa yang Anda cari, 
            dan kami akan membantu menemukannya.
          </p>
          
          <div className="mt-12 flex items-center gap-6">
            <div className="glass-card rounded-xl p-4">
              <p className="text-3xl font-bold gradient-text">10K+</p>
              <p className="text-xs text-muted-foreground">Pengguna Aktif</p>
            </div>
            <div className="glass-card rounded-xl p-4">
              <p className="text-3xl font-bold gradient-text">50K+</p>
              <p className="text-xs text-muted-foreground">Produk</p>
            </div>
            <div className="glass-card rounded-xl p-4">
              <p className="text-3xl font-bold gradient-text">99%</p>
              <p className="text-xs text-muted-foreground">Kepuasan</p>
            </div>
          </div>
        </div>

        <p className="text-sm text-muted-foreground">
          &copy; 2024 Chat2Pay by Code Crafters
        </p>
      </div>

      {/* Right side - Form */}
      <div className="flex-1 flex flex-col items-center justify-center p-6 relative z-10">
        {/* Mobile logo */}
        <Link href="/" className="mb-8 lg:hidden">
          <div className="bg-white rounded-xl px-4 py-2">
            <Image 
              src="/logo-chat2pay.png" 
              alt="Chat2Pay" 
              width={160} 
              height={52}
              className="h-12 w-auto"
            />
          </div>
        </Link>

        {/* Login Card */}
        <div className="w-full max-w-md">
          <div className="glass-card rounded-2xl p-8">
            <div className="text-center mb-8">
              <h2 className="text-2xl font-bold mb-2">Selamat Datang! 👋</h2>
              <p className="text-sm text-muted-foreground">
                Masuk ke akun Anda untuk melanjutkan
              </p>
            </div>

            <LoginForm />
          </div>

          <p className="mt-6 text-center text-sm text-muted-foreground lg:hidden">
            &copy; 2024 Chat2Pay by Code Crafters
          </p>
        </div>
      </div>
    </div>
  );
}
