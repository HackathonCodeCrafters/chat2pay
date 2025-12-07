"use client";

import Link from "next/link";
import { useState, type FormEvent } from "react";
import { ArrowLeftCircle, Lock, Sparkles } from "lucide-react";

import { Button } from "@/components/ui/button";
import { ApiError, apiClient } from "@/lib/api";

type LoginPayload = {
  email: string;
  password: string;
};

type Status =
  | { type: "idle" }
  | { type: "loading" }
  | { type: "success"; message: string }
  | { type: "error"; message: string };

const inputClass =
  "w-full rounded-lg border border-zinc-200 bg-white px-4 py-3 text-sm shadow-sm outline-none ring-0 transition focus:border-zinc-800 focus:ring-2 focus:ring-zinc-900/10 dark:border-zinc-700 dark:bg-zinc-900 dark:text-zinc-50 dark:focus:border-zinc-100 dark:focus:ring-zinc-100/10";

export default function LoginMerchantPage() {
  const [payload, setPayload] = useState<LoginPayload>({
    email: "",
    password: "",
  });
  const [status, setStatus] = useState<Status>({ type: "idle" });

  const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setStatus({ type: "loading" });

    try {
      await apiClient.post("/merchants/login", {
        json: {
          email: payload.email,
          password: payload.password,
        },
      });

      setStatus({
        type: "success",
        message: "Login berhasil. Mengarahkan ke dashboard...",
      });
    } catch (error) {
      const message =
        error instanceof ApiError
          ? error.payload && typeof error.payload === "object"
            ? (error.payload as { message?: string }).message ??
              "Login gagal, periksa kredensial Anda."
            : "Login gagal, periksa kredensial Anda."
          : "Login gagal, periksa kredensial Anda.";

      setStatus({ type: "error", message });
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-50 via-white to-amber-50 text-zinc-900 dark:from-black dark:via-zinc-900 dark:to-amber-950/20 dark:text-zinc-100">
      <div className="mx-auto flex min-h-screen max-w-6xl flex-col gap-12 px-6 py-12 lg:flex-row lg:items-center lg:gap-16">
        <div className="flex-1 space-y-6">
          <Link
            href="/"
            className="inline-flex items-center gap-2 text-sm font-semibold text-amber-700 underline-offset-4 hover:underline dark:text-amber-300"
          >
            <ArrowLeftCircle className="h-4 w-4" />
            Kembali ke beranda
          </Link>
          <div className="space-y-4">
            <div className="inline-flex items-center gap-2 rounded-full bg-zinc-900 px-3 py-1 text-xs font-semibold uppercase tracking-[0.16em] text-white shadow-lg ring-1 ring-zinc-900/10 dark:bg-white dark:text-zinc-950">
              <Lock className="h-4 w-4" />
              Login Merchant
            </div>
            <h1 className="text-4xl font-semibold leading-tight">
              Masuk ke dashboard pembayaran Anda
            </h1>
            <p className="max-w-2xl text-base text-zinc-600 dark:text-zinc-300">
              Kelola pembayaran, pantau transaksi, dan atur channel chat dari
              satu tempat. Gunakan email bisnis yang sudah terdaftar.
            </p>
          </div>
          <div className="grid gap-4 sm:grid-cols-2">
            {[
              "Realtime insight transaksi",
              "Pengaturan webhook yang mudah",
              "Sesi aman dengan proteksi token",
              "Dukungan 24/7 lewat chat",
            ].map((item) => (
              <div
                key={item}
                className="flex items-center gap-3 rounded-xl border border-zinc-200/80 bg-white/70 px-4 py-3 text-sm font-medium shadow-sm backdrop-blur dark:border-zinc-800/80 dark:bg-zinc-900/70"
              >
                <Sparkles className="h-5 w-5 text-amber-500" />
                <span>{item}</span>
              </div>
            ))}
          </div>
        </div>

        <div className="flex-1">
          <div className="rounded-2xl border border-zinc-200/80 bg-white/80 p-6 shadow-2xl backdrop-blur dark:border-zinc-800/80 dark:bg-zinc-900/80 sm:p-8">
            <div className="mb-6 space-y-1">
              <p className="text-sm font-semibold text-amber-600 dark:text-amber-300">
                Akses akun
              </p>
              <h2 className="text-2xl font-semibold">Masuk sebagai merchant</h2>
              <p className="text-sm text-zinc-500 dark:text-zinc-400">
                Masukkan kredensial yang Anda gunakan saat mendaftar.
              </p>
            </div>

            <form className="space-y-4" onSubmit={handleSubmit}>
              <div className="space-y-2">
                <label className="text-sm font-medium text-zinc-700 dark:text-zinc-200">
                  Email bisnis
                </label>
                <input
                  required
                  type="email"
                  className={inputClass}
                  placeholder="ops@bisnis.co.id"
                  value={payload.email}
                  onChange={(event) =>
                    setPayload((prev) => ({ ...prev, email: event.target.value }))
                  }
                />
              </div>

              <div className="space-y-2">
                <label className="text-sm font-medium text-zinc-700 dark:text-zinc-200">
                  Password
                </label>
                <input
                  required
                  type="password"
                  className={inputClass}
                  placeholder="Masukkan password"
                  value={payload.password}
                  onChange={(event) =>
                    setPayload((prev) => ({
                      ...prev,
                      password: event.target.value,
                    }))
                  }
                />
              </div>

              {status.type === "error" ? (
                <div className="rounded-lg border border-rose-200 bg-rose-50 px-4 py-3 text-sm text-rose-800 dark:border-rose-900/60 dark:bg-rose-950/40 dark:text-rose-100">
                  {status.message}
                </div>
              ) : null}

              {status.type === "success" ? (
                <div className="rounded-lg border border-emerald-200 bg-emerald-50 px-4 py-3 text-sm text-emerald-800 dark:border-emerald-900/60 dark:bg-emerald-950/40 dark:text-emerald-100">
                  {status.message}
                </div>
              ) : null}

              <Button
                type="submit"
                className="w-full text-base font-semibold"
                disabled={status.type === "loading"}
              >
                {status.type === "loading" ? "Memproses..." : "Masuk"}
              </Button>
            </form>

            <p className="mt-6 text-sm text-zinc-600 dark:text-zinc-400">
              Belum punya akun?{" "}
              <Link
                href="/register"
                className="font-semibold text-amber-700 underline-offset-4 hover:underline dark:text-amber-300"
              >
                Daftar sekarang
              </Link>
              .
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}
