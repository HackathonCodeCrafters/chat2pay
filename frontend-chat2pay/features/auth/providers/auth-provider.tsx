"use client";

import * as React from "react";
import { useRouter } from "next/navigation";
import {
  User,
  AuthState,
  LoginCredentials,
  CustomerRegisterData,
  CustomerAuthResponse,
} from "../types";
import { loginCustomer, registerCustomer } from "../api";

const TOKEN_KEY = "chat2pay_token";
const USER_KEY = "chat2pay_user";

interface AuthContextValue extends AuthState {
  login: (credentials: LoginCredentials) => Promise<void>;
  register: (data: CustomerRegisterData) => Promise<void>;
  logout: () => void;
  getToken: () => string | null;
}

export const AuthContext = React.createContext<AuthContextValue | null>(null);

function getStoredAuth(): { token: string | null; user: User | null } {
  if (typeof window === "undefined") {
    return { token: null, user: null };
  }

  const token = localStorage.getItem(TOKEN_KEY);
  const userStr = localStorage.getItem(USER_KEY);
  const user = userStr ? JSON.parse(userStr) : null;

  return { token, user };
}

function setStoredAuth(token: string, user: User) {
  localStorage.setItem(TOKEN_KEY, token);
  localStorage.setItem(USER_KEY, JSON.stringify(user));
}

function clearStoredAuth() {
  localStorage.removeItem(TOKEN_KEY);
  localStorage.removeItem(USER_KEY);
}

function mapResponseToUser(response: CustomerAuthResponse): User {
  return {
    id: response.id,
    name: response.name,
    email: response.email || "",
    phone: response.phone || undefined,
    role: "customer",
  };
}

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const router = useRouter();
  const [state, setState] = React.useState<AuthState>({
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
    async (credentials: LoginCredentials) => {
      setState((prev) => ({ ...prev, isLoading: true }));

      try {
        const response = await loginCustomer(credentials);
        const user = mapResponseToUser(response);

        setStoredAuth(response.access_token, user);
        setState({
          user,
          token: response.access_token,
          isLoading: false,
          isAuthenticated: true,
        });

        router.push("/chat");
      } catch (error) {
        setState((prev) => ({ ...prev, isLoading: false }));
        throw error;
      }
    },
    [router]
  );

  const register = React.useCallback(
    async (data: CustomerRegisterData) => {
      setState((prev) => ({ ...prev, isLoading: true }));

      try {
        const response = await registerCustomer(data);
        const user = mapResponseToUser(response);

        setStoredAuth(response.access_token, user);
        setState({
          user,
          token: response.access_token,
          isLoading: false,
          isAuthenticated: true,
        });

        router.push("/chat");
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
    router.push("/login");
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

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}
