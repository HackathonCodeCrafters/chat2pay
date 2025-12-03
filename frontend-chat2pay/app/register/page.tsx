"use client";

import Link from "next/link";
import { useState, type FormEvent } from "react";
import { CheckCircle2, ShieldCheck } from "lucide-react";

import { Button } from "@/components/ui/button";
import { ApiError, apiClient } from "@/lib/api";

type RegisterPayload = {
  businessName: string;
  email: string;
  phone: string;
  password: string;
  confirmPassword: string;
};

type Status =
  | { type: "idle" }
  | { type: "loading" }
  | { type: "success"; message: string }
  | { type: "error"; message: string };

const inputClass =
  "w-full rounded-lg border border-zinc-200 bg-white px-4 py-3 text-sm shadow-sm outline-none ring-0 transition focus:border-zinc-800 focus:ring-2 focus:ring-zinc-900/10 dark:border-zinc-700 dark:bg-zinc-900 dark:text-zinc-50 dark:focus:border-zinc-100 dark:focus:ring-zinc-100/10";

export default function RegisterMerchantPage() {
  const [payload, setPayload] = useState<RegisterPayload>({
    businessName: "",
    email: "",
    phone: "",
    password: "",
    confirmPassword: "",
  });
  const [status, setStatus] = useState<Status>({ type: "idle" });

  const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();

    if (payload.password !== payload.confirmPassword) {
      setStatus({ type: "error", message: "Konfirmasi password tidak cocok." });
      return;
    }

    setStatus({ type: "loading" });
    try {
      await apiClient.post("/merchants/register", {
        json: {
          business_name: payload.businessName,
          email: payload.email,
          phone: payload.phone,
          password: payload.password,
        },
      });

      setStatus({
        type: "success",
        message: "Akun merchant berhasil dibuat. Silakan login.",
      });
    } catch (error) {
      const message =
        error instanceof ApiError
          ? error.payload && typeof error.payload === "object"
            ? (error.payload as { message?: string }).message ??
              "Pendaftaran gagal, coba lagi."
            : "Pendaftaran gagal, coba lagi."
          : "Pendaftaran gagal, coba lagi.";

      setStatus({ type: "error", message });
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-amber-50 via-white to-slate-50 text-zinc-900 dark:from-zinc-950 dark:via-black dark:to-zinc-900 dark:text-zinc-100">
      <div className="mx-auto flex min-h-screen max-w-6xl flex-col gap-10 px-6 py-10 lg:flex-row lg:items-center lg:gap-16 lg:py-14">
        <div className="flex-1 space-y-6">
          <Link
            href="/"
            className="inline-flex items-center gap-2 text-sm font-semibold text-amber-700 underline-offset-4 hover:underline dark:text-amber-300"
          >
            &larr; Kembali ke beranda
          </Link>
          <div className="space-y-4">
            <span className="inline-flex items-center gap-2 rounded-full bg-amber-100 px-3 py-1 text-xs font-semibold uppercase tracking-[0.16em] text-amber-700 shadow-sm dark:bg-amber-500/15 dark:text-amber-200">
              <ShieldCheck className="h-4 w-4" />
              Merchant Onboarding
            </span>
            <div className="space-y-3">
              <h1 className="text-4xl font-semibold leading-tight">
                Daftar Merchant Chat2Pay
              </h1>
              <p className="max-w-2xl text-base text-zinc-600 dark:text-zinc-300">
                Buat akun merchant untuk mulai menerima pembayaran lewat chat.
                Kami hanya butuh beberapa detail bisnis Anda untuk menyiapkan
                dashboard dan kredensial API.
              </p>
            </div>
          </div>
          <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
            {[
              "Verifikasi cepat & sandbox siap pakai",
              "Atur tim & role dengan granular",
              "Koneksi API yang konsisten",
              "Rekonsiliasi payout otomatis",
            ].map((item) => (
              <div
                key={item}
                className="flex items-start gap-3 rounded-xl border border-zinc-200/80 bg-white/70 px-4 py-3 text-sm font-medium shadow-sm backdrop-blur dark:border-zinc-800/80 dark:bg-zinc-900/60"
              >
                <CheckCircle2 className="mt-0.5 h-5 w-5 text-emerald-500" />
                <span>{item}</span>
              </div>
            ))}
          </div>
        </div>

        <div className="flex-1">
          <div className="rounded-2xl border border-zinc-200/80 bg-white/80 p-6 shadow-xl backdrop-blur dark:border-zinc-800/80 dark:bg-zinc-900/80 sm:p-8">
            <div className="mb-6 space-y-1">
              <p className="text-sm font-semibold text-amber-600 dark:text-amber-300">
                Formulir Pendaftaran
              </p>
              <h2 className="text-2xl font-semibold">Informasi bisnis</h2>
              <p className="text-sm text-zinc-500 dark:text-zinc-400">
                Lengkapi detail di bawah untuk membuat akun merchant Anda.
              </p>
            </div>

            <form className="space-y-4" onSubmit={handleSubmit}>
              <div className="space-y-2">
                <label className="text-sm font-medium text-zinc-700 dark:text-zinc-200">
                  Nama bisnis
                </label>
                <input
                  required
                  className={inputClass}
                  placeholder="Contoh: PT Maju Jaya"
                  value={payload.businessName}
                  onChange={(event) =>
                    setPayload((prev) => ({
                      ...prev,
                      businessName: event.target.value,
                    }))
                  }
                />
              </div>

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
                    setPayload((prev) => ({
                      ...prev,
                      email: event.target.value,
                    }))
                  }
                />
              </div>

              <div className="space-y-2">
                <label className="text-sm font-medium text-zinc-700 dark:text-zinc-200">
                  Nomor kontak
                </label>
                <input
                  required
                  type="tel"
                  className={inputClass}
                  placeholder="+62 812 3456 7890"
                  value={payload.phone}
                  onChange={(event) =>
                    setPayload((prev) => ({
                      ...prev,
                      phone: event.target.value,
                    }))
                  }
                />
              </div>

              <div className="grid gap-4 sm:grid-cols-2">
                <div className="space-y-2">
                  <label className="text-sm font-medium text-zinc-700 dark:text-zinc-200">
                    Password
                  </label>
                  <input
                    required
                    type="password"
                    className={inputClass}
                    placeholder="Minimal 8 karakter"
                    value={payload.password}
                    onChange={(event) =>
                      setPayload((prev) => ({
                        ...prev,
                        password: event.target.value,
                      }))
                    }
                  />
                </div>
                <div className="space-y-2">
                  <label className="text-sm font-medium text-zinc-700 dark:text-zinc-200">
                    Konfirmasi password
                  </label>
                  <input
                    required
                    type="password"
                    className={inputClass}
                    placeholder="Ulangi password"
                    value={payload.confirmPassword}
                    onChange={(event) =>
                      setPayload((prev) => ({
                        ...prev,
                        confirmPassword: event.target.value,
                      }))
                    }
                  />
                </div>
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
                {status.type === "loading" ? "Memproses..." : "Buat akun merchant"}
              </Button>
            </form>

            <p className="mt-6 text-sm text-zinc-600 dark:text-zinc-400">
              Sudah punya akun?{" "}
              <Link
                href="/login"
                className="font-semibold text-amber-700 underline-offset-4 hover:underline dark:text-amber-300"
              >
                Masuk di sini
              </Link>
              .
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}
