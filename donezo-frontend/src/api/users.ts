import { apiClient } from "./client";
import { type ApiResponse } from "../types/api";
import { type UpdateProfileInput } from "../types/user";
import { type User } from "../types/auth";

export const usersApi = {
  getProfile: () => apiClient.get<ApiResponse<User>>("/api/users/v1/profile"),

  updateProfile: (data: UpdateProfileInput) =>
    apiClient.put<ApiResponse<User>>("/api/users/v1/profile", data),
};
