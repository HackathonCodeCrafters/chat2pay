"use client";

import { useState, type FormEvent } from "react";
import { ArrowLeft, Pencil, Plus, RefreshCcw, Shield, Trash2 } from "lucide-react";
import Link from "next/link";

import { Button } from "@/components/ui/button";
import { ApiError, apiClient, endpoints } from "@/lib/api";
import { cn } from "@/lib/utils";

type Order = {
  id: string;
  customer_id?: string;
  amount?: number;
  status?: string;
  created_at?: string;
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

export default function OrderManagementPage() {
  const [orders, setOrders] = useState<Order[]>([]);
  const [selectedId, setSelectedId] = useState<string | null>(null);

  const [tableStatus, setTableStatus] = useState<Status>({ type: "idle" });
  const [formStatus, setFormStatus] = useState<Status>({ type: "idle" });
  const [rowStatus, setRowStatus] = useState<Status>({ type: "idle" });

  const [form, setForm] = useState({
    id: "",
    customer_id: "",
    amount: "",
    status: "",
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
      const { data } = await apiClient.get<Order[]>(endpoints.orders.root(), {
        next: { revalidate: 0 },
      });
      setOrders(data);
      setTableStatus({
        type: "success",
        message: `Berhasil memuat ${data.length} order.`,
      });
    } catch (error) {
      setTableStatus({
        type: "error",
        message: parseError(error, "Gagal memuat daftar order."),
      });
    }
  };

  const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setFormStatus({ type: "loading" });
    setRowStatus({ type: "idle" });

    const payload = {
      customer_id: form.customer_id || undefined,
      amount: form.amount ? Number(form.amount) : undefined,
      status: form.status || undefined,
    };

    if (isEditing) {
      try {
        const { data } = await apiClient.put<Order>(
          endpoints.orders.byId(selectedId!),
          { json: payload }
        );
        setOrders((prev) =>
          prev.map((order) =>
            order.id === selectedId ? { ...order, ...data } : order
          )
        );
        setFormStatus({
          type: "success",
          message: `Order ${selectedId} berhasil diperbarui.`,
        });
      } catch (error) {
        setFormStatus({
          type: "error",
          message: parseError(error, "Gagal memperbarui order."),
        });
      }
      return;
    }

    try {
      const { data } = await apiClient.post<Order>(endpoints.orders.root(), {
        json: payload,
      });
      setOrders((prev) => [data, ...prev]);
      setFormStatus({
        type: "success",
        message: "Order baru berhasil dibuat.",
      });
      setForm({ id: "", customer_id: "", amount: "", status: "" });
    } catch (error) {
      setFormStatus({
        type: "error",
        message: parseError(error, "Gagal membuat order."),
      });
    }
  };

  const handleSelect = (order: Order) => {
    setSelectedId(order.id);
    setForm({
      id: order.id,
      customer_id: order.customer_id ?? "",
      amount: order.amount ? String(order.amount) : "",
      status: order.status ?? "",
    });
    setFormStatus({ type: "idle" });
    setRowStatus({ type: "idle" });
  };

  const handleCancelEdit = () => {
    setSelectedId(null);
    setForm({ id: "", customer_id: "", amount: "", status: "" });
    setFormStatus({ type: "idle" });
  };

  const handleDelete = async (orderId: string) => {
    setRowStatus({ type: "loading" });
    try {
      await apiClient.delete<void>(endpoints.orders.byId(orderId));
      setOrders((prev) => prev.filter((order) => order.id !== orderId));
      if (selectedId === orderId) handleCancelEdit();
      setRowStatus({
        type: "success",
        message: `Order ${orderId} berhasil dihapus.`,
      });
    } catch (error) {
      setRowStatus({
        type: "error",
        message: parseError(error, "Gagal menghapus order."),
      });
    }
  };

  const handleRefreshRow = async (orderId: string) => {
    setRowStatus({ type: "loading" });
    try {
      const { data } = await apiClient.get<Order>(
        endpoints.orders.byId(orderId),
        { next: { revalidate: 0 } }
      );
      setOrders((prev) =>
        prev.map((order) =>
          order.id === orderId ? { ...order, ...data } : order
        )
      );
      setRowStatus({
        type: "success",
        message: `Detail order ${orderId} diperbarui.`,
      });
    } catch (error) {
      setRowStatus({
        type: "error",
        message: parseError(error, "Gagal mengambil detail order."),
      });
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-indigo-50 via-white to-amber-50 text-zinc-900 dark:from-black dark:via-zinc-950 dark:to-zinc-900 dark:text-zinc-100">
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
            <Shield className="h-4 w-4" />
            Order Management
          </div>
          <div className="space-y-2">
            <h1 className="text-4xl font-semibold leading-tight">Tabel Order</h1>
            <p className="max-w-3xl text-base text-zinc-600 dark:text-zinc-300">
              Kelola pesanan: lihat, tambah, ubah, dan hapus tanpa token.
            </p>
          </div>
        </header>

        <div className="grid gap-6 lg:grid-cols-[2fr_1fr]">
          <div className={panelClass}>
            <div className="mb-4 flex flex-wrap items-center justify-between gap-3">
              <div className="space-y-1">
                <p className="text-sm font-semibold text-amber-700 dark:text-amber-300">
                  Data order
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
              <div className="grid grid-cols-[1.3fr_1.1fr_1fr_0.9fr_0.9fr] gap-3 bg-zinc-100 px-4 py-3 text-xs font-semibold uppercase tracking-[0.08em] text-zinc-600 dark:bg-zinc-900 dark:text-zinc-300">
                <span>ID Order</span>
                <span>ID Customer</span>
                <span>Nominal</span>
                <span>Status</span>
                <span className="text-right">Aksi</span>
              </div>
              {orders.length === 0 ? (
                <div className="px-4 py-6 text-sm text-zinc-500 dark:text-zinc-400">
                  Belum ada data. Gunakan tombol Refresh atau buat order baru.
                </div>
              ) : (
                <div className="divide-y divide-zinc-200 dark:divide-zinc-800">
                  {orders.map((order) => (
                    <div
                      key={order.id}
                      className={cn(
                        "grid grid-cols-[1.3fr_1.1fr_1fr_0.9fr_0.9fr] items-center gap-3 px-4 py-3 text-sm",
                        selectedId === order.id
                          ? "bg-amber-50/80 dark:bg-amber-500/10"
                          : "bg-white dark:bg-zinc-900"
                      )}
                    >
                      <div className="truncate font-medium text-zinc-900 dark:text-zinc-100">
                        {order.id}
                        <p className="text-xs text-zinc-500 dark:text-zinc-400">
                          {order.created_at ? new Date(order.created_at).toLocaleString("id-ID") : "—"}
                        </p>
                      </div>
                      <span className="truncate text-zinc-700 dark:text-zinc-200">
                        {order.customer_id || "—"}
                      </span>
                      <span className="truncate text-zinc-700 dark:text-zinc-200">
                        {order.amount !== undefined ? `Rp ${order.amount.toLocaleString("id-ID")}` : "—"}
                      </span>
                      <span className="truncate text-zinc-600 dark:text-zinc-300">
                        {order.status || "—"}
                      </span>
                      <div className="flex items-center justify-end gap-2 text-xs">
                        <Button
                          size="sm"
                          variant="ghost"
                          onClick={() => handleRefreshRow(order.id)}
                          title="GET /orders/:id"
                        >
                          <RefreshCcw className="h-4 w-4" />
                        </Button>
                        <Button
                          size="sm"
                          variant="ghost"
                          onClick={() => handleSelect(order)}
                          title="Edit order"
                        >
                          <Pencil className="h-4 w-4" />
                        </Button>
                        <Button
                          size="sm"
                          variant="ghost"
                          className="text-rose-600 hover:text-rose-700 dark:text-rose-300"
                          onClick={() => handleDelete(order.id)}
                          title="DELETE /orders/:id"
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
                  {isEditing ? "Edit order" : "Tambah order"}
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
                  ID Customer (opsional)
                </label>
                <input
                  className={inputClass}
                  placeholder="cust_123"
                  value={form.customer_id}
                  onChange={(event) =>
                    setForm((prev) => ({ ...prev, customer_id: event.target.value }))
                  }
                />
              </div>
              <div className="grid gap-3 sm:grid-cols-2">
                <div className="space-y-1">
                  <label className="text-sm font-medium text-zinc-700 dark:text-zinc-200">
                    Nominal (angka)
                  </label>
                  <input
                    className={inputClass}
                    type="number"
                    min="0"
                    step="0.01"
                    placeholder="250000"
                    value={form.amount}
                    onChange={(event) =>
                      setForm((prev) => ({ ...prev, amount: event.target.value }))
                    }
                  />
                </div>
                <div className="space-y-1">
                  <label className="text-sm font-medium text-zinc-700 dark:text-zinc-200">
                    Status (opsional)
                  </label>
                  <input
                    className={inputClass}
                    placeholder="pending / paid / failed"
                    value={form.status}
                    onChange={(event) =>
                      setForm((prev) => ({ ...prev, status: event.target.value }))
                    }
                  />
                </div>
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
                      : "Buat order"}
                </Button>
                {!isEditing ? (
                  <Button
                    type="button"
                    variant="outline"
                    onClick={() => setForm({ id: "", customer_id: "", amount: "", status: "" })}
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
