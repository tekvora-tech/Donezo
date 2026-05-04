import { apiClient } from "./client";
import { type ApiResponse } from "../types/api";
import {
  type LoginCredentials,
  type RegisterCredentials,
  type AuthResult,
} from "../types/auth";

export const authApi = {
  register: (data: RegisterCredentials) =>
    apiClient.post<ApiResponse<AuthResult>>("/api/auth/v1/register", data),

  login: (data: LoginCredentials) =>
    apiClient.post<ApiResponse<AuthResult>>("/api/auth/v1/login", data),

  logout: () => apiClient.post<ApiResponse<null>>("/api/auth/v1/logout", {}),

  refresh: () =>
    apiClient.post<ApiResponse<{ access_token: string; expires_in: number }>>(
      "/api/auth/v1/refresh",
      {},
    ),
};
