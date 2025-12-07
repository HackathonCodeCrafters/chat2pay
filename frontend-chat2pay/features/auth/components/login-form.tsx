"use client";

import * as React from "react";
import Link from "next/link";
import { Eye, EyeOff, Mail, Lock, AlertCircle, ArrowRight } from "lucide-react";
import { Spinner } from "@/shared/components/atoms";
import { useAuth } from "../hooks/use-auth";
import { cn } from "@/shared/lib/utils";

export function LoginForm() {
  const { login, isLoading } = useAuth();
  const [showPassword, setShowPassword] = React.useState(false);
  const [error, setError] = React.useState<string | null>(null);
  const [focusedField, setFocusedField] = React.useState<string | null>(null);

  const [formData, setFormData] = React.useState({
    email: "",
    password: "",
  });

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
    setError(null);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);

    if (!formData.email || !formData.password) {
      setError("Harap isi semua field");
      return;
    }

    try {
      await login(formData);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Login gagal. Silakan coba lagi.");
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-5">
      {error && (
        <div className="flex items-center gap-3 p-3 rounded-xl bg-red-500/10 border border-red-500/20">
          <AlertCircle className="h-4 w-4 text-red-400 shrink-0" />
          <p className="text-sm text-red-400">{error}</p>
        </div>
      )}

      <div className="space-y-2">
        <label htmlFor="email" className="text-sm font-medium">
          Email
        </label>
        <div 
          className={cn(
            "flex items-center gap-3 px-4 py-3 rounded-xl border transition-all duration-300",
            focusedField === "email" 
              ? "border-violet-500/50 bg-violet-500/5 shadow-lg shadow-violet-500/10" 
              : "border-white/10 bg-white/5"
          )}
        >
          <Mail className={cn(
            "h-4 w-4 transition-colors",
            focusedField === "email" ? "text-violet-400" : "text-muted-foreground"
          )} />
          <input
            id="email"
            name="email"
            type="email"
            placeholder="nama@email.com"
            value={formData.email}
            onChange={handleChange}
            onFocus={() => setFocusedField("email")}
            onBlur={() => setFocusedField(null)}
            className="flex-1 bg-transparent border-none outline-none text-sm placeholder:text-muted-foreground"
            autoComplete="email"
            disabled={isLoading}
          />
        </div>
      </div>

      <div className="space-y-2">
        <label htmlFor="password" className="text-sm font-medium">
          Password
        </label>
        <div 
          className={cn(
            "flex items-center gap-3 px-4 py-3 rounded-xl border transition-all duration-300",
            focusedField === "password" 
              ? "border-violet-500/50 bg-violet-500/5 shadow-lg shadow-violet-500/10" 
              : "border-white/10 bg-white/5"
          )}
        >
          <Lock className={cn(
            "h-4 w-4 transition-colors",
            focusedField === "password" ? "text-violet-400" : "text-muted-foreground"
          )} />
          <input
            id="password"
            name="password"
            type={showPassword ? "text" : "password"}
            placeholder="Masukkan password"
            value={formData.password}
            onChange={handleChange}
            onFocus={() => setFocusedField("password")}
            onBlur={() => setFocusedField(null)}
            className="flex-1 bg-transparent border-none outline-none text-sm placeholder:text-muted-foreground"
            autoComplete="current-password"
            disabled={isLoading}
          />
          <button
            type="button"
            onClick={() => setShowPassword(!showPassword)}
            className="text-muted-foreground hover:text-foreground transition-colors"
          >
            {showPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
          </button>
        </div>
      </div>

      <button 
        type="submit" 
        disabled={isLoading}
        className="w-full flex items-center justify-center gap-2 py-3 rounded-xl font-medium bg-gradient-to-r from-violet-500 to-cyan-500 text-white hover:shadow-lg hover:shadow-violet-500/25 transition-all duration-300 disabled:opacity-50"
      >
        {isLoading ? (
          <Spinner size="sm" />
        ) : (
          <>
            Masuk
            <ArrowRight className="h-4 w-4" />
          </>
        )}
      </button>

      <div className="relative">
        <div className="absolute inset-0 flex items-center">
          <div className="w-full border-t border-white/10" />
        </div>
        <div className="relative flex justify-center text-xs">
          <span className="bg-background px-2 text-muted-foreground">atau</span>
        </div>
      </div>

      <p className="text-center text-sm text-muted-foreground">
        Belum punya akun?{" "}
        <Link
          href="/register"
          className="font-medium text-violet-400 hover:text-violet-300 transition-colors"
        >
          Daftar sekarang
        </Link>
      </p>

      <p className="text-center text-xs text-muted-foreground pt-2">
        <Link href="/merchant/login" className="hover:text-foreground transition-colors">
          Login sebagai Merchant →
        </Link>
      </p>
    </form>
  );
}
