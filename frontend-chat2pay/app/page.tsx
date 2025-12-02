import { ShieldCheck, Signal, Workflow } from "lucide-react";

import { Button } from "@/components/ui/button";
import { getRuntimeConfig } from "@/lib/config/runtime";

const codeSample = `import { apiClient, endpoints } from "@/lib/api";

type Payment = { id: string; amount: number; status: string };

export async function getPayments() {
  const { data } = await apiClient.get<Payment[]>(
    endpoints.payments.root(),
    { query: { page: 1, per_page: 10 } }
  );

  return data;
}

export async function createPayment(payload: {
  amount: number;
  currency: string;
}) {
  const { data } = await apiClient.post<Payment>(
    endpoints.payments.root(),
    { json: payload }
  );

  return data;
}`;

const checklist = [
  "Set NEXT_PUBLIC_API_BASE_URL di .env lokal/hosting.",
  "Gunakan apiClient.get/post/put/patch/delete dengan penanganan error terpadu.",
  "Definisikan path di lib/api/endpoints.ts supaya konsisten.",
  "Suntikkan token dengan opsi getAuthToken bila API butuh auth.",
];

export default function Home() {
  const config = getRuntimeConfig();
  const baseUrl = config.apiBaseUrl || "Belum di-set";

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-50 via-white to-slate-100">
      <main className="mx-auto flex w-full max-w-5xl flex-col gap-10 px-6 py-16">
        <section className="flex flex-col gap-3">
          <div className="inline-flex items-center gap-3 rounded-full bg-slate-900 text-white px-4 py-2 text-xs font-semibold uppercase tracking-[0.12em]">
            <Signal className="size-4" />
            API Integration Ready
          </div>
          <div className="flex flex-col gap-4 md:flex-row md:items-end md:justify-between">
            <div className="space-y-2">
              <h1 className="text-4xl font-semibold leading-tight text-slate-900">
                Infrastruktur integrasi API siap dipakai
              </h1>
              <p className="max-w-2xl text-lg text-slate-600">
                Base URL terkelola dari environment, client HTTP dengan
                penanganan error yang seragam, dan registry endpoint yang
                rapi—siap disambungkan ke layanan pembayaran atau chatbot Anda.
              </p>
            </div>
            <Button className="self-start md:self-auto" size="lg">
              Lihat struktur API
            </Button>
          </div>
          <div className="flex flex-wrap gap-3 text-sm text-slate-600">
            <span className="rounded-full bg-slate-200 px-3 py-1">
              Base URL: {baseUrl}
            </span>
            <span className="rounded-full bg-slate-200 px-3 py-1">
              Lib: lib/api/*
            </span>
            <span className="rounded-full bg-slate-200 px-3 py-1">
              Error class: ApiError
            </span>
          </div>
        </section>

        <section className="grid gap-6 md:grid-cols-3">
          <article className="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm">
            <div className="flex items-center gap-3 text-slate-900">
              <Workflow className="size-5" />
              <h2 className="text-lg font-semibold">Alur standar</h2>
            </div>
            <p className="mt-3 text-sm text-slate-600">
              Request builder dengan query params, JSON serializer, dan helper
              HTTP verb (GET/POST/PUT/PATCH/DELETE) untuk menjaga konsistensi.
            </p>
          </article>

          <article className="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm">
            <div className="flex items-center gap-3 text-slate-900">
              <ShieldCheck className="size-5" />
              <h2 className="text-lg font-semibold">Penanganan error</h2>
            </div>
            <p className="mt-3 text-sm text-slate-600">
              ApiError membawa status code, payload error (jika ada), dan
              request-id untuk tracing, sehingga mudah di-observasi.
            </p>
          </article>

          <article className="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm">
            <div className="flex items-center gap-3 text-slate-900">
              <Signal className="size-5" />
              <h2 className="text-lg font-semibold">Endpoint registry</h2>
            </div>
            <p className="mt-3 text-sm text-slate-600">
              Semua path API terpusat di <code>lib/api/endpoints.ts</code> agar
              mudah dipakai ulang dan meminimalkan typo.
            </p>
          </article>
        </section>

        <section className="grid gap-6 lg:grid-cols-[1.2fr_0.8fr]">
          <article className="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm">
            <div className="flex items-center justify-between gap-3">
              <div>
                <h3 className="text-lg font-semibold text-slate-900">
                  Contoh pemanggilan API
                </h3>
                <p className="text-sm text-slate-600">
                  Cukup pakai helper bawaan, tanpa perlu menulis boilerplate
                  fetch.
                </p>
              </div>
              <span className="rounded-full bg-slate-900 px-3 py-1 text-xs font-medium text-white">
                Type-safe
              </span>
            </div>
            <pre className="mt-4 overflow-x-auto rounded-xl bg-slate-900 p-5 text-xs text-slate-100">
              <code>{codeSample}</code>
            </pre>
          </article>

          <article className="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm">
            <div className="flex items-center gap-3 text-slate-900">
              <ShieldCheck className="size-5" />
              <h3 className="text-lg font-semibold">Checklist implementasi</h3>
            </div>
            <ul className="mt-4 space-y-3 text-sm text-slate-600">
              {checklist.map((item) => (
                <li
                  key={item}
                  className="flex items-start gap-2 rounded-xl bg-slate-50 px-3 py-2"
                >
                  <span className="mt-0.5 size-2 rounded-full bg-emerald-500" />
                  <span>{item}</span>
                </li>
              ))}
            </ul>
          </article>
        </section>
      </main>
    </div>
  );
}
