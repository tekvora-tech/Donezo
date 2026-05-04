export interface ApiResponse<T> {
  message: string;
  code: number;
  result: T;
  error?: string;
}

export interface ApiError {
  message: string;
  code: number;
  error: string;
}

export interface PaginationMetadata {
  current_page: number;
  page_size: number;
  total_pages: number;
  total_records: number;
  has_next: boolean;
  has_prev: boolean;
}

export interface PaginatedResponse<T> {
  data: T[];
  metadata: PaginationMetadata;
}

export type Priority = "low" | "medium" | "high" | "urgent";
export type Status = "pending" | "in_progress" | "completed" | "cancelled";
