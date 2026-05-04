import { motion } from "framer-motion";
import { Link } from "react-router-dom";
import {
  IoCalendar,
  IoCheckmarkCircle,
  IoEllipse,
  IoTrash,
  IoPencil,
} from "react-icons/io5";
import Card from "../ui/Card";
import PriorityBadge from "./PriorityBadge";
import TagBadge from "./TagBadge";
import { formatRelativeDate, isOverdue } from "../../utils/formatDate";
import { cn } from "../../utils/cn";

const TodoCard = ({ todo, onDelete, onStatusChange, index = 0 }) => {
  const overdue = isOverdue(todo.due_date);
  const completed = todo.status === "completed";
  const progress = todo.sub_task_count
    ? Math.round(
        ((todo.completed_sub_task_count || 0) / todo.sub_task_count) * 100,
      )
    : 0;

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay: index * 0.05, duration: 0.3 }}
    >
      <Card hover className="group relative overflow-hidden">
        {/* Progress bar top */}
        {todo.sub_task_count > 0 && (
          <div className="absolute top-0 left-0 right-0 h-1 bg-[var(--bg-tertiary)]">
            <motion.div
              initial={{ width: 0 }}
              animate={{ width: `${progress}%` }}
              transition={{ duration: 0.5, delay: 0.2 }}
              className={cn(
                "h-full rounded-full",
                progress === 100 ? "bg-pastel-mint-400" : "bg-pastel-blue-400",
              )}
            />
          </div>
        )}

        <div className="p-5 pt-6">
          <div className="flex items-start justify-between gap-3">
            <div className="flex items-start gap-3 flex-1 min-w-0">
              {/* Checkbox */}
              <motion.button
                whileTap={{ scale: 0.9 }}
                onClick={() =>
                  onStatusChange?.(todo.id, completed ? "pending" : "completed")
                }
                className={cn(
                  "mt-0.5 w-6 h-6 rounded-full border-2 flex items-center justify-center shrink-0 transition-all duration-200",
                  completed
                    ? "bg-pastel-mint-400 border-pastel-mint-400 text-white"
                    : "border-[var(--border-color)] hover:border-pastel-blue-400 text-transparent",
                )}
              >
                <IoCheckmarkCircle size={16} />
              </motion.button>

              <div className="flex-1 min-w-0">
                <Link to={`/todos/${todo.id}`}>
                  <h3
                    className={cn(
                      "font-semibold text-[var(--text-primary)] truncate hover:text-pastel-blue-400 transition-colors",
                      completed && "line-through text-[var(--text-tertiary)]",
                    )}
                  >
                    {todo.title}
                  </h3>
                </Link>
                {todo.description && (
                  <p className="mt-1 text-sm text-[var(--text-secondary)] line-clamp-2">
                    {todo.description}
                  </p>
                )}
              </div>
            </div>

            {/* Actions */}
            <div className="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
              <Link
                to={`/todos/${todo.id}`}
                className="p-1.5 rounded-lg hover:bg-[var(--bg-tertiary)] text-[var(--text-tertiary)] hover:text-pastel-blue-400 transition-colors"
              >
                <IoPencil size={16} />
              </Link>
              <button
                onClick={() => onDelete?.(todo.id)}
                className="p-1.5 rounded-lg hover:bg-pastel-rose-50 dark:hover:bg-pastel-rose-900/20 text-[var(--text-tertiary)] hover:text-pastel-rose-400 transition-colors"
              >
                <IoTrash size={16} />
              </button>
            </div>
          </div>

          {/* Meta row */}
          <div className="mt-4 flex flex-wrap items-center gap-2">
            <PriorityBadge priority={todo.priority} size="sm" />

            {todo.due_date && (
              <span
                className={cn(
                  "inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs font-medium",
                  overdue
                    ? "bg-pastel-rose-100 text-pastel-rose-600 dark:bg-pastel-rose-900/20"
                    : "bg-[var(--bg-tertiary)] text-[var(--text-tertiary)]",
                )}
              >
                <IoCalendar size={12} />
                {formatRelativeDate(todo.due_date)}
              </span>
            )}

            {todo.sub_task_count > 0 && (
              <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs font-medium bg-[var(--bg-tertiary)] text-[var(--text-tertiary)]">
                <IoEllipse
                  size={10}
                  className={
                    progress === 100
                      ? "text-pastel-mint-400"
                      : "text-pastel-blue-400"
                  }
                />
                {todo.completed_sub_task_count}/{todo.sub_task_count}
              </span>
            )}

            {/* Tags */}
            {todo.tags?.map((tag) => (
              <TagBadge key={tag.id} tag={tag} size="sm" />
            ))}
          </div>
        </div>
      </Card>
    </motion.div>
  );
};

export default TodoCard;
