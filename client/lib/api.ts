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

  let response = await fetch(`${BASE_URL}${endpoint}`, {
    ...options,
    headers,
  });

  // Check if response has body (e.g. 204 has no body)
  let data: any = {};
  const contentType = response.headers.get("content-type");
  if (contentType && contentType.includes("application/json")) {
    data = await response.json();
  }

  if (!response.ok) {
    // If unauthorized, try to refresh token (avoid loop if already refreshing/logging in)
    if (
      response.status === 401 &&
      typeof window !== "undefined" &&
      endpoint !== "/auth/login" &&
      endpoint !== "/auth/refresh"
    ) {
      const refreshToken = localStorage.getItem("refresh_token");
      if (refreshToken) {
        try {
          const refreshRes = await fetch(`${BASE_URL}/auth/refresh`, {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
            },
            body: JSON.stringify({ refresh_token: refreshToken }),
          });

          if (refreshRes.ok) {
            const refreshData = await refreshRes.json();
            if (refreshData.success && refreshData.data) {
              const { access_token, refresh_token: newRefreshToken } = refreshData.data;
              localStorage.setItem("access_token", access_token);
              if (newRefreshToken) {
                localStorage.setItem("refresh_token", newRefreshToken);
              }

              // Retry the original request with the new token
              headers.set("Authorization", `Bearer ${access_token}`);
              response = await fetch(`${BASE_URL}${endpoint}`, {
                ...options,
                headers,
              });

              const retryContentType = response.headers.get("content-type");
              if (retryContentType && retryContentType.includes("application/json")) {
                data = await response.json();
              } else {
                data = {};
              }

              if (response.ok) {
                return data;
              }
            }
          }
        } catch (e) {
          console.error("Token refresh failed", e);
        }
      }

      // If refresh failed or was not possible, clean up
      localStorage.removeItem("access_token");
      localStorage.removeItem("refresh_token");
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
