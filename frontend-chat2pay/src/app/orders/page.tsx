"use client";

import * as React from "react";
import { useRouter } from "next/navigation";
import { Package, ArrowLeft, ShoppingBag } from "lucide-react";
import { Button, Spinner, Badge } from "@/shared/components/atoms";
import { EmptyState } from "@/shared/components/molecules";
import { useAuth } from "@/features/auth";
import { getCustomerOrders, type Order } from "@/features/checkout";
import { formatCurrency } from "@/shared/lib/utils";

const statusColors: Record<string, string> = {
  pending: "bg-yellow-100 text-yellow-800",
  paid: "bg-blue-100 text-blue-800",
  processing: "bg-purple-100 text-purple-800",
  shipped: "bg-indigo-100 text-indigo-800",
  delivered: "bg-green-100 text-green-800",
  cancelled: "bg-red-100 text-red-800",
};

export default function OrdersPage() {
  const router = useRouter();
  const { token, isAuthenticated, isLoading: authLoading } = useAuth();

  const [orders, setOrders] = React.useState<Order[]>([]);
  const [isLoading, setIsLoading] = React.useState(true);
  const [error, setError] = React.useState<string | null>(null);

  React.useEffect(() => {
    if (!authLoading && !isAuthenticated) {
      router.push("/login");
    }
  }, [authLoading, isAuthenticated, router]);

  React.useEffect(() => {
    async function loadOrders() {
      if (!token) return;

      setIsLoading(true);
      try {
        const data = await getCustomerOrders(token);
        setOrders(data);
      } catch (err) {
        setError(err instanceof Error ? err.message : "Failed to load orders");
      } finally {
        setIsLoading(false);
      }
    }

    if (isAuthenticated && token) {
      loadOrders();
    }
  }, [isAuthenticated, token]);

  if (authLoading || isLoading) {
    return (
      <div className="container py-8 flex justify-center">
        <Spinner size="lg" />
      </div>
    );
  }

  return (
    <div className="container py-8 max-w-2xl">
      <Button variant="ghost" onClick={() => router.push("/chat")} className="mb-6">
        <ArrowLeft className="h-4 w-4 mr-2" />
        Back to Chat
      </Button>

      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold flex items-center gap-2">
          <ShoppingBag className="h-6 w-6" />
          Pesanan Saya
        </h1>
      </div>

      {error && (
        <div className="mb-4 p-4 bg-red-50 text-red-600 rounded-lg">{error}</div>
      )}

      {orders.length === 0 ? (
        <EmptyState
          icon={Package}
          title="Belum Ada Pesanan"
          description="Pesanan Anda akan muncul di sini setelah checkout."
        />
      ) : (
        <div className="space-y-4">
          {orders.map((order) => (
            <div
              key={order.id}
              className="border rounded-lg p-4 hover:border-primary/50 cursor-pointer transition-colors"
              onClick={() => router.push(`/orders/${order.id}`)}
            >
              <div className="flex items-start justify-between mb-3">
                <div>
                  <p className="text-sm text-muted-foreground">
                    Order #{order.id.slice(0, 8)}
                  </p>
                  <p className="text-xs text-muted-foreground">
                    {new Date(order.created_at).toLocaleDateString("id-ID", {
                      day: "numeric",
                      month: "long",
                      year: "numeric",
                    })}
                  </p>
                </div>
                <Badge className={statusColors[order.status] || "bg-gray-100"}>
                  {order.status.toUpperCase()}
                </Badge>
              </div>

              <div className="space-y-1 mb-3">
                {order.items?.slice(0, 2).map((item) => (
                  <p key={item.id} className="text-sm truncate">
                    {item.product_name} x {item.quantity}
                  </p>
                ))}
                {order.items && order.items.length > 2 && (
                  <p className="text-sm text-muted-foreground">
                    +{order.items.length - 2} item lainnya
                  </p>
                )}
              </div>

              <div className="flex items-center justify-between pt-3 border-t">
                <span className="text-sm text-muted-foreground">Total</span>
                <span className="font-semibold text-primary">
                  {formatCurrency(order.total)}
                </span>
              </div>

              {order.tracking_number && (
                <div className="mt-2 pt-2 border-t">
                  <p className="text-xs text-muted-foreground">
                    Resi: <span className="font-mono">{order.tracking_number}</span>
                  </p>
                </div>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
