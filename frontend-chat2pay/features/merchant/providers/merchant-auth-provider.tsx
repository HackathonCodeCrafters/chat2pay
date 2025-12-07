"use client";

import * as React from "react";
import { useRouter } from "next/navigation";
import {
  MerchantUser,
  MerchantLoginCredentials,
  MerchantRegisterData,
  MerchantAuthResponse,
} from "../types";
import { loginMerchant, registerMerchant } from "../api";

const TOKEN_KEY = "chat2pay_merchant_token";
const USER_KEY = "chat2pay_merchant_user";

interface MerchantAuthState {
  user: MerchantUser | null;
  token: string | null;
  isLoading: boolean;
  isAuthenticated: boolean;
}

interface MerchantAuthContextValue extends MerchantAuthState {
  login: (credentials: MerchantLoginCredentials) => Promise<void>;
  register: (data: MerchantRegisterData) => Promise<void>;
  logout: () => void;
  getToken: () => string | null;
}

export const MerchantAuthContext = React.createContext<MerchantAuthContextValue | null>(null);

function getStoredAuth(): { token: string | null; user: MerchantUser | null } {
  if (typeof window === "undefined") {
    return { token: null, user: null };
  }

  const token = localStorage.getItem(TOKEN_KEY);
  const userStr = localStorage.getItem(USER_KEY);
  const user = userStr ? JSON.parse(userStr) : null;

  return { token, user };
}

function setStoredAuth(token: string, user: MerchantUser) {
  localStorage.setItem(TOKEN_KEY, token);
  localStorage.setItem(USER_KEY, JSON.stringify(user));
}

function clearStoredAuth() {
  localStorage.removeItem(TOKEN_KEY);
  localStorage.removeItem(USER_KEY);
}

function mapResponseToUser(response: MerchantAuthResponse): MerchantUser {
  return {
    id: response.id,
    merchant_id: response.merchant_id,
    merchant: response.merchant,
    name: response.name,
    email: response.email,
    role: response.role,
    status: response.status,
  };
}

export function MerchantAuthProvider({ children }: { children: React.ReactNode }) {
  const router = useRouter();
  const [state, setState] = React.useState<MerchantAuthState>({
    user: null,
    token: null,
    isLoading: true,
    isAuthenticated: false,
  });

  React.useEffect(() => {
    const { token, user } = getStoredAuth();
    setState({
      user,
      token,
      isLoading: false,
      isAuthenticated: !!token && !!user,
    });
  }, []);

  const login = React.useCallback(
    async (credentials: MerchantLoginCredentials) => {
      setState((prev) => ({ ...prev, isLoading: true }));

      try {
        const response = await loginMerchant(credentials);
        const user = mapResponseToUser(response);

        setStoredAuth(response.access_token, user);
        setState({
          user,
          token: response.access_token,
          isLoading: false,
          isAuthenticated: true,
        });

        router.push("/merchant/dashboard");
      } catch (error) {
        setState((prev) => ({ ...prev, isLoading: false }));
        throw error;
      }
    },
    [router]
  );

  const register = React.useCallback(
    async (data: MerchantRegisterData) => {
      setState((prev) => ({ ...prev, isLoading: true }));

      try {
        const response = await registerMerchant(data);
        const user = mapResponseToUser(response);

        setStoredAuth(response.access_token, user);
        setState({
          user,
          token: response.access_token,
          isLoading: false,
          isAuthenticated: true,
        });

        router.push("/merchant/dashboard");
      } catch (error) {
        setState((prev) => ({ ...prev, isLoading: false }));
        throw error;
      }
    },
    [router]
  );

  const logout = React.useCallback(() => {
    clearStoredAuth();
    setState({
      user: null,
      token: null,
      isLoading: false,
      isAuthenticated: false,
    });
    router.push("/merchant/login");
  }, [router]);

  const getToken = React.useCallback(() => {
    return state.token;
  }, [state.token]);

  const value = React.useMemo(
    () => ({
      ...state,
      login,
      register,
      logout,
      getToken,
    }),
    [state, login, register, logout, getToken]
  );

  return (
    <MerchantAuthContext.Provider value={value}>
      {children}
    </MerchantAuthContext.Provider>
  );
}
