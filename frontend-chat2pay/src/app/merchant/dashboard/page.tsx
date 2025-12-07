"use client";

import * as React from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import Image from "next/image";
import { 
  Store, LogOut, Package, BarChart3, Settings, ShoppingCart, 
  TrendingUp, Users, DollarSign, ArrowUpRight, Sparkles, Bell,
  Plus, Search, MoreHorizontal
} from "lucide-react";
import { Button, Spinner, Badge } from "@/shared/components/atoms";
import { useMerchantAuth, ProductList, AddProductModal, getProducts, createProduct } from "@/features/merchant";
import type { Product, CreateProductData } from "@/features/merchant";

export default function MerchantDashboardPage() {
  const router = useRouter();
  const { user, token, isAuthenticated, isLoading: authLoading, logout } = useMerchantAuth();
  
  const [products, setProducts] = React.useState<Product[]>([]);
  const [isLoadingProducts, setIsLoadingProducts] = React.useState(true);
  const [isAddModalOpen, setIsAddModalOpen] = React.useState(false);
  const [isCreating, setIsCreating] = React.useState(false);

  React.useEffect(() => {
    if (!authLoading && !isAuthenticated) {
      router.push("/merchant/login");
    }
  }, [authLoading, isAuthenticated, router]);

  React.useEffect(() => {
    if (token && user?.merchant_id) {
      loadProducts();
    }
  }, [token, user?.merchant_id]);

  const loadProducts = async () => {
    if (!token || !user?.merchant_id) return;
    
    setIsLoadingProducts(true);
    try {
      const data = await getProducts(token, user.merchant_id);
      setProducts(data || []);
    } catch (error) {
      console.error("Failed to load products:", error);
      setProducts([]);
    } finally {
      setIsLoadingProducts(false);
    }
  };

  const handleCreateProduct = async (data: CreateProductData) => {
    if (!token || !user?.merchant_id) return;
    
    setIsCreating(true);
    try {
      await createProduct(token, user.merchant_id, data);
      await loadProducts();
    } finally {
      setIsCreating(false);
    }
  };

  if (authLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <Spinner size="lg" />
      </div>
    );
  }

  if (!isAuthenticated) {
    return null;
  }

  const stats = [
    {
      label: "Total Produk",
      value: products.length,
      icon: Package,
      trend: "+12%",
      trendUp: true,
      color: "violet"
    },
    {
      label: "Produk Aktif",
      value: products.filter(p => p.status === "active").length,
      icon: TrendingUp,
      trend: "+8%",
      trendUp: true,
      color: "emerald"
    },
    {
      label: "Total Views",
      value: "2.4K",
      icon: Users,
      trend: "+24%",
      trendUp: true,
      color: "cyan"
    },
    {
      label: "Revenue",
      value: "Rp 12.5M",
      icon: DollarSign,
      trend: "+18%",
      trendUp: true,
      color: "amber"
    },
  ];

  return (
    <div className="min-h-screen">
      {/* Animated background */}
      <div className="fixed inset-0 overflow-hidden pointer-events-none">
        <div className="absolute top-0 left-1/4 w-96 h-96 bg-violet-500/10 rounded-full blur-3xl" />
        <div className="absolute bottom-0 right-1/4 w-96 h-96 bg-cyan-500/10 rounded-full blur-3xl" />
      </div>

      {/* Navbar */}
      <header className="sticky top-0 z-50 glass border-b border-white/5">
        <div className="container flex h-16 items-center justify-between">
          <Link href="/merchant/dashboard" className="flex items-center gap-3">
            <div className="bg-white rounded-xl px-3 py-1.5">
              <Image 
                src="/logo-chat2pay.png" 
                alt="Chat2Pay" 
                width={120} 
                height={40}
                className="h-9 w-auto"
              />
            </div>
            <div className="hidden sm:block h-6 w-px bg-white/10" />
            <span className="hidden sm:block text-sm text-muted-foreground">Merchant Portal</span>
          </Link>

          <div className="flex items-center gap-3">
            {/* Search */}
            <div className="hidden md:flex items-center gap-2 px-3 py-2 rounded-xl glass-card">
              <Search className="h-4 w-4 text-muted-foreground" />
              <input 
                type="text" 
                placeholder="Search..." 
                className="bg-transparent border-none outline-none text-sm w-40 placeholder:text-muted-foreground"
              />
            </div>

            {/* Notifications */}
            <button className="relative flex items-center justify-center h-10 w-10 rounded-xl glass-card-hover">
              <Bell className="h-4 w-4 text-muted-foreground" />
              <span className="absolute -top-1 -right-1 h-4 w-4 rounded-full bg-gradient-to-r from-violet-500 to-cyan-500 text-[10px] font-bold text-white flex items-center justify-center">
                3
              </span>
            </button>

            {/* Divider */}
            <div className="h-8 w-px bg-white/10 mx-1" />

            {/* User */}
            <div className="hidden sm:flex items-center gap-3">
              <div className="flex flex-col items-end">
                <span className="text-sm font-medium">{user?.merchant?.name || user?.name}</span>
                <span className="text-[10px] text-muted-foreground">{user?.email}</span>
              </div>
              <div className="relative">
                <div className="h-10 w-10 rounded-xl bg-gradient-to-br from-violet-500/20 to-cyan-500/20 flex items-center justify-center border border-white/10">
                  <span className="text-sm font-medium gradient-text">
                    {(user?.name || "M").charAt(0).toUpperCase()}
                  </span>
                </div>
                <div className="absolute bottom-0 right-0 h-2.5 w-2.5 rounded-full bg-emerald-500 border-2 border-background" />
              </div>
            </div>
            
            <button 
              onClick={logout}
              className="flex items-center justify-center h-10 w-10 rounded-xl text-muted-foreground hover:text-red-400 hover:bg-red-500/10 transition-all"
            >
              <LogOut className="h-4 w-4" />
            </button>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="container py-8 relative z-10">
        {/* Welcome Section */}
        <div className="mb-8">
          <h1 className="text-2xl font-bold mb-1">
            Selamat datang, <span className="gradient-text">{user?.name || "Merchant"}</span>! 👋
          </h1>
          <p className="text-muted-foreground">
            Berikut adalah ringkasan performa toko Anda hari ini.
          </p>
        </div>

        {/* Stats Grid */}
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4 mb-8">
          {stats.map((stat, index) => (
            <div 
              key={stat.label}
              className="stat-card p-5 glass-card-hover cursor-pointer"
              style={{ animationDelay: `${index * 100}ms` }}
            >
              <div className="flex items-start justify-between mb-4">
                <div className={`flex h-12 w-12 items-center justify-center rounded-xl bg-${stat.color}-500/10 border border-${stat.color}-500/20`}>
                  <stat.icon className={`h-6 w-6 text-${stat.color}-400`} />
                </div>
                <div className={`flex items-center gap-1 text-xs font-medium ${stat.trendUp ? 'text-emerald-400' : 'text-red-400'}`}>
                  <ArrowUpRight className="h-3 w-3" />
                  {stat.trend}
                </div>
              </div>
              <p className="text-2xl font-bold mb-1">{stat.value}</p>
              <p className="text-sm text-muted-foreground">{stat.label}</p>
            </div>
          ))}
        </div>

        {/* Quick Actions */}
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3 mb-8">
          <button
            onClick={() => setIsAddModalOpen(true)}
            className="glass-card glass-card-hover rounded-xl p-5 text-left group"
          >
            <div className="flex items-center gap-4">
              <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-gradient-to-br from-violet-500 to-cyan-500 shadow-lg shadow-violet-500/20 group-hover:shadow-violet-500/40 transition-all">
                <Plus className="h-6 w-6 text-white" />
              </div>
              <div>
                <p className="font-medium mb-0.5">Tambah Produk</p>
                <p className="text-xs text-muted-foreground">Tambahkan produk baru ke toko</p>
              </div>
            </div>
          </button>

          <button
            onClick={() => router.push("/merchant/orders")}
            className="glass-card glass-card-hover rounded-xl p-5 text-left group"
          >
            <div className="flex items-center gap-4">
              <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-gradient-to-br from-cyan-500 to-emerald-500 shadow-lg shadow-cyan-500/20 group-hover:shadow-cyan-500/40 transition-all">
                <ShoppingCart className="h-6 w-6 text-white" />
              </div>
              <div>
                <p className="font-medium mb-0.5">Kelola Pesanan</p>
                <p className="text-xs text-muted-foreground">Lihat dan proses pesanan masuk</p>
              </div>
            </div>
          </button>

          <button
            className="glass-card glass-card-hover rounded-xl p-5 text-left group"
          >
            <div className="flex items-center gap-4">
              <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-gradient-to-br from-amber-500 to-orange-500 shadow-lg shadow-amber-500/20 group-hover:shadow-amber-500/40 transition-all">
                <BarChart3 className="h-6 w-6 text-white" />
              </div>
              <div>
                <p className="font-medium mb-0.5">Analitik</p>
                <p className="text-xs text-muted-foreground">Lihat performa penjualan</p>
              </div>
            </div>
          </button>
        </div>

        {/* Store Status Banner */}
        <div className="glass-card rounded-xl p-4 mb-8 flex items-center justify-between">
          <div className="flex items-center gap-4">
            <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-emerald-500/10 border border-emerald-500/20">
              <Store className="h-5 w-5 text-emerald-400" />
            </div>
            <div>
              <p className="font-medium">Status Toko</p>
              <p className="text-xs text-muted-foreground">{user?.merchant?.name || "Merchant Store"}</p>
            </div>
          </div>
          <div className="flex items-center gap-3">
            <span className={`px-3 py-1.5 rounded-lg text-xs font-medium ${
              user?.merchant?.status === "active" 
                ? "badge-success" 
                : "badge-warning"
            }`}>
              {user?.merchant?.status === "active" ? "Aktif" : "Pending"}
            </span>
            <button className="flex items-center justify-center h-8 w-8 rounded-lg hover:bg-white/5 transition-all">
              <MoreHorizontal className="h-4 w-4 text-muted-foreground" />
            </button>
          </div>
        </div>

        {/* Products Section */}
        <div className="glass-card rounded-xl overflow-hidden">
          <div className="p-5 border-b border-white/5 flex items-center justify-between">
            <div>
              <h2 className="font-semibold text-lg">Produk Anda</h2>
              <p className="text-sm text-muted-foreground">{products.length} produk terdaftar</p>
            </div>
            <button
              onClick={() => setIsAddModalOpen(true)}
              className="flex items-center gap-2 px-4 py-2 rounded-xl text-sm font-medium bg-gradient-to-r from-violet-500 to-cyan-500 text-white hover:shadow-lg hover:shadow-violet-500/25 transition-all"
            >
              <Plus className="h-4 w-4" />
              Tambah Produk
            </button>
          </div>
          
          <div className="p-5">
            <ProductList
              products={products}
              isLoading={isLoadingProducts}
              onAddProduct={() => setIsAddModalOpen(true)}
            />
          </div>
        </div>
      </main>

      {/* Add Product Modal */}
      <AddProductModal
        isOpen={isAddModalOpen}
        isLoading={isCreating}
        onClose={() => setIsAddModalOpen(false)}
        onSubmit={handleCreateProduct}
      />
    </div>
  );
}
