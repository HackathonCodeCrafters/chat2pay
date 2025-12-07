"use client";

import { useContext } from "react";
import { MerchantAuthContext } from "../providers/merchant-auth-provider";

export function useMerchantAuth() {
  const context = useContext(MerchantAuthContext);

  if (!context) {
    throw new Error("useMerchantAuth must be used within a MerchantAuthProvider");
  }

  return context;
}
