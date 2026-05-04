import { useState } from "react";
import { motion, AnimatePresence } from "framer-motion";
import { IoSearch, IoFunnel, IoClose, IoChevronDown } from "react-icons/io5";
import { PRIORITY_OPTIONS, STATUS_OPTIONS } from "../../utils/constants";
import { useTags } from "../../hooks/useTags";
import { cn } from "../../utils/cn";

const TodoFilters = ({ filters, onChange }) => {
  const [showAdvanced, setShowAdvanced] = useState(false);
  const { tags } = useTags();

  const hasActiveFilters = Object.values(filters).some(
    (v) => v !== undefined && v !== "",
  );

  const clearFilters = () => {
    onChange({});
  };

  return (
    <motion.div
      initial={{ opacity: 0, y: -10 }}
      animate={{ opacity: 1, y: 0 }}
      className="space-y-3"
    >
      {/* Search + Quick filters */}
      <div className="flex flex-wrap items-center gap-2">
        <div className="relative flex-1 min-w-[200px]">
          <IoSearch
            size={18}
            className="absolute left-3.5 top-1/2 -translate-y-1/2 text-[var(--text-tertiary)]"
          />
          <input
            type="text"
            placeholder="Search todos..."
            value={filters.search || ""}
            onChange={(e) =>
              onChange({ ...filters, search: e.target.value, page: 1 })
            }
            className="w-full pl-10 pr-4 py-2.5 rounded-xl bg-[var(--bg-secondary)] border border-[var(--border-color)] text-[var(--text-primary)] placeholder:text-[var(--text-tertiary)] focus:border-pastel-blue-400 focus:ring-2 focus:ring-pastel-blue-400/20 outline-none transition-all text-sm"
          />
          {filters.search && (
            <button
              onClick={() => onChange({ ...filters, search: undefined })}
              className="absolute right-3 top-1/2 -translate-y-1/2 text-[var(--text-tertiary)] hover:text-[var(--text-primary)]"
            >
              <IoClose size={16} />
            </button>
          )}
        </div>

        <button
          onClick={() => setShowAdvanced(!showAdvanced)}
          className={cn(
            "flex items-center gap-2 px-4 py-2.5 rounded-xl border text-sm font-medium transition-all",
            showAdvanced || hasActiveFilters
              ? "border-pastel-blue-400 text-pastel-blue-400 bg-pastel-blue-50 dark:bg-pastel-blue-900/20"
              : "border-[var(--border-color)] text-[var(--text-secondary)] hover:bg-[var(--bg-tertiary)]",
          )}
        >
          <IoFunnel size={16} />
          Filters
          {(filters.status || filters.priority || filters.tag_id) && (
            <span className="w-5 h-5 rounded-full bg-pastel-blue-400 text-white text-xs flex items-center justify-center">
              {
                [filters.status, filters.priority, filters.tag_id].filter(
                  Boolean,
                ).length
              }
            </span>
          )}
          <IoChevronDown
            size={14}
            className={cn("transition-transform", showAdvanced && "rotate-180")}
          />
        </button>

        {hasActiveFilters && (
          <button
            onClick={clearFilters}
            className="text-sm text-pastel-rose-400 hover:text-pastel-rose-500 font-medium"
          >
            Clear all
          </button>
        )}
      </div>

      {/* Advanced filters */}
      <AnimatePresence>
        {showAdvanced && (
          <motion.div
            initial={{ height: 0, opacity: 0 }}
            animate={{ height: "auto", opacity: 1 }}
            exit={{ height: 0, opacity: 0 }}
            transition={{ duration: 0.2 }}
            className="overflow-hidden"
          >
            <div className="flex flex-wrap gap-3 p-4 rounded-xl bg-[var(--bg-secondary)] border border-[var(--border-color)]">
              {/* Status */}
              <div className="space-y-1.5">
                <label className="text-xs font-semibold text-[var(--text-tertiary)] uppercase">
                  Status
                </label>
                <div className="flex flex-wrap gap-1.5">
                  {STATUS_OPTIONS.map((s) => (
                    <button
                      key={s.value}
                      onClick={() =>
                        onChange({
                          ...filters,
                          status:
                            filters.status === s.value ? undefined : s.value,
                          page: 1,
                        })
                      }
                      className={cn(
                        "px-3 py-1.5 rounded-lg text-xs font-medium transition-all",
                        filters.status === s.value
                          ? "bg-pastel-blue-400 text-white shadow-md shadow-pastel-blue-400/25"
                          : "bg-[var(--bg-tertiary)] text-[var(--text-secondary)] hover:bg-[var(--border-color)]",
                      )}
                    >
                      {s.label}
                    </button>
                  ))}
                </div>
              </div>

              {/* Priority */}
              <div className="space-y-1.5">
                <label className="text-xs font-semibold text-[var(--text-tertiary)] uppercase">
                  Priority
                </label>
                <div className="flex flex-wrap gap-1.5">
                  {PRIORITY_OPTIONS.map((p) => (
                    <button
                      key={p.value}
                      onClick={() =>
                        onChange({
                          ...filters,
                          priority:
                            filters.priority === p.value ? undefined : p.value,
                          page: 1,
                        })
                      }
                      className={cn(
                        "px-3 py-1.5 rounded-lg text-xs font-medium transition-all",
                        filters.priority === p.value
                          ? "bg-pastel-lavender-400 text-white shadow-md shadow-pastel-lavender-400/25"
                          : "bg-[var(--bg-tertiary)] text-[var(--text-secondary)] hover:bg-[var(--border-color)]",
                      )}
                    >
                      {p.icon} {p.label}
                    </button>
                  ))}
                </div>
              </div>

              {/* Tags */}
              {tags.length > 0 && (
                <div className="space-y-1.5">
                  <label className="text-xs font-semibold text-[var(--text-tertiary)] uppercase">
                    Tags
                  </label>
                  <div className="flex flex-wrap gap-1.5">
                    {tags.map((tag) => (
                      <button
                        key={tag.id}
                        onClick={() =>
                          onChange({
                            ...filters,
                            tag_id:
                              filters.tag_id === tag.id ? undefined : tag.id,
                            page: 1,
                          })
                        }
                        className={cn(
                          "px-3 py-1.5 rounded-lg text-xs font-medium transition-all flex items-center gap-1.5",
                          filters.tag_id === tag.id
                            ? "ring-2 ring-offset-1"
                            : "bg-[var(--bg-tertiary)] hover:bg-[var(--border-color)]",
                        )}
                        style={
                          filters.tag_id === tag.id
                            ? {
                                backgroundColor: `${tag.color}20`,
                                color: tag.color,
                                boxShadow: `0 0 0 2px ${tag.color}40`,
                              }
                            : {}
                        }
                      >
                        <span
                          className="w-2 h-2 rounded-full"
                          style={{ backgroundColor: tag.color }}
                        />
                        {tag.name}
                      </button>
                    ))}
                  </div>
                </div>
              )}

              {/* Sort */}
              <div className="space-y-1.5">
                <label className="text-xs font-semibold text-[var(--text-tertiary)] uppercase">
                  Sort
                </label>
                <div className="flex gap-2">
                  <select
                    value={filters.sort_by || "created_at"}
                    onChange={(e) =>
                      onChange({ ...filters, sort_by: e.target.value })
                    }
                    className="px-3 py-1.5 rounded-lg bg-[var(--bg-tertiary)] text-[var(--text-secondary)] text-xs border border-[var(--border-color)] outline-none focus:border-pastel-blue-400"
                  >
                    <option value="created_at">Created</option>
                    <option value="due_date">Due Date</option>
                    <option value="priority">Priority</option>
                    <option value="title">Title</option>
                  </select>
                  <button
                    onClick={() =>
                      onChange({
                        ...filters,
                        sort_order:
                          filters.sort_order === "asc" ? "desc" : "asc",
                      })
                    }
                    className="px-3 py-1.5 rounded-lg bg-[var(--bg-tertiary)] text-[var(--text-secondary)] text-xs border border-[var(--border-color)] hover:bg-[var(--border-color)] transition-colors"
                  >
                    {filters.sort_order === "asc" ? "↑ Asc" : "↓ Desc"}
                  </button>
                </div>
              </div>
            </div>
          </motion.div>
        )}
      </AnimatePresence>
    </motion.div>
  );
};

export default TodoFilters;
