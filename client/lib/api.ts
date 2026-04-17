export const BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8090";

type ApiResponse<T> = {
  success: boolean;
  message: string;
  data?: T;
  error?: string;
};

export async function apiRequest<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<ApiResponse<T>> {
  const token = typeof window !== "undefined" ? localStorage.getItem("access_token") : null;

  const headers = new Headers(options.headers);
  if (token) {
    headers.set("Authorization", `Bearer ${token}`);
  }
  if (!headers.has("Content-Type") && !(options.body instanceof FormData)) {
    headers.set("Content-Type", "application/json");
  }

  const response = await fetch(`${BASE_URL}${endpoint}`, {
    ...options,
    headers,
  });

  const data = await response.json();

  if (!response.ok) {
    // If unauthorized, could clear token and redirect to login
    if (response.status === 401 && typeof window !== "undefined") {
      localStorage.removeItem("access_token");
    }
    return {
      success: false,
      message: data.message || "An error occurred",
      error: data.error || response.statusText,
    };
  }

  return data;
}

export const setAuthToken = (token: string) => {
  localStorage.setItem("access_token", token);
};

export const clearAuthToken = () => {
  localStorage.removeItem("access_token");
};

export const getAuthToken = () => {
  return typeof window !== "undefined" ? localStorage.getItem("access_token") : null;
};
