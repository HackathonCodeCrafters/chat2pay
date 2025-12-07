"use client";

import * as React from "react";
import { useRouter } from "next/navigation";
import { Package, Truck, Clock, CheckCircle, XCircle } from "lucide-react";
import { Button, Spinner, Badge, Input } from "@/shared/components/atoms";
import { Alert, EmptyState } from "@/shared/components/molecules";
import { useMerchantAuth } from "@/features/merchant";
import { getMerchantOrders, type Order } from "@/features/checkout";
import { formatCurrency } from "@/shared/lib/utils";
import { apiClient } from "@/shared/api/client";
import { endpoints } from "@/shared/api/endpoints";

const statusConfig: Record<string, { color: string; icon: React.ElementType }> = {
  pending: { color: "bg-yellow-100 text-yellow-800", icon: Clock },
  paid: { color: "bg-blue-100 text-blue-800", icon: CheckCircle },
  processing: { color: "bg-purple-100 text-purple-800", icon: Package },
  shipped: { color: "bg-indigo-100 text-indigo-800", icon: Truck },
  delivered: { color: "bg-green-100 text-green-800", icon: CheckCircle },
  cancelled: { color: "bg-red-100 text-red-800", icon: XCircle },
};

export default function MerchantOrdersPage() {
  const router = useRouter();
  const { user, token, isAuthenticated, isLoading: authLoading } = useMerchantAuth();
  
  const [orders, setOrders] = React.useState<Order[]>([]);
  const [isLoading, setIsLoading] = React.useState(true);
  const [error, setError] = React.useState<string | null>(null);
  const [updatingId, setUpdatingId] = React.useState<string | null>(null);

  React.useEffect(() => {
    if (!authLoading && !isAuthenticated) {
      router.push("/merchant/login");
    }
  }, [authLoading, isAuthenticated, router]);

  React.useEffect(() => {
    async function loadOrders() {
      if (!token) return;
      
      setIsLoading(true);
      try {
        const data = await getMerchantOrders(token);
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

  const handleUpdateStatus = async (orderId: string, status: string, trackingNumber?: string) => {
    if (!token) return;
    
    setUpdatingId(orderId);
    try {
      await apiClient.patch(endpoints.orders.updateStatus(orderId), {
        json: { status, tracking_number: trackingNumber },
        headers: { Authorization: `Bearer ${token}` },
      });
      
      // Refresh orders
      const data = await getMerchantOrders(token);
      setOrders(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to update order");
    } finally {
      setUpdatingId(null);
    }
  };

  if (authLoading || isLoading) {
    return (
      <div className="container py-8 flex justify-center">
        <Spinner size="lg" />
      </div>
    );
  }

  return (
    <div className="container py-8">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold">Orders</h1>
          <p className="text-muted-foreground">Manage incoming orders</p>
        </div>
        <Button variant="outline" onClick={() => router.push("/merchant/dashboard")}>
          Back to Dashboard
        </Button>
      </div>

      {error && (
        <Alert variant="destructive" className="mb-6">{error}</Alert>
      )}

      {orders.length === 0 ? (
        <EmptyState
          icon={Package}
          title="No Orders Yet"
          description="When customers place orders, they will appear here."
        />
      ) : (
        <div className="space-y-4">
          {orders.map((order) => {
            const StatusIcon = statusConfig[order.status]?.icon || Clock;
            return (
              <div key={order.id} className="border rounded-lg overflow-hidden">
                {/* Header */}
                <div className="bg-muted/50 p-4 flex items-center justify-between">
                  <div>
                    <p className="text-sm text-muted-foreground">Order #{order.id.slice(0, 8)}</p>
                    <p className="text-xs text-muted-foreground">
                      {new Date(order.created_at).toLocaleDateString("id-ID", {
                        day: "numeric",
                        month: "long",
                        year: "numeric",
                        hour: "2-digit",
                        minute: "2-digit",
                      })}
                    </p>
                  </div>
                  <Badge className={statusConfig[order.status]?.color || "bg-gray-100"}>
                    <StatusIcon className="h-3 w-3 mr-1" />
                    {order.status.toUpperCase()}
                  </Badge>
                </div>

                {/* Items */}
                <div className="p-4 border-b">
                  <div className="space-y-2">
                    {order.items?.map((item) => (
                      <div key={item.id} className="flex justify-between text-sm">
                        <span>{item.product_name} x {item.quantity}</span>
                        <span>{formatCurrency(item.subtotal)}</span>
                      </div>
                    ))}
                  </div>
                </div>

                {/* Shipping */}
                <div className="p-4 border-b text-sm">
                  <p className="font-medium mb-1">Shipping to:</p>
                  <p className="text-muted-foreground">
                    {order.shipping_address}, {order.shipping_city}, {order.shipping_province}
                  </p>
                  <p className="text-muted-foreground">
                    {order.courier} - {order.courier_service} ({order.shipping_etd} days)
                  </p>
                </div>

                {/* Summary & Actions */}
                <div className="p-4 bg-muted/30 flex items-center justify-between">
                  <div>
                    <p className="text-sm text-muted-foreground">Total</p>
                    <p className="font-semibold text-primary">{formatCurrency(order.total)}</p>
                  </div>
                  <div className="flex gap-2">
                    {order.status === "pending" && (
                      <Button
                        size="sm"
                        variant="outline"
                        onClick={() => handleUpdateStatus(order.id, "processing")}
                        loading={updatingId === order.id}
                      >
                        Process Order
                      </Button>
                    )}
                    {order.status === "processing" && (
                      <div className="flex gap-2 items-center">
                        <Input
                          placeholder="Tracking number"
                          className="w-40 h-8"
                          id={`tracking-${order.id}`}
                        />
                        <Button
                          size="sm"
                          onClick={() => {
                            const input = document.getElementById(`tracking-${order.id}`) as HTMLInputElement;
                            handleUpdateStatus(order.id, "shipped", input?.value);
                          }}
                          loading={updatingId === order.id}
                        >
                          Ship Order
                        </Button>
                      </div>
                    )}
                    {order.status === "shipped" && (
                      <Button
                        size="sm"
                        variant="outline"
                        onClick={() => handleUpdateStatus(order.id, "delivered")}
                        loading={updatingId === order.id}
                      >
                        Mark Delivered
                      </Button>
                    )}
                  </div>
                </div>
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
}
