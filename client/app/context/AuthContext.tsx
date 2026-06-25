"use client";

import React, { createContext, useContext, useEffect, useState } from "react";
import { apiRequest, clearAuthToken, getAuthToken } from "@/lib/api";

type User = {
  id: string;
  email: string;
  role: string;
};

type AuthContextType = {
  user: User | null;
  loading: boolean;
  login: (token: string, refreshToken?: string) => Promise<void>;
  logout: () => Promise<void>;
};

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  const fetchMe = async () => {
    const res = await apiRequest<User>("/auth/me");
    if (res.success && res.data) {
      setUser(res.data);
    } else {
      setUser(null);
    }
    setLoading(false);
  };

  useEffect(() => {
    const token = getAuthToken();
    if (token) {
      fetchMe();
    } else {
      setLoading(false);
    }
  }, []);

  const login = async (token: string, refreshToken?: string) => {
    localStorage.setItem("access_token", token);
    if (refreshToken) {
      localStorage.setItem("refresh_token", refreshToken);
    }
    await fetchMe();
  };

  const logout = async () => {
    const refreshToken = localStorage.getItem("refresh_token");
    if (refreshToken) {
      try {
        await apiRequest("/auth/logout", {
          method: "POST",
          body: JSON.stringify({ refresh_token: refreshToken }),
        });
      } catch (e) {
        console.error("Logout request failed", e);
      }
    }
    clearAuthToken();
    localStorage.removeItem("refresh_token");
    setUser(null);
  };

  return (
    <AuthContext.Provider value={{ user, loading, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
};
