import apiClient, { extractApiError } from "./client";
import type {
  ApiResponse,
  LoginInput,
  LoginResponse,
  ChangePasswordInput,
  User,
} from "@/types/api.types";

// Login
export async function login(input: LoginInput): Promise<LoginResponse> {
  try {
    const { data } = await apiClient.post<ApiResponse<LoginResponse>>(
      "/auth/login",
      input
    );
    if (!data.data) throw new Error("No data returned");
    return data.data;
  } catch (error) {
    throw extractApiError(error);
  }
}

// Logout
export async function logout(): Promise<void> {
  try {
    await apiClient.post("/auth/logout");
  } catch (error) {
    // Ignore errors on logout
    console.error("Logout error:", error);
  }
}

// Refresh token
export async function refreshToken(): Promise<{
  access_token: string;
  expires_at: string;
}> {
  try {
    const { data } = await apiClient.post<
      ApiResponse<{ access_token: string; expires_at: string }>
    >("/auth/refresh");
    if (!data.data) throw new Error("No data returned");
    return data.data;
  } catch (error) {
    throw extractApiError(error);
  }
}

// Get current user
export async function getMe(): Promise<User> {
  try {
    const { data } = await apiClient.get<ApiResponse<User>>("/auth/me");
    if (!data.data) throw new Error("No data returned");
    return data.data;
  } catch (error) {
    throw extractApiError(error);
  }
}

// Change password
export async function changePassword(
  input: ChangePasswordInput
): Promise<void> {
  try {
    await apiClient.post("/auth/change-password", input);
  } catch (error) {
    throw extractApiError(error);
  }
}
