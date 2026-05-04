import { useState } from "react";
import { useSearchParams } from "react-router-dom";
import { motion, AnimatePresence } from "framer-motion";
import {
  IoAdd,
  IoChevronUp,
  IoChevronDown,
  IoFilter,
  IoGrid,
  IoList,
} from "react-icons/io5";

import TodoCard from "../components/todos/TodoCard";
import TodoFilters from "../components/todos/TodoFilters";
import TodoForm from "../components/todos/TodoForm";
import Skeleton from "../components/ui/Skeleton";
import Button from "../components/ui/Button";

import { useTodos } from "../hooks/useTodos";
import { cn } from "../utils/cn";

import type { Status, Priority } from "../types/api";
import type { TodoFilters as TFilters } from "../types/todo";

// 🔥 HARUS ADA (runtime value, bukan type)
const STATUS_VALUES: Status[] = ["pending", "in_progress", "completed"];
const PRIORITY_VALUES: Priority[] = ["low", "medium", "high"];

// helper generic
const parseEnum = <T extends string>(
  value: string | null,
  allowed: readonly T[],
): T | undefined => {
  return value && allowed.includes(value as T) ? (value as T) : undefined;
};

const DashboardPage = () => {
  const [searchParams, setSearchParams] = useSearchParams();
  const [formOpen, setFormOpen] = useState(false);
  const [viewMode, setViewMode] = useState<"list" | "grid">("list");
  const SORT_BY_VALUES = [
    "priority",
    "created_at",
    "updated_at",
    "due_date",
    "title",
  ] as const;

  type SortBy = (typeof SORT_BY_VALUES)[number];
  const filters: TFilters = {
    page: parseInt(searchParams.get("page") || "1"),
    page_size: 10,

    status: parseEnum(searchParams.get("status"), STATUS_VALUES),
    priority: parseEnum(searchParams.get("priority"), PRIORITY_VALUES),

    tag_id: searchParams.get("tag_id") || undefined,
    search: searchParams.get("search") || undefined,

    sort_by:
      parseEnum(searchParams.get("sort_by"), SORT_BY_VALUES) ?? "created_at",
    sort_order: searchParams.get("sort_order") === "asc" ? "asc" : "desc",
  };

  const {
    todos,
    metadata,
    isLoading,
    createTodo,
    updateTodo,
    deleteTodo,
    isCreating,
  } = useTodos(filters);

  const updateFilters = (newFilters: Partial<TFilters>) => {
    const params = new URLSearchParams();

    Object.entries(newFilters).forEach(([key, value]) => {
      if (value !== undefined && value !== "" && value !== null) {
        params.set(key, String(value));
      }
    });

    setSearchParams(params);
  };

  const handlePageChange = (page: number) => {
    updateFilters({ ...filters, page });
    window.scrollTo({ top: 0, behavior: "smooth" });
  };

  const handleStatusChange = async (id: string, status: Status) => {
    await updateTodo({ id, data: { status } });
  };

  return (
    <div className="max-w-5xl mx-auto space-y-6">
      {/* Header */}
      <motion.div
        initial={{ opacity: 0, y: -10 }}
        animate={{ opacity: 1, y: 0 }}
        className="flex flex-wrap items-center justify-between gap-4"
      >
        <div>
          <h1 className="text-2xl font-bold text-[var(--text-primary)]">
            My Todos
          </h1>
          <p className="text-sm text-[var(--text-secondary)] mt-0.5">
            {metadata?.total_records || 0} tasks total
          </p>
        </div>

        <div className="flex items-center gap-2">
          <div className="flex items-center bg-[var(--bg-secondary)] rounded-xl border border-[var(--border-color)] p-1">
            <button
              onClick={() => setViewMode("list")}
              className={cn(
                "p-2 rounded-lg transition-all",
                viewMode === "list"
                  ? "bg-[var(--bg-primary)] shadow-sm text-pastel-blue-400"
                  : "text-[var(--text-tertiary)]",
              )}
            >
              <IoList size={18} />
            </button>

            <button
              onClick={() => setViewMode("grid")}
              className={cn(
                "p-2 rounded-lg transition-all",
                viewMode === "grid"
                  ? "bg-[var(--bg-primary)] shadow-sm text-pastel-blue-400"
                  : "text-[var(--text-tertiary)]",
              )}
            >
              <IoGrid size={18} />
            </button>
          </div>

          <Button onClick={() => setFormOpen(true)}>
            <IoAdd size={18} /> New Todo
          </Button>
        </div>
      </motion.div>

      {/* Filters */}
      <TodoFilters filters={filters} onChange={updateFilters} />

      {/* Content */}
      {isLoading ? (
        <div
          className={cn("grid gap-4", viewMode === "grid" && "sm:grid-cols-2")}
        >
          {[1, 2, 3, 4].map((i) => (
            <Skeleton key={i} className="h-40 rounded-2xl" />
          ))}
        </div>
      ) : todos.length === 0 ? (
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          className="text-center py-20"
        >
          <div className="w-20 h-20 rounded-full bg-pastel-blue-50 dark:bg-pastel-blue-900/20 flex items-center justify-center mx-auto mb-4">
            <IoFilter size={32} className="text-pastel-blue-400" />
          </div>

          <h3 className="text-lg font-semibold text-[var(--text-primary)] mb-2">
            No todos found
          </h3>

          <p className="text-[var(--text-secondary)] mb-6">
            {filters.search || filters.status || filters.priority
              ? "Try adjusting your filters"
              : "Get started by creating your first todo"}
          </p>

          {!filters.search && !filters.status && !filters.priority && (
            <Button onClick={() => setFormOpen(true)}>
              <IoAdd size={18} /> Create Todo
            </Button>
          )}
        </motion.div>
      ) : (
        <>
          <div
            className={cn(
              "grid gap-4",
              viewMode === "grid" && "sm:grid-cols-2",
            )}
          >
            <AnimatePresence mode="popLayout">
              {todos.map((todo, i) => (
                <TodoCard
                  key={todo.id}
                  todo={todo}
                  index={i}
                  onDelete={deleteTodo}
                  onStatusChange={handleStatusChange}
                />
              ))}
            </AnimatePresence>
          </div>

          {/* Pagination */}
          {metadata && metadata.total_pages > 1 && (
            <div className="flex items-center justify-center gap-2 pt-4">
              <Button
                variant="secondary"
                size="sm"
                onClick={() => handlePageChange(filters.page - 1)}
                disabled={!metadata.has_prev}
              >
                <IoChevronUp size={16} />
              </Button>

              <div className="flex items-center gap-1">
                {Array.from(
                  { length: metadata.total_pages },
                  (_, i) => i + 1,
                ).map((page) => (
                  <button
                    key={page}
                    onClick={() => handlePageChange(page)}
                    className={cn(
                      "w-9 h-9 rounded-lg text-sm font-medium transition-all",
                      page === filters.page
                        ? "bg-pastel-blue-400 text-white shadow-md shadow-pastel-blue-400/25"
                        : "text-[var(--text-secondary)] hover:bg-[var(--bg-tertiary)]",
                    )}
                  >
                    {page}
                  </button>
                ))}
              </div>

              <Button
                variant="secondary"
                size="sm"
                onClick={() => handlePageChange(filters.page + 1)}
                disabled={!metadata.has_next}
              >
                <IoChevronDown size={16} />
              </Button>
            </div>
          )}
        </>
      )}

      {/* Create Modal */}
      <TodoForm
        isOpen={formOpen}
        onClose={() => setFormOpen(false)}
        onSubmit={async (data) => {
          await createTodo(data);
          setFormOpen(false);
        }}
        isLoading={isCreating}
        initialData={undefined}
      />
    </div>
  );
};

export default DashboardPage;
