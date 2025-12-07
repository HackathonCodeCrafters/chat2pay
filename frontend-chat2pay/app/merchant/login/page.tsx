"use client";

import Link from "next/link";
import { Store } from "lucide-react";
import { MerchantLoginForm } from "@/features/merchant";

export default function MerchantLoginPage() {
  return (
    <div className="min-h-screen flex flex-col items-center justify-center p-4 bg-muted/30">
      <Link href="/" className="mb-8 flex items-center gap-2">
        <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-primary">
          <Store className="h-6 w-6 text-primary-foreground" />
        </div>
        <span className="text-2xl font-bold">Chat2Pay</span>
      </Link>

      <MerchantLoginForm />

      <p className="mt-8 text-center text-sm text-muted-foreground">
        &copy; 2024 Chat2Pay by Code Crafters
      </p>
    </div>
  );
}
