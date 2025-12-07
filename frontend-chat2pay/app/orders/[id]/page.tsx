"use client";

import * as React from "react";
import { Suspense } from "react";
import { useParams, useSearchParams, useRouter } from "next/navigation";
import { CheckCircle, Package, Truck, MapPin, ArrowLeft, Search } from "lucide-react";
import { Button, Spinner, Badge } from "@/shared/components/atoms";
import { Alert } from "@/shared/components/molecules";
import { useAuth } from "@/features/auth";
import { apiClient } from "@/shared/api/client";
import type { BackendResponse } from "@/shared/api/types";
import { endpoints } from "@/shared/api/endpoints";
import { formatCurrency } from "@/shared/lib/utils";
import type { Order } from "@/features/checkout";
import { trackShipment, type TrackingResult } from "@/features/checkout";

const statusColors: Record<string, string> = {
  pending: "bg-yellow-100 text-yellow-800",
  paid: "bg-blue-100 text-blue-800",
  processing: "bg-purple-100 text-purple-800",
  shipped: "bg-indigo-100 text-indigo-800",
  delivered: "bg-green-100 text-green-800",
  cancelled: "bg-red-100 text-red-800",
};

function OrderDetailContent() {
  const params = useParams();
  const searchParams = useSearchParams();
  const router = useRouter();
  const { token, isAuthenticated } = useAuth();
  
  const [order, setOrder] = React.useState<Order | null>(null);
  const [isLoading, setIsLoading] = React.useState(true);
  const [error, setError] = React.useState<string | null>(null);
  const [tracking, setTracking] = React.useState<TrackingResult | null>(null);
  const [isTrackingLoading, setIsTrackingLoading] = React.useState(false);
  const [trackingError, setTrackingError] = React.useState<string | null>(null);
  
  const isSuccess = searchParams.get("success") === "true";
  const orderId = params.id as string;

  React.useEffect(() => {
    if (!isAuthenticated) {
      router.push("/login");
      return;
    }

    async function loadOrder() {
      if (!token || !orderId) return;
      
      setIsLoading(true);
      try {
        const response = await apiClient.get<BackendResponse<Order>>(
          endpoints.orders.byId(orderId),
          {
            headers: {
              Authorization: `Bearer ${token}`,
            },
          }
        );
        if (response.data.status && response.data.data) {
          setOrder(response.data.data);
        } else {
          setError("Order not found");
        }
      } catch {
        setError("Failed to load order");
      } finally {
        setIsLoading(false);
      }
    }
    loadOrder();
  }, [token, orderId, isAuthenticated, router]);

  if (!isAuthenticated || isLoading) {
    return (
      <div className="container py-8 flex justify-center">
        <Spinner size="lg" />
      </div>
    );
  }

  const handleTrackShipment = async () => {
    if (!order?.tracking_number || !order?.courier) return;
    
    setIsTrackingLoading(true);
    setTrackingError(null);
    
    try {
      const result = await trackShipment(order.tracking_number, order.courier.toLowerCase());
      setTracking(result);
    } catch (err) {
      setTrackingError(err instanceof Error ? err.message : "Failed to track shipment");
    } finally {
      setIsTrackingLoading(false);
    }
  };

  if (error || !order) {
    return (
      <div className="container py-8 max-w-2xl">
        <Alert variant="destructive">{error || "Order not found"}</Alert>
        <Button onClick={() => router.push("/chat")} className="mt-4">
          Back to Chat
        </Button>
      </div>
    );
  }

  return (
    <div className="container py-8 max-w-2xl">
      <Button variant="ghost" onClick={() => router.push("/chat")} className="mb-6">
        <ArrowLeft className="h-4 w-4 mr-2" />
        Back to Chat
      </Button>

      {isSuccess && (
        <Alert className="mb-6 bg-green-50 border-green-200 text-green-800">
          <CheckCircle className="h-5 w-5 inline mr-2" />
          Order placed successfully! Please complete payment to process your order.
        </Alert>
      )}

      <div className="border rounded-lg overflow-hidden">
        {/* Header */}
        <div className="bg-muted/50 p-4 border-b">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-muted-foreground">Order ID</p>
              <p className="font-mono text-sm">{order.id}</p>
            </div>
            <Badge className={statusColors[order.status] || "bg-gray-100"}>
              {order.status.toUpperCase()}
            </Badge>
          </div>
        </div>

        {/* Order Items */}
        <div className="p-4 border-b">
          <h3 className="font-semibold flex items-center gap-2 mb-3">
            <Package className="h-4 w-4" />
            Items
          </h3>
          <div className="space-y-2">
            {order.items?.map((item) => (
              <div key={item.id} className="flex justify-between text-sm">
                <span>
                  {item.product_name} x {item.quantity}
                </span>
                <span>{formatCurrency(item.subtotal)}</span>
              </div>
            ))}
          </div>
        </div>

        {/* Shipping */}
        <div className="p-4 border-b">
          <h3 className="font-semibold flex items-center gap-2 mb-3">
            <Truck className="h-4 w-4" />
            Shipping
          </h3>
          <div className="text-sm space-y-1">
            <p>
              <span className="text-muted-foreground">Courier:</span>{" "}
              {order.courier} - {order.courier_service}
            </p>
            {order.shipping_etd && (
              <p>
                <span className="text-muted-foreground">Estimated:</span>{" "}
                {order.shipping_etd} days
              </p>
            )}
            {order.tracking_number && (
              <div className="space-y-2">
                <p>
                  <span className="text-muted-foreground">Tracking:</span>{" "}
                  <span className="font-mono">{order.tracking_number}</span>
                </p>
                <Button
                  size="sm"
                  variant="outline"
                  onClick={handleTrackShipment}
                  loading={isTrackingLoading}
                >
                  <Search className="h-4 w-4 mr-1" />
                  Track Shipment
                </Button>
                {trackingError && (
                  <p className="text-sm text-red-500">{trackingError}</p>
                )}
              </div>
            )}
          </div>
        </div>

        {/* Address */}
        <div className="p-4 border-b">
          <h3 className="font-semibold flex items-center gap-2 mb-3">
            <MapPin className="h-4 w-4" />
            Delivery Address
          </h3>
          <p className="text-sm">
            {order.shipping_address}
            <br />
            {order.shipping_city}, {order.shipping_province} {order.shipping_postal_code}
          </p>
        </div>

        {/* Summary */}
        <div className="p-4 bg-muted/30">
          <div className="space-y-2 text-sm">
            <div className="flex justify-between">
              <span>Subtotal</span>
              <span>{formatCurrency(order.subtotal)}</span>
            </div>
            <div className="flex justify-between">
              <span>Shipping</span>
              <span>{formatCurrency(order.shipping_cost)}</span>
            </div>
            <div className="flex justify-between font-semibold text-base pt-2 border-t">
              <span>Total</span>
              <span className="text-primary">{formatCurrency(order.total)}</span>
            </div>
          </div>

          {order.payment_status === "pending" && (
            <div className="mt-4 p-3 bg-yellow-50 border border-yellow-200 rounded-lg">
              <p className="text-sm text-yellow-800">
                Payment pending. Please transfer <strong>{formatCurrency(order.total)}</strong> to complete your order.
              </p>
            </div>
          )}
        </div>
      </div>

      {order.notes && (
        <div className="mt-4 p-4 border rounded-lg">
          <h3 className="font-semibold mb-2">Notes</h3>
          <p className="text-sm text-muted-foreground">{order.notes}</p>
        </div>
      )}

      {/* Tracking Info */}
      {tracking && (
        <div className="mt-4 border rounded-lg overflow-hidden">
          <div className="bg-muted/50 p-4 border-b">
            <h3 className="font-semibold flex items-center gap-2">
              <Truck className="h-4 w-4" />
              Tracking Information
            </h3>
            <p className="text-sm text-muted-foreground mt-1">
              {tracking.summary.courier_name} - {tracking.summary.service_code}
            </p>
          </div>
          
          <div className="p-4 border-b">
            <div className="flex items-center justify-between">
              <span className="text-sm">Status</span>
              <Badge className={tracking.delivered ? "bg-green-100 text-green-800" : "bg-blue-100 text-blue-800"}>
                {tracking.delivered ? "Delivered" : tracking.summary.status}
              </Badge>
            </div>
            {tracking.delivery_status.pod_receiver && (
              <p className="text-sm text-muted-foreground mt-2">
                Received by: {tracking.delivery_status.pod_receiver}
              </p>
            )}
          </div>

          <div className="p-4">
            <h4 className="font-medium text-sm mb-3">Shipment History</h4>
            <div className="space-y-3">
              {tracking.manifest.map((item, index) => (
                <div key={index} className="flex gap-3 text-sm">
                  <div className="flex flex-col items-center">
                    <div className={`w-2 h-2 rounded-full ${index === 0 ? "bg-primary" : "bg-gray-300"}`} />
                    {index < tracking.manifest.length - 1 && (
                      <div className="w-px h-full bg-gray-200 my-1" />
                    )}
                  </div>
                  <div className="flex-1 pb-3">
                    <p className="font-medium">{item.manifest_description}</p>
                    <p className="text-muted-foreground text-xs">
                      {item.manifest_date} {item.manifest_time} - {item.city_name}
                    </p>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

export default function OrderDetailPage() {
  return (
    <Suspense fallback={<div className="container py-8 flex justify-center"><Spinner size="lg" /></div>}>
      <OrderDetailContent />
    </Suspense>
  );
}
