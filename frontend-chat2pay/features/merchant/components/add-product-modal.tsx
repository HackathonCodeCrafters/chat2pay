"use client";

import * as React from "react";
import { X, Upload, Image as ImageIcon } from "lucide-react";
import { Button, Input, Label } from "@/shared/components/atoms";
import { Alert } from "@/shared/components/molecules";
import { CreateProductData } from "../types";

interface AddProductModalProps {
  isOpen: boolean;
  isLoading: boolean;
  onClose: () => void;
  onSubmit: (data: CreateProductData) => Promise<void>;
}

export function AddProductModal({ isOpen, isLoading, onClose, onSubmit }: AddProductModalProps) {
  const [error, setError] = React.useState<string | null>(null);
  const [imagePreview, setImagePreview] = React.useState<string | null>(null);
  const fileInputRef = React.useRef<HTMLInputElement>(null);
  const [formData, setFormData] = React.useState({
    name: "",
    description: "",
    sku: "",
    price: "",
    stock: "",
    image: "",
    weight: "",
    length: "",
    width: "",
    height: "",
  });

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
    setError(null);
  };

  const handleImageChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    if (file.size > 2 * 1024 * 1024) {
      setError("Image size must be less than 2MB");
      return;
    }

    const reader = new FileReader();
    reader.onloadend = () => {
      const base64 = reader.result as string;
      setImagePreview(base64);
      setFormData((prev) => ({ ...prev, image: base64 }));
    };
    reader.readAsDataURL(file);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);

    if (!formData.name || !formData.price || !formData.stock) {
      setError("Please fill in all required fields");
      return;
    }

    try {
      await onSubmit({
        name: formData.name,
        description: formData.description || undefined,
        sku: formData.sku || undefined,
        price: parseFloat(formData.price),
        stock: parseInt(formData.stock),
        image: formData.image || undefined,
        weight: formData.weight ? parseInt(formData.weight) : undefined,
        length: formData.length ? parseInt(formData.length) : undefined,
        width: formData.width ? parseInt(formData.width) : undefined,
        height: formData.height ? parseInt(formData.height) : undefined,
      });
      
      // Reset form
      setFormData({ name: "", description: "", sku: "", price: "", stock: "", image: "", weight: "", length: "", width: "", height: "" });
      setImagePreview(null);
      onClose();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to create product");
    }
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center">
      <div className="absolute inset-0 bg-black/50" onClick={onClose} />
      <div className="relative bg-background rounded-lg shadow-lg w-full max-w-md mx-4 p-6">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-lg font-semibold">Add New Product</h2>
          <button onClick={onClose} className="text-muted-foreground hover:text-foreground">
            <X className="h-5 w-5" />
          </button>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4">
          {error && <Alert variant="destructive">{error}</Alert>}

          <div className="space-y-2">
            <Label htmlFor="name">Product Name *</Label>
            <Input
              id="name"
              name="name"
              value={formData.name}
              onChange={handleChange}
              placeholder="Laptop Gaming ASUS"
              disabled={isLoading}
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="description">Description</Label>
            <textarea
              id="description"
              name="description"
              value={formData.description}
              onChange={handleChange}
              placeholder="Product description..."
              className="w-full min-h-[80px] rounded-md border border-input bg-background px-3 py-2 text-sm"
              disabled={isLoading}
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="sku">SKU</Label>
            <Input
              id="sku"
              name="sku"
              value={formData.sku}
              onChange={handleChange}
              placeholder="SKU-001"
              disabled={isLoading}
            />
          </div>

          <div className="space-y-2">
            <Label>Product Image</Label>
            <input
              ref={fileInputRef}
              type="file"
              accept="image/*"
              onChange={handleImageChange}
              className="hidden"
            />
            <div
              onClick={() => fileInputRef.current?.click()}
              className="border-2 border-dashed border-input rounded-lg p-4 text-center cursor-pointer hover:border-primary transition-colors"
            >
              {imagePreview ? (
                <img
                  src={imagePreview}
                  alt="Preview"
                  className="mx-auto h-24 w-24 object-cover rounded"
                />
              ) : (
                <div className="flex flex-col items-center gap-2 text-muted-foreground">
                  <Upload className="h-8 w-8" />
                  <span className="text-sm">Click to upload image</span>
                  <span className="text-xs">Max 2MB</span>
                </div>
              )}
            </div>
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="price">Price (Rp) *</Label>
              <Input
                id="price"
                name="price"
                type="number"
                value={formData.price}
                onChange={handleChange}
                placeholder="15000000"
                disabled={isLoading}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="stock">Stock *</Label>
              <Input
                id="stock"
                name="stock"
                type="number"
                value={formData.stock}
                onChange={handleChange}
                placeholder="100"
                disabled={isLoading}
              />
            </div>
          </div>

          {/* Shipping Dimensions */}
          <div className="space-y-2">
            <Label className="text-sm font-medium">Shipping Info (for delivery cost)</Label>
            <div className="grid grid-cols-2 gap-3">
              <div className="space-y-1">
                <Label htmlFor="weight" className="text-xs text-muted-foreground">Weight (gram)</Label>
                <Input
                  id="weight"
                  name="weight"
                  type="number"
                  value={formData.weight}
                  onChange={handleChange}
                  placeholder="1000"
                  disabled={isLoading}
                />
              </div>
              <div className="space-y-1">
                <Label htmlFor="length" className="text-xs text-muted-foreground">Length (cm)</Label>
                <Input
                  id="length"
                  name="length"
                  type="number"
                  value={formData.length}
                  onChange={handleChange}
                  placeholder="30"
                  disabled={isLoading}
                />
              </div>
              <div className="space-y-1">
                <Label htmlFor="width" className="text-xs text-muted-foreground">Width (cm)</Label>
                <Input
                  id="width"
                  name="width"
                  type="number"
                  value={formData.width}
                  onChange={handleChange}
                  placeholder="20"
                  disabled={isLoading}
                />
              </div>
              <div className="space-y-1">
                <Label htmlFor="height" className="text-xs text-muted-foreground">Height (cm)</Label>
                <Input
                  id="height"
                  name="height"
                  type="number"
                  value={formData.height}
                  onChange={handleChange}
                  placeholder="10"
                  disabled={isLoading}
                />
              </div>
            </div>
          </div>

          <div className="flex gap-3 pt-4">
            <Button type="button" variant="outline" onClick={onClose} className="flex-1">
              Cancel
            </Button>
            <Button type="submit" loading={isLoading} className="flex-1">
              Add Product
            </Button>
          </div>
        </form>
      </div>
    </div>
  );
}
