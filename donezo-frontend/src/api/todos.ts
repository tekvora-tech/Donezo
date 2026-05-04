import { apiClient } from "./client";
import { type ApiResponse, type PaginatedResponse } from "../types/api";
import {
  type Todo,
  type CreateTodoInput,
  type UpdateTodoInput,
  type TodoFilters,
  type SubTask,
} from "../types/todo";

export const todosApi = {
  create: (data: CreateTodoInput) =>
    apiClient.post<ApiResponse<Todo>>("/api/todos/v1/", data),

  list: (filters: TodoFilters = {}) => {
    const params: Record<string, string> = {};
    Object.entries(filters).forEach(([key, value]) => {
      if (value !== undefined && value !== null) {
        params[key] = String(value);
      }
    });
    return apiClient.get<ApiResponse<PaginatedResponse<Todo>>>(
      "/api/todos/v1/",
      params,
    );
  },

  getById: (id: string) =>
    apiClient.get<ApiResponse<Todo>>(`/api/todos/v1/${id}`),

  update: (id: string, data: UpdateTodoInput) =>
    apiClient.put<ApiResponse<Todo>>(`/api/todos/v1/${id}`, data),

  delete: (id: string) =>
    apiClient.delete<ApiResponse<null>>(`/api/todos/v1/${id}`),

  // Sub-tasks
  createSubTask: (todoId: string, data: { title: string }) =>
    apiClient.post<ApiResponse<SubTask>>(
      `/api/todos/v1/${todoId}/sub-tasks`,
      data,
    ),

  listSubTasks: (todoId: string) =>
    apiClient.get<ApiResponse<{ data: SubTask[] }>>(
      `/api/todos/v1/${todoId}/sub-tasks`,
    ),

  updateSubTask: (
    todoId: string,
    subTaskId: string,
    data: { title?: string; is_completed?: boolean },
  ) =>
    apiClient.put<ApiResponse<SubTask>>(
      `/api/todos/v1/${todoId}/sub-tasks/${subTaskId}`,
      data,
    ),

  deleteSubTask: (todoId: string, subTaskId: string) =>
    apiClient.delete<ApiResponse<null>>(
      `/api/todos/v1/${todoId}/sub-tasks/${subTaskId}`,
    ),

  // Tags
  assignTag: (todoId: string, tagId: string) =>
    apiClient.post<ApiResponse<Todo>>(`/api/todos/v1/${todoId}/tags`, {
      tag_id: tagId,
    }),

  removeTag: (todoId: string, tagId: string) =>
    apiClient.delete<ApiResponse<Todo>>(
      `/api/todos/v1/${todoId}/tags/${tagId}`,
    ),
};
