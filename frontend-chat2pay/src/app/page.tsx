"use client";

import Link from "next/link";
import Image from "next/image";
import { 
  MessageSquare, Sparkles, ShoppingCart, Zap, ArrowRight, 
  Bot, Shield, Clock, Star, ChevronRight, Play,
  Cpu, Layers, Globe, CheckCircle2
} from "lucide-react";
import { useAuth } from "@/features/auth";

const features = [
  {
    icon: Bot,
    title: "AI-Powered Search",
    description: "Cukup ceritakan apa yang Anda cari, AI kami akan menemukan produk yang sempurna untuk Anda.",
    color: "violet"
  },
  {
    icon: Zap,
    title: "Respons Instan",
    description: "Dapatkan rekomendasi produk dalam hitungan detik, tanpa perlu scroll panjang.",
    color: "cyan"
  },
  {
    icon: Shield,
    title: "Transaksi Aman",
    description: "Pembayaran terenkripsi dan perlindungan pembeli untuk setiap transaksi.",
    color: "emerald"
  },
  {
    icon: Clock,
    title: "24/7 Available",
    description: "AI Assistant kami siap membantu kapanpun Anda butuhkan.",
    color: "amber"
  }
];

const stats = [
  { value: "10K+", label: "Active Users" },
  { value: "50K+", label: "Products" },
  { value: "99%", label: "Satisfaction" },
  { value: "24/7", label: "AI Support" },
];

const steps = [
  { step: "01", title: "Daftar Akun", description: "Buat akun gratis dalam hitungan detik" },
  { step: "02", title: "Mulai Chat", description: "Ceritakan produk yang Anda cari ke AI" },
  { step: "03", title: "Pilih & Bayar", description: "Pilih produk dan checkout dengan mudah" },
];

export default function HomePage() {
  const { isAuthenticated, isLoading } = useAuth();

  return (
    <div className="min-h-screen flex flex-col relative overflow-hidden">
      {/* Animated Background */}
      <div className="fixed inset-0 overflow-hidden pointer-events-none">
        <div className="absolute top-0 left-1/4 w-[600px] h-[600px] bg-violet-500/20 rounded-full blur-[120px] animate-pulse" />
        <div className="absolute bottom-0 right-1/4 w-[500px] h-[500px] bg-cyan-500/20 rounded-full blur-[120px] animate-pulse" style={{ animationDelay: '1s' }} />
        <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[800px] h-[800px] bg-violet-500/5 rounded-full blur-[100px]" />
        {/* Grid overlay */}
        <div className="absolute inset-0 bg-[linear-gradient(rgba(255,255,255,0.02)_1px,transparent_1px),linear-gradient(90deg,rgba(255,255,255,0.02)_1px,transparent_1px)] bg-[size:64px_64px]" />
      </div>

      {/* Navbar */}
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
            <Link 
              href="/merchant/login"
              className="hidden sm:flex px-4 py-2 rounded-lg text-sm text-muted-foreground hover:text-foreground hover:bg-white/5 transition-all"
            >
              Merchant Portal
            </Link>
            {!isLoading && (
              <>
                {isAuthenticated ? (
                  <Link 
                    href="/chat"
                    className="flex items-center gap-2 px-5 py-2.5 rounded-xl text-sm font-medium bg-gradient-to-r from-violet-500 to-cyan-500 text-white hover:shadow-lg hover:shadow-violet-500/25 transition-all duration-300"
                  >
                    Mulai Chat
                    <ArrowRight className="h-4 w-4" />
                  </Link>
                ) : (
                  <>
                    <Link 
                      href="/login"
                      className="px-4 py-2 rounded-lg text-sm text-muted-foreground hover:text-foreground hover:bg-white/5 transition-all"
                    >
                      Login
                    </Link>
                    <Link 
                      href="/register"
                      className="flex items-center gap-2 px-5 py-2.5 rounded-xl text-sm font-medium bg-gradient-to-r from-violet-500 to-cyan-500 text-white hover:shadow-lg hover:shadow-violet-500/25 transition-all duration-300"
                    >
                      Get Started
                      <ArrowRight className="h-4 w-4" />
                    </Link>
                  </>
                )}
              </>
            )}
          </nav>
        </div>
      </header>

      <main className="flex-1 relative z-10">
        {/* Hero Section */}
        <section className="container pt-20 pb-32 md:pt-32 md:pb-40">
          <div className="mx-auto max-w-4xl text-center">
            {/* Badge */}
            <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full glass-card mb-8 animate-in fade-in slide-in-from-bottom-4 duration-700">
              <span className="flex h-2 w-2 rounded-full bg-emerald-500 animate-pulse" />
              <span className="text-sm text-muted-foreground">Powered by GPT-4 AI Technology</span>
              <ChevronRight className="h-4 w-4 text-muted-foreground" />
            </div>

            {/* Headline */}
            <h1 className="text-4xl sm:text-5xl md:text-6xl lg:text-7xl font-bold tracking-tight mb-6 animate-in fade-in slide-in-from-bottom-4 duration-700" style={{ animationDelay: '100ms' }}>
              Belanja Lebih Cerdas{" "}
              <br className="hidden sm:block" />
              dengan{" "}
              <span className="gradient-text">AI Assistant</span>
            </h1>

            {/* Subheadline */}
            <p className="text-lg md:text-xl text-muted-foreground max-w-2xl mx-auto mb-10 animate-in fade-in slide-in-from-bottom-4 duration-700" style={{ animationDelay: '200ms' }}>
              Temukan produk impian Anda hanya dengan bercerita. AI kami memahami 
              kebutuhan Anda dan merekomendasikan produk yang sempurna.
            </p>

            {/* CTA Buttons */}
            <div className="flex flex-col sm:flex-row gap-4 justify-center mb-16 animate-in fade-in slide-in-from-bottom-4 duration-700" style={{ animationDelay: '300ms' }}>
              <Link 
                href={isAuthenticated ? "/chat" : "/register"}
                className="group flex items-center justify-center gap-2 px-8 py-4 rounded-2xl text-base font-medium bg-gradient-to-r from-violet-500 to-cyan-500 text-white hover:shadow-2xl hover:shadow-violet-500/30 transition-all duration-300 hover:-translate-y-1"
              >
                <MessageSquare className="h-5 w-5" />
                {isAuthenticated ? "Mulai Chat Sekarang" : "Coba Gratis Sekarang"}
                <ArrowRight className="h-4 w-4 group-hover:translate-x-1 transition-transform" />
              </Link>
              <Link 
                href="#demo"
                className="group flex items-center justify-center gap-2 px-8 py-4 rounded-2xl text-base font-medium glass-card glass-card-hover"
              >
                <Play className="h-5 w-5" />
                Lihat Demo
              </Link>
            </div>

            {/* Stats */}
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4 max-w-2xl mx-auto animate-in fade-in slide-in-from-bottom-4 duration-700" style={{ animationDelay: '400ms' }}>
              {stats.map((stat) => (
                <div key={stat.label} className="glass-card rounded-xl p-4 text-center">
                  <p className="text-2xl md:text-3xl font-bold gradient-text">{stat.value}</p>
                  <p className="text-xs text-muted-foreground">{stat.label}</p>
                </div>
              ))}
            </div>
          </div>
        </section>

        {/* Features Section */}
        <section className="container py-24">
          <div className="text-center mb-16">
            <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full glass-card mb-6">
              <Layers className="h-4 w-4 text-violet-400" />
              <span className="text-sm text-muted-foreground">Features</span>
            </div>
            <h2 className="text-3xl md:text-4xl font-bold mb-4">
              Mengapa Memilih{" "}
              <span className="gradient-text">Chat2Pay</span>?
            </h2>
            <p className="text-muted-foreground max-w-xl mx-auto">
              Platform belanja berbasis AI pertama di Indonesia yang memahami kebutuhan Anda
            </p>
          </div>

          <div className="grid md:grid-cols-2 lg:grid-cols-4 gap-6 max-w-6xl mx-auto">
            {features.map((feature, index) => (
              <div 
                key={feature.title}
                className="glass-card glass-card-hover rounded-2xl p-6 group"
                style={{ animationDelay: `${index * 100}ms` }}
              >
                <div className={`flex h-14 w-14 items-center justify-center rounded-xl bg-${feature.color}-500/10 border border-${feature.color}-500/20 mb-5 group-hover:scale-110 transition-transform duration-300`}>
                  <feature.icon className={`h-7 w-7 text-${feature.color}-400`} />
                </div>
                <h3 className="text-lg font-semibold mb-2">{feature.title}</h3>
                <p className="text-sm text-muted-foreground leading-relaxed">{feature.description}</p>
              </div>
            ))}
          </div>
        </section>

        {/* How It Works */}
        <section className="container py-24">
          <div className="text-center mb-16">
            <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full glass-card mb-6">
              <Cpu className="h-4 w-4 text-cyan-400" />
              <span className="text-sm text-muted-foreground">How It Works</span>
            </div>
            <h2 className="text-3xl md:text-4xl font-bold mb-4">
              Semudah{" "}
              <span className="gradient-text">1, 2, 3</span>
            </h2>
            <p className="text-muted-foreground max-w-xl mx-auto">
              Mulai belanja dengan AI dalam 3 langkah sederhana
            </p>
          </div>

          <div className="grid md:grid-cols-3 gap-8 max-w-4xl mx-auto">
            {steps.map((step, index) => (
              <div key={step.step} className="relative">
                {index < steps.length - 1 && (
                  <div className="hidden md:block absolute top-12 left-full w-full h-px bg-gradient-to-r from-violet-500/50 to-transparent -translate-x-1/2 z-0" />
                )}
                <div className="glass-card rounded-2xl p-6 text-center relative z-10">
                  <div className="flex h-16 w-16 items-center justify-center rounded-2xl bg-gradient-to-br from-violet-500 to-cyan-500 mx-auto mb-5 shadow-lg shadow-violet-500/25">
                    <span className="text-2xl font-bold text-white">{step.step}</span>
                  </div>
                  <h3 className="text-lg font-semibold mb-2">{step.title}</h3>
                  <p className="text-sm text-muted-foreground">{step.description}</p>
                </div>
              </div>
            ))}
          </div>
        </section>

        {/* Demo/Preview Section */}
        <section id="demo" className="container py-24">
          <div className="glass-card rounded-3xl p-8 md:p-12 max-w-5xl mx-auto overflow-hidden relative">
            {/* Decorative elements */}
            <div className="absolute top-0 right-0 w-64 h-64 bg-violet-500/20 rounded-full blur-3xl" />
            <div className="absolute bottom-0 left-0 w-64 h-64 bg-cyan-500/20 rounded-full blur-3xl" />
            
            <div className="relative z-10 grid md:grid-cols-2 gap-8 items-center">
              <div>
                <div className="inline-flex items-center gap-2 px-3 py-1.5 rounded-full bg-emerald-500/10 border border-emerald-500/30 text-emerald-400 text-sm mb-6">
                  <CheckCircle2 className="h-4 w-4" />
                  Live Demo
                </div>
                <h2 className="text-3xl md:text-4xl font-bold mb-4">
                  Lihat AI Beraksi
                </h2>
                <p className="text-muted-foreground mb-6">
                  Tanyakan apa saja tentang produk yang Anda cari. AI kami akan memberikan 
                  rekomendasi personal berdasarkan kebutuhan Anda.
                </p>
                <ul className="space-y-3 mb-8">
                  {[
                    "Rekomendasi produk berdasarkan budget",
                    "Perbandingan fitur antar produk",
                    "Saran produk sesuai kebutuhan spesifik"
                  ].map((item) => (
                    <li key={item} className="flex items-center gap-3 text-sm">
                      <div className="flex h-5 w-5 items-center justify-center rounded-full bg-gradient-to-r from-violet-500 to-cyan-500">
                        <CheckCircle2 className="h-3 w-3 text-white" />
                      </div>
                      {item}
                    </li>
                  ))}
                </ul>
                <Link 
                  href={isAuthenticated ? "/chat" : "/register"}
                  className="inline-flex items-center gap-2 px-6 py-3 rounded-xl text-sm font-medium bg-gradient-to-r from-violet-500 to-cyan-500 text-white hover:shadow-lg hover:shadow-violet-500/25 transition-all"
                >
                  Coba Sekarang
                  <ArrowRight className="h-4 w-4" />
                </Link>
              </div>

              {/* Chat Preview */}
              <div className="glass-card rounded-2xl p-4 border border-white/10">
                <div className="flex items-center gap-3 pb-4 border-b border-white/10">
                  <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-gradient-to-br from-violet-500 to-cyan-500">
                    <Bot className="h-5 w-5 text-white" />
                  </div>
                  <div>
                    <p className="font-medium text-sm">AI Assistant</p>
                    <p className="text-xs text-emerald-400 flex items-center gap-1">
                      <span className="h-1.5 w-1.5 rounded-full bg-emerald-500" />
                      Online
                    </p>
                  </div>
                </div>
                <div className="py-4 space-y-4">
                  {/* User message */}
                  <div className="flex justify-end">
                    <div className="chat-bubble-user rounded-2xl rounded-tr-md px-4 py-2.5 max-w-[80%]">
                      <p className="text-sm">Carikan laptop gaming budget 15 juta</p>
                    </div>
                  </div>
                  {/* AI response */}
                  <div className="flex gap-3">
                    <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-violet-500/20 shrink-0">
                      <Bot className="h-4 w-4 text-violet-400" />
                    </div>
                    <div className="chat-bubble-assistant rounded-2xl rounded-tl-md px-4 py-2.5">
                      <p className="text-sm">Saya menemukan 3 laptop gaming terbaik untuk budget Anda! Berikut rekomendasinya... 🎮</p>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </section>

        {/* Final CTA */}
        <section className="container py-24">
          <div className="relative rounded-3xl overflow-hidden">
            {/* Animated gradient background */}
            <div className="absolute inset-0 bg-gradient-to-r from-violet-600 via-cyan-600 to-violet-600 bg-[length:200%_100%] animate-[gradient-shift_5s_ease_infinite]" />
            <div className="absolute inset-0 bg-[linear-gradient(rgba(255,255,255,0.1)_1px,transparent_1px),linear-gradient(90deg,rgba(255,255,255,0.1)_1px,transparent_1px)] bg-[size:32px_32px]" />
            
            <div className="relative z-10 text-center py-16 md:py-24 px-6">
              <h2 className="text-3xl md:text-4xl lg:text-5xl font-bold text-white mb-4">
                Siap Belanja dengan Cara Baru?
              </h2>
              <p className="text-white/80 max-w-xl mx-auto mb-8 text-lg">
                Bergabung dengan ribuan pengguna yang sudah merasakan kemudahan 
                belanja dengan AI Assistant
              </p>
              <div className="flex flex-col sm:flex-row gap-4 justify-center">
                <Link 
                  href="/register"
                  className="inline-flex items-center justify-center gap-2 px-8 py-4 rounded-2xl text-base font-medium bg-white text-violet-600 hover:bg-white/90 transition-all shadow-2xl shadow-black/20"
                >
                  <ShoppingCart className="h-5 w-5" />
                  Mulai Gratis Sekarang
                </Link>
                <Link 
                  href="/merchant/register"
                  className="inline-flex items-center justify-center gap-2 px-8 py-4 rounded-2xl text-base font-medium bg-white/10 text-white border border-white/20 hover:bg-white/20 transition-all"
                >
                  <Globe className="h-5 w-5" />
                  Daftar sebagai Merchant
                </Link>
              </div>
            </div>
          </div>
        </section>
      </main>

      {/* Footer */}
      <footer className="relative z-10 border-t border-white/5 glass">
        <div className="container py-12">
          <div className="grid md:grid-cols-4 gap-8 mb-8">
            <div className="md:col-span-2">
              <Link href="/" className="flex items-center mb-4">
                <div className="bg-white rounded-xl px-4 py-2">
                  <Image 
                    src="/logo-chat2pay.png" 
                    alt="Chat2Pay" 
                    width={140} 
                    height={45}
                    className="h-10 w-auto"
                  />
                </div>
              </Link>
              <p className="text-sm text-muted-foreground max-w-xs">
                Platform belanja berbasis AI pertama di Indonesia. Temukan produk impian Anda dengan mudah.
              </p>
            </div>
            <div>
              <h4 className="font-medium mb-4">Product</h4>
              <ul className="space-y-2 text-sm text-muted-foreground">
                <li><Link href="/chat" className="hover:text-foreground transition-colors">AI Chat</Link></li>
                <li><Link href="/orders" className="hover:text-foreground transition-colors">Orders</Link></li>
                <li><Link href="/cart" className="hover:text-foreground transition-colors">Cart</Link></li>
              </ul>
            </div>
            <div>
              <h4 className="font-medium mb-4">Merchant</h4>
              <ul className="space-y-2 text-sm text-muted-foreground">
                <li><Link href="/merchant/register" className="hover:text-foreground transition-colors">Daftar Merchant</Link></li>
                <li><Link href="/merchant/login" className="hover:text-foreground transition-colors">Login Merchant</Link></li>
                <li><Link href="/merchant/dashboard" className="hover:text-foreground transition-colors">Dashboard</Link></li>
              </ul>
            </div>
          </div>
          <div className="pt-8 border-t border-white/5 flex flex-col md:flex-row justify-between items-center gap-4">
            <p className="text-sm text-muted-foreground">
              &copy; 2024 Chat2Pay by Code Crafters. All rights reserved.
            </p>
            <div className="flex items-center gap-4">
              <span className="text-xs text-muted-foreground">Powered by</span>
              <div className="flex items-center gap-2 px-3 py-1.5 rounded-lg glass-card">
                <Cpu className="h-4 w-4 text-violet-400" />
                <span className="text-xs font-medium">GPT-4</span>
              </div>
            </div>
          </div>
        </div>
      </footer>
    </div>
  );
}
