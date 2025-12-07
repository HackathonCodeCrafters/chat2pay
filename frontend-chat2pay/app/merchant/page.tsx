"use client";

import { useMemo, useState, type FormEvent } from "react";
import {
  ArrowLeft,
  Pencil,
  Plus,
  RefreshCcw,
  ShieldCheck,
  Trash2,
} from "lucide-react";
import Link from "next/link";

import { Button } from "@/components/ui/button";
import { ApiError, apiClient, endpoints } from "@/lib/api";
import { cn } from "@/lib/utils";

type Merchant = {
  id: string;
  name: string;
  email: string;
  phone?: string;
  status?: string;
};

type Status =
  | { type: "idle" }
  | { type: "loading" }
  | { type: "success"; message: string }
  | { type: "error"; message: string };

const inputClass =
  "w-full rounded-lg border border-zinc-200 bg-white px-3 py-2.5 text-sm shadow-sm outline-none ring-0 transition focus:border-zinc-800 focus:ring-2 focus:ring-zinc-900/10 dark:border-zinc-700 dark:bg-zinc-900 dark:text-zinc-50 dark:focus:border-zinc-100 dark:focus:ring-zinc-100/10";

const panelClass =
  "rounded-2xl border border-zinc-200/80 bg-white/80 p-6 shadow-xl backdrop-blur dark:border-zinc-800/80 dark:bg-zinc-900/80";

const parseError = (error: unknown, fallback: string) => {
  if (error instanceof ApiError) {
    if (error.payload && typeof error.payload === "object") {
      const message = (error.payload as { message?: string }).message;
      if (message) return message;
    }
    if (error.message) return error.message;
  }
  return fallback;
};

export default function MerchantTablePage() {
  const [merchants, setMerchants] = useState<Merchant[]>([]);
  const [selectedId, setSelectedId] = useState<string | null>(null);

  const [tableStatus, setTableStatus] = useState<Status>({ type: "idle" });
  const [formStatus, setFormStatus] = useState<Status>({ type: "idle" });
  const [rowStatus, setRowStatus] = useState<Status>({ type: "idle" });

  const [form, setForm] = useState({
    name: "",
    email: "",
    phone: "",
    id: "",
  });

  const isEditing = Boolean(selectedId);

  const badgeMessage = (status: Status) => {
    if (status.type === "loading") return "Sedang memproses...";
    if (status.type === "success") return status.message;
    if (status.type === "error") return status.message;
    return null;
  };

  const fetchAll = async () => {
    setTableStatus({ type: "loading" });
    try {
      const { data } = await apiClient.get<Merchant[]>(
        endpoints.merchants.root(),
        { next: { revalidate: 0 } }
      );
      setMerchants(data);
      setTableStatus({
        type: "success",
        message: `Berhasil memuat ${data.length} merchant.`,
      });
    } catch (error) {
      setTableStatus({
        type: "error",
        message: parseError(error, "Gagal memuat daftar merchant."),
      });
    }
  };

  const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setFormStatus({ type: "loading" });
    setRowStatus({ type: "idle" });

    if (isEditing) {
      try {
        const { data } = await apiClient.put<Merchant>(
          endpoints.merchants.byId(selectedId!),
          {
            json: {
              name: form.name || undefined,
              email: form.email || undefined,
              phone: form.phone || undefined,
            },
          }
        );

        setMerchants((prev) =>
          prev.map((merchant) =>
            merchant.id === selectedId ? { ...merchant, ...data } : merchant
          )
        );
        setFormStatus({
          type: "success",
          message: `Merchant ${selectedId} berhasil diperbarui.`,
        });
      } catch (error) {
        setFormStatus({
          type: "error",
          message: parseError(error, "Gagal memperbarui merchant."),
        });
      }
      return;
    }

    try {
      const { data } = await apiClient.post<Merchant>(
        endpoints.merchants.root(),
        {
          json: { name: form.name, email: form.email, phone: form.phone },
        }
      );
      setMerchants((prev) => [data, ...prev]);
      setFormStatus({
        type: "success",
        message: "Merchant baru berhasil dibuat.",
      });
      setForm({ id: "", name: "", email: "", phone: "" });
    } catch (error) {
      setFormStatus({
        type: "error",
        message: parseError(error, "Gagal membuat merchant."),
      });
    }
  };

  const handleSelect = (merchant: Merchant) => {
    setSelectedId(merchant.id);
    setForm({
      id: merchant.id,
      name: merchant.name ?? "",
      email: merchant.email ?? "",
      phone: merchant.phone ?? "",
    });
    setFormStatus({ type: "idle" });
    setRowStatus({ type: "idle" });
  };

  const handleCancelEdit = () => {
    setSelectedId(null);
    setForm({ id: "", name: "", email: "", phone: "" });
    setFormStatus({ type: "idle" });
  };

  const handleDelete = async (merchantId: string) => {
    setRowStatus({ type: "loading" });
    try {
      await apiClient.delete<void>(endpoints.merchants.byId(merchantId), {
      });
      setMerchants((prev) => prev.filter((merchant) => merchant.id !== merchantId));
      if (selectedId === merchantId) handleCancelEdit();
      setRowStatus({
        type: "success",
        message: `Merchant ${merchantId} berhasil dihapus.`,
      });
    } catch (error) {
      setRowStatus({
        type: "error",
        message: parseError(error, "Gagal menghapus merchant."),
      });
    }
  };

  const handleRefreshRow = async (merchantId: string) => {
    setRowStatus({ type: "loading" });
    try {
      const { data } = await apiClient.get<Merchant>(
        endpoints.merchants.byId(merchantId),
        { next: { revalidate: 0 } }
      );
      setMerchants((prev) =>
        prev.map((merchant) =>
          merchant.id === merchantId ? { ...merchant, ...data } : merchant
        )
      );
      setRowStatus({
        type: "success",
        message: `Detail merchant ${merchantId} diperbarui.`,
      });
    } catch (error) {
      setRowStatus({
        type: "error",
        message: parseError(error, "Gagal mengambil detail merchant."),
      });
    }
  };

  const quickStats = useMemo(() => {
    if (!merchants.length) return null;
    const active = merchants.filter((item) => item.status === "active").length;
    return { total: merchants.length, active };
  }, [merchants]);

  return (
    <div className="min-h-screen bg-gradient-to-br from-amber-50 via-white to-slate-50 text-zinc-900 dark:from-black dark:via-zinc-950 dark:to-zinc-900 dark:text-zinc-100">
      <div className="mx-auto flex max-w-6xl flex-col gap-10 px-6 py-12">
        <div className="flex items-center gap-3 text-sm text-amber-700 dark:text-amber-300">
          <ArrowLeft className="h-4 w-4" />
          <Link
            href="/"
            className="font-semibold underline-offset-4 hover:underline"
          >
            Kembali ke beranda
          </Link>
        </div>

        <header className="space-y-4">
          <div className="inline-flex items-center gap-2 rounded-full bg-amber-100 px-3 py-1 text-xs font-semibold uppercase tracking-[0.16em] text-amber-700 shadow-sm dark:bg-amber-500/15 dark:text-amber-200">
            <ShieldCheck className="h-4 w-4" />
            Merchant Management
          </div>
          <div className="space-y-2">
            <h1 className="text-4xl font-semibold leading-tight">
              Tabel Merchant untuk Admin
            </h1>
            <p className="max-w-3xl text-base text-zinc-600 dark:text-zinc-300">
              Panel admin untuk melihat dan memelihara data merchant tanpa token.
            </p>
          </div>
          {quickStats ? (
            <p className="text-sm text-zinc-600 dark:text-zinc-400">
              Total merchant: <strong>{quickStats.total}</strong>{" "}
              {quickStats.active ? (
                <span className="text-emerald-600 dark:text-emerald-400">
                  ({quickStats.active} active)
                </span>
              ) : null}
            </p>
          ) : null}
        </header>

        <div className="grid gap-6 lg:grid-cols-[2fr_1fr]">
          <div className={panelClass}>
            <div className="mb-4 flex flex-wrap items-center justify-between gap-3">
              <div className="space-y-1">
                <p className="text-sm font-semibold text-amber-700 dark:text-amber-300">
                  Data merchant
                </p>
                <h3 className="text-xl font-semibold">Daftar & aksi cepat</h3>
              </div>
              <div className="flex gap-2">
                <Button variant="outline" size="sm" onClick={fetchAll} disabled={tableStatus.type === "loading"}>
                  <RefreshCcw className="h-4 w-4" />
                  Refresh
                </Button>
                <Button
                  variant="secondary"
                  size="sm"
                  onClick={() => setSelectedId(null)}
                  className="hidden sm:inline-flex"
                >
                  <Plus className="h-4 w-4" />
                  Mode tambah
                </Button>
              </div>
            </div>

            {badgeMessage(tableStatus) ? (
              <p
                className={cn(
                  "mb-3 text-sm",
                  tableStatus.type === "error"
                    ? "text-rose-500"
                    : "text-emerald-600 dark:text-emerald-400"
                )}
              >
                {badgeMessage(tableStatus)}
              </p>
            ) : null}
            {badgeMessage(rowStatus) ? (
              <p
                className={cn(
                  "mb-3 text-sm",
                  rowStatus.type === "error"
                    ? "text-rose-500"
                    : "text-emerald-600 dark:text-emerald-400"
                )}
              >
                {badgeMessage(rowStatus)}
              </p>
            ) : null}

            <div className="overflow-hidden rounded-xl border border-zinc-200/80 shadow-sm dark:border-zinc-800/80">
              <div className="grid grid-cols-[1.5fr_1.5fr_1fr_0.8fr] gap-3 bg-zinc-100 px-4 py-3 text-xs font-semibold uppercase tracking-[0.08em] text-zinc-600 dark:bg-zinc-900 dark:text-zinc-300">
                <span>Nama</span>
                <span>Email</span>
                <span>Status</span>
                <span className="text-right">Aksi</span>
              </div>
              {merchants.length === 0 ? (
                <div className="px-4 py-6 text-sm text-zinc-500 dark:text-zinc-400">
                  Belum ada data. Gunakan tombol Refresh atau buat merchant baru.
                </div>
              ) : (
                <div className="divide-y divide-zinc-200 dark:divide-zinc-800">
                  {merchants.map((merchant) => (
                    <div
                      key={merchant.id}
                      className={cn(
                        "grid grid-cols-[1.5fr_1.5fr_1fr_0.8fr] items-center gap-3 px-4 py-3 text-sm",
                        selectedId === merchant.id
                          ? "bg-amber-50/80 dark:bg-amber-500/10"
                          : "bg-white dark:bg-zinc-900"
                      )}
                    >
                      <div className="truncate font-medium text-zinc-900 dark:text-zinc-100">
                        {merchant.name || "—"}
                        <p className="text-xs text-zinc-500 dark:text-zinc-400">
                          {merchant.id}
                        </p>
                      </div>
                      <span className="truncate text-zinc-700 dark:text-zinc-200">
                        {merchant.email || "—"}
                      </span>
                      <span className="truncate text-zinc-600 dark:text-zinc-300">
                        {merchant.status || "—"}
                      </span>
                      <div className="flex items-center justify-end gap-2 text-xs">
                        <Button
                          size="sm"
                          variant="ghost"
                          onClick={() => handleRefreshRow(merchant.id)}
                          title="GET /merchants/:id"
                        >
                          <RefreshCcw className="h-4 w-4" />
                        </Button>
                        <Button
                          size="sm"
                          variant="ghost"
                          onClick={() => handleSelect(merchant)}
                          title="Edit merchant"
                        >
                          <Pencil className="h-4 w-4" />
                        </Button>
                        <Button
                          size="sm"
                          variant="ghost"
                          className="text-rose-600 hover:text-rose-700 dark:text-rose-300"
                          onClick={() => handleDelete(merchant.id)}
                          title="DELETE /merchants/:id"
                        >
                          <Trash2 className="h-4 w-4" />
                        </Button>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>

          <div className={panelClass}>
            <div className="mb-4 flex items-center justify-between gap-3">
              <div>
                <p className="text-sm font-semibold text-amber-700 dark:text-amber-300">
                  {isEditing ? "Edit merchant" : "Tambah merchant"}
                </p>
                <h3 className="text-xl font-semibold">
                  {isEditing ? "Update data" : "Formulir baru"}
                </h3>
              </div>
              {isEditing ? (
                <Button variant="outline" size="sm" onClick={handleCancelEdit}>
                  Batal edit
                </Button>
              ) : null}
            </div>

            <form className="space-y-3" onSubmit={handleSubmit}>
              {isEditing ? (
                <div className="space-y-1">
                  <label className="text-xs font-medium text-zinc-500 dark:text-zinc-400">
                    ID (readonly)
                  </label>
                  <input className={cn(inputClass, "bg-zinc-100 dark:bg-zinc-800")} value={form.id} readOnly />
                </div>
              ) : null}
              <div className="space-y-1">
                <label className="text-sm font-medium text-zinc-700 dark:text-zinc-200">
                  Nama merchant
                </label>
                <input
                  className={inputClass}
                  placeholder="Nama"
                  required={!isEditing}
                  value={form.name}
                  onChange={(event) =>
                    setForm((prev) => ({ ...prev, name: event.target.value }))
                  }
                />
              </div>
              <div className="space-y-1">
                <label className="text-sm font-medium text-zinc-700 dark:text-zinc-200">
                  Email
                </label>
                <input
                  className={inputClass}
                  type="email"
                  placeholder="ops@merchant.id"
                  required={!isEditing}
                  value={form.email}
                  onChange={(event) =>
                    setForm((prev) => ({ ...prev, email: event.target.value }))
                  }
                />
              </div>
              <div className="space-y-1">
                <label className="text-sm font-medium text-zinc-700 dark:text-zinc-200">
                  Telepon (opsional)
                </label>
                <input
                  className={inputClass}
                  placeholder="+62..."
                  value={form.phone}
                  onChange={(event) =>
                    setForm((prev) => ({ ...prev, phone: event.target.value }))
                  }
                />
              </div>

              {badgeMessage(formStatus) ? (
                <p
                  className={cn(
                    "text-sm",
                    formStatus.type === "error"
                      ? "text-rose-500"
                      : "text-emerald-600 dark:text-emerald-400"
                  )}
                >
                  {badgeMessage(formStatus)}
                </p>
              ) : null}

              <div className="flex items-center gap-3">
                <Button type="submit" disabled={formStatus.type === "loading"}>
                  {formStatus.type === "loading"
                    ? "Memproses..."
                    : isEditing
                      ? "Simpan perubahan"
                      : "Buat merchant"}
                </Button>
                {!isEditing ? (
                  <Button
                    type="button"
                    variant="outline"
                    onClick={() => setForm({ id: "", name: "", email: "", phone: "" })}
                  >
                    Reset form
                  </Button>
                ) : null}
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>
  );
}
