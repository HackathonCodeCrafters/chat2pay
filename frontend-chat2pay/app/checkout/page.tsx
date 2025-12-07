"use client";

import * as React from "react";
import { Suspense } from "react";
import { useRouter } from "next/navigation";
import { ArrowLeft, MapPin, Truck, CreditCard, Package } from "lucide-react";
import { Button, Input, Label, Spinner } from "@/shared/components/atoms";
import { Alert } from "@/shared/components/molecules";
import { useAuth } from "@/features/auth";
import { useCart } from "@/features/cart";
import { 
  getProvinces, 
  getCities, 
  getShippingCost, 
  createOrder,
  type Province,
  type City,
  type CourierResult,
} from "@/features/checkout";
import { formatCurrency } from "@/shared/lib/utils";

function CheckoutContent() {
  const router = useRouter();
  const { user, token, isAuthenticated } = useAuth();
  const { items: cartItems, isLoaded: isCartLoaded, clearCart } = useCart();

  const [provinces, setProvinces] = React.useState<Province[]>([]);
  const [cities, setCities] = React.useState<City[]>([]);
  const [shippingOptions, setShippingOptions] = React.useState<CourierResult[]>([]);
  
  const [selectedProvince, setSelectedProvince] = React.useState("");
  const [selectedCity, setSelectedCity] = React.useState("");
  const [address, setAddress] = React.useState("");
  const [postalCode, setPostalCode] = React.useState("");
  const [selectedCourier, setSelectedCourier] = React.useState("");
  const [selectedService, setSelectedService] = React.useState("");
  const [notes, setNotes] = React.useState("");
  
  const [isLoadingProvinces, setIsLoadingProvinces] = React.useState(false);
  const [isLoadingCities, setIsLoadingCities] = React.useState(false);
  const [isLoadingShipping, setIsLoadingShipping] = React.useState(false);
  const [isSubmitting, setIsSubmitting] = React.useState(false);
  const [error, setError] = React.useState<string | null>(null);

  // Redirect if cart is empty
  React.useEffect(() => {
    if (isCartLoaded && cartItems.length === 0) {
      router.push("/cart");
    }
  }, [isCartLoaded, cartItems.length, router]);

  // Redirect if not authenticated
  React.useEffect(() => {
    if (!isAuthenticated) {
      router.push("/login");
    }
  }, [isAuthenticated, router]);

  // Load provinces on mount
  React.useEffect(() => {
    async function loadProvinces() {
      setIsLoadingProvinces(true);
      try {
        const data = await getProvinces();
        setProvinces(data);
      } catch (err) {
        setError("Failed to load provinces");
      } finally {
        setIsLoadingProvinces(false);
      }
    }
    loadProvinces();
  }, []);

  // Load cities when province changes
  React.useEffect(() => {
    if (!selectedProvince) {
      setCities([]);
      return;
    }

    async function loadCities() {
      setIsLoadingCities(true);
      try {
        const data = await getCities(selectedProvince);
        setCities(data);
      } catch (err) {
        setError("Failed to load cities");
      } finally {
        setIsLoadingCities(false);
      }
    }
    loadCities();
  }, [selectedProvince]);

  // Calculate total weight
  const totalWeight = cartItems.reduce((sum, item) => sum + (item.product.weight || 1000) * item.quantity, 0);

  // Load shipping cost when city is selected
  React.useEffect(() => {
    if (!selectedCity || cartItems.length === 0) {
      setShippingOptions([]);
      return;
    }

    async function loadShippingCost() {
      setIsLoadingShipping(true);
      try {
        // Assume merchant origin is Jakarta (city_id: 152) - TODO: get from merchant data
        const data = await getShippingCost("152", selectedCity, totalWeight);
        setShippingOptions(data);
      } catch (err) {
        setError("Failed to load shipping options");
      } finally {
        setIsLoadingShipping(false);
      }
    }
    loadShippingCost();
  }, [selectedCity, totalWeight, cartItems.length]);

  // Get selected shipping cost
  const getSelectedShippingCost = () => {
    if (!selectedCourier || !selectedService) return 0;
    const courier = shippingOptions.find(c => c.code === selectedCourier);
    const service = courier?.costs.find(s => s.service === selectedService);
    return service?.cost[0]?.value || 0;
  };

  const getSelectedEtd = () => {
    if (!selectedCourier || !selectedService) return "";
    const courier = shippingOptions.find(c => c.code === selectedCourier);
    const service = courier?.costs.find(s => s.service === selectedService);
    return service?.cost[0]?.etd || "";
  };

  // Calculate totals
  const subtotal = cartItems.reduce((sum, item) => sum + item.product.price * item.quantity, 0);
  const shippingCost = getSelectedShippingCost();
  const total = subtotal + shippingCost;

  // Get city/province name
  const getCityName = () => cities.find(c => c.city_id === selectedCity)?.city_name || "";
  const getProvinceName = () => provinces.find(p => p.province_id === selectedProvince)?.province || "";

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);

    if (!address || !selectedCity || !selectedCourier || !selectedService) {
      setError("Please fill in all required fields");
      return;
    }

    if (!token) {
      setError("Please login first");
      return;
    }

    setIsSubmitting(true);
    try {
      const courier = shippingOptions.find(c => c.code === selectedCourier);
      
      const order = await createOrder(token, {
        items: cartItems.map(item => ({
          product_id: item.product.id,
          quantity: item.quantity,
        })),
        shipping_address: address,
        shipping_city: getCityName(),
        shipping_city_id: selectedCity,
        shipping_province: getProvinceName(),
        shipping_postal_code: postalCode,
        courier: selectedCourier.toUpperCase(),
        courier_service: selectedService,
        shipping_cost: shippingCost,
        shipping_etd: getSelectedEtd(),
        notes,
      });

      // Clear cart and redirect to order success page
      clearCart();
      router.push(`/orders/${order?.id}?success=true`);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to create order");
    } finally {
      setIsSubmitting(false);
    }
  };

  if (!isAuthenticated || !isCartLoaded || cartItems.length === 0) {
    return (
      <div className="container py-8 flex justify-center">
        <Spinner size="lg" />
      </div>
    );
  }

  return (
    <div className="container py-8 max-w-4xl">
      <Button variant="ghost" onClick={() => router.back()} className="mb-6">
        <ArrowLeft className="h-4 w-4 mr-2" />
        Back
      </Button>

      <h1 className="text-2xl font-bold mb-6">Checkout</h1>

      {error && (
        <Alert variant="destructive" className="mb-6">{error}</Alert>
      )}

      <form onSubmit={handleSubmit} className="space-y-6">
        {/* Order Items */}
        <div className="border rounded-lg p-4">
          <h2 className="font-semibold flex items-center gap-2 mb-4">
            <Package className="h-5 w-5" />
            Order Items
          </h2>
          <div className="space-y-3">
            {cartItems.map((item, index) => (
              <div key={index} className="flex items-center gap-4">
                <div className="flex-1">
                  <p className="font-medium">{item.product.name}</p>
                  <p className="text-sm text-muted-foreground">
                    {formatCurrency(item.product.price)} x {item.quantity}
                  </p>
                </div>
                <p className="font-semibold">
                  {formatCurrency(item.product.price * item.quantity)}
                </p>
              </div>
            ))}
          </div>
        </div>

        {/* Shipping Address */}
        <div className="border rounded-lg p-4">
          <h2 className="font-semibold flex items-center gap-2 mb-4">
            <MapPin className="h-5 w-5" />
            Shipping Address
          </h2>
          <div className="grid gap-4 md:grid-cols-2">
            <div className="space-y-2">
              <Label>Province *</Label>
              <select
                className="w-full h-10 rounded-md border border-input bg-background px-3 text-sm"
                value={selectedProvince}
                onChange={(e) => {
                  setSelectedProvince(e.target.value);
                  setSelectedCity("");
                }}
                disabled={isLoadingProvinces}
              >
                <option value="">Select Province</option>
                {provinces.map((p) => (
                  <option key={p.province_id} value={p.province_id}>
                    {p.province}
                  </option>
                ))}
              </select>
            </div>
            <div className="space-y-2">
              <Label>City *</Label>
              <select
                className="w-full h-10 rounded-md border border-input bg-background px-3 text-sm"
                value={selectedCity}
                onChange={(e) => setSelectedCity(e.target.value)}
                disabled={isLoadingCities || !selectedProvince}
              >
                <option value="">Select City</option>
                {cities.map((c) => (
                  <option key={c.city_id} value={c.city_id}>
                    {c.type} {c.city_name}
                  </option>
                ))}
              </select>
            </div>
            <div className="space-y-2 md:col-span-2">
              <Label>Address *</Label>
              <textarea
                className="w-full min-h-[80px] rounded-md border border-input bg-background px-3 py-2 text-sm"
                placeholder="Enter your complete address"
                value={address}
                onChange={(e) => setAddress(e.target.value)}
              />
            </div>
            <div className="space-y-2">
              <Label>Postal Code</Label>
              <Input
                placeholder="12345"
                value={postalCode}
                onChange={(e) => setPostalCode(e.target.value)}
              />
            </div>
          </div>
        </div>

        {/* Shipping Method */}
        <div className="border rounded-lg p-4">
          <h2 className="font-semibold flex items-center gap-2 mb-4">
            <Truck className="h-5 w-5" />
            Shipping Method
          </h2>
          {isLoadingShipping ? (
            <div className="flex justify-center py-4">
              <Spinner />
            </div>
          ) : shippingOptions.length === 0 ? (
            <p className="text-sm text-muted-foreground">
              {selectedCity ? "No shipping options available" : "Select city to see shipping options"}
            </p>
          ) : (
            <div className="space-y-3">
              {shippingOptions.map((courier) => (
                <div key={courier.code}>
                  <p className="font-medium text-sm mb-2">{courier.name}</p>
                  <div className="grid gap-2 md:grid-cols-2">
                    {courier.costs.map((service) => (
                      <label
                        key={service.service}
                        className={`flex items-center justify-between p-3 border rounded-lg cursor-pointer transition-colors ${
                          selectedCourier === courier.code && selectedService === service.service
                            ? "border-primary bg-primary/5"
                            : "hover:border-primary/50"
                        }`}
                      >
                        <div className="flex items-center gap-3">
                          <input
                            type="radio"
                            name="shipping"
                            className="h-4 w-4"
                            checked={selectedCourier === courier.code && selectedService === service.service}
                            onChange={() => {
                              setSelectedCourier(courier.code);
                              setSelectedService(service.service);
                            }}
                          />
                          <div>
                            <p className="font-medium text-sm">{service.service}</p>
                            <p className="text-xs text-muted-foreground">
                              {service.cost[0]?.etd} days
                            </p>
                          </div>
                        </div>
                        <p className="font-semibold text-sm">
                          {formatCurrency(service.cost[0]?.value || 0)}
                        </p>
                      </label>
                    ))}
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>

        {/* Notes */}
        <div className="border rounded-lg p-4">
          <h2 className="font-semibold mb-4">Notes (Optional)</h2>
          <textarea
            className="w-full min-h-[60px] rounded-md border border-input bg-background px-3 py-2 text-sm"
            placeholder="Add notes for the seller..."
            value={notes}
            onChange={(e) => setNotes(e.target.value)}
          />
        </div>

        {/* Order Summary */}
        <div className="border rounded-lg p-4">
          <h2 className="font-semibold flex items-center gap-2 mb-4">
            <CreditCard className="h-5 w-5" />
            Order Summary
          </h2>
          <div className="space-y-2">
            <div className="flex justify-between text-sm">
              <span>Subtotal</span>
              <span>{formatCurrency(subtotal)}</span>
            </div>
            <div className="flex justify-between text-sm">
              <span>Shipping ({totalWeight}g)</span>
              <span>{shippingCost > 0 ? formatCurrency(shippingCost) : "-"}</span>
            </div>
            <div className="border-t pt-2 mt-2">
              <div className="flex justify-between font-semibold">
                <span>Total</span>
                <span className="text-primary">{formatCurrency(total)}</span>
              </div>
            </div>
          </div>
        </div>

        <Button
          type="submit"
          size="lg"
          className="w-full"
          disabled={isSubmitting || !selectedCourier || !selectedCity || !address}
          loading={isSubmitting}
        >
          Place Order
        </Button>
      </form>
    </div>
  );
}

export default function CheckoutPage() {
  return (
    <Suspense fallback={<div className="container py-8 flex justify-center"><Spinner size="lg" /></div>}>
      <CheckoutContent />
    </Suspense>
  );
}
