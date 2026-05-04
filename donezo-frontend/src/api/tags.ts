import { apiClient } from "./client";
import { type ApiResponse } from "../types/api";
import {
  type Tag,
  type CreateTagInput,
  type UpdateTagInput,
} from "../types/tag";

export const tagsApi = {
  create: (data: CreateTagInput) =>
    apiClient.post<ApiResponse<Tag>>("/api/tags/v1/", data),

  list: () => apiClient.get<ApiResponse<{ data: Tag[] }>>("/api/tags/v1/"),

  update: (id: string, data: UpdateTagInput) =>
    apiClient.put<ApiResponse<Tag>>(`/api/tags/v1/${id}`, data),

  delete: (id: string) =>
    apiClient.delete<ApiResponse<null>>(`/api/tags/v1/${id}`),
};
