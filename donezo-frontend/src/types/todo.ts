import { type Priority, type Status } from "./api";

export interface SubTask {
  id: string;
  todo_id: string;
  title: string;
  is_completed: boolean;
  created_at: string;
}

export interface TodoTag {
  id: string;
  name: string;
  color: string;
}

export interface Todo {
  id: string;
  user_id: string;
  title: string;
  description: string | null;
  status: Status;
  priority: Priority;
  due_date: string | null;
  tags: TodoTag[];
  sub_tasks: SubTask[];
  sub_task_count?: number;
  completed_sub_task_count?: number;
  created_at: string;
  updated_at: string;
}

export interface CreateTodoInput {
  title: string;
  description?: string;
  priority?: Priority;
  due_date?: string;
  tag_ids?: string[];
}

export interface UpdateTodoInput {
  title?: string;
  description?: string;
  status?: Status;
  priority?: Priority;
  due_date?: string;
  tag_ids?: string[];
}

export interface TodoFilters {
  page?: number;
  page_size?: number;
  status?: Status;
  priority?: Priority;
  tag_id?: string;
  search?: string;
  due_date_from?: string;
  due_date_to?: string;
  sort_by?: "created_at" | "updated_at" | "due_date" | "priority" | "title";
  sort_order?: "asc" | "desc";
}
