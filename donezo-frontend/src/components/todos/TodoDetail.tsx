import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { motion } from "framer-motion";
import {
  IoArrowBack,
  IoCalendar,
  IoFlag,
  IoPricetag,
  IoTrash,
  IoPencil,
  IoCheckmarkCircle,
} from "react-icons/io5";
import Card from "../ui/Card";
import Button from "../ui/Button";
import Badge from "../ui/Badge";
import PriorityBadge from "./PriorityBadge";
import TagBadge from "./TagBadge";
import SubTaskList from "./SubTaskList";
import TodoForm from "./TodoForm";
import { useTodoDetail } from "../../hooks/useTodos";
import { formatDate } from "../../utils/formatDate";
import { STATUS_COLORS } from "../../utils/constants";
import { cn } from "../../utils/cn";

const TodoDetail = ({ todoId }) => {
  const navigate = useNavigate();
  const [editModalOpen, setEditModalOpen] = useState(false);
  const { todo, isLoading, updateTodo, deleteTodo, isUpdating, isDeleting } =
    useTodoDetail(todoId);

  if (isLoading) {
    return (
      <div className="space-y-4">
        {[1, 2, 3].map((i) => (
          <div
            key={i}
            className="h-32 bg-[var(--bg-tertiary)] rounded-2xl animate-pulse"
          />
        ))}
      </div>
    );
  }

  if (!todo) {
    return (
      <div className="text-center py-20">
        <p className="text-[var(--text-tertiary)] text-lg">Todo not found</p>
        <Button onClick={() => navigate("/dashboard")} className="mt-4">
          Back to Dashboard
        </Button>
      </div>
    );
  }

  const statusColor = STATUS_COLORS[todo.status] || STATUS_COLORS.pending;

  return (
    <div className="max-w-3xl mx-auto space-y-6">
      {/* Header */}
      <div className="flex items-center gap-3">
        <button
          onClick={() => navigate("/dashboard")}
          className="p-2.5 rounded-xl hover:bg-[var(--bg-tertiary)] text-[var(--text-secondary)] transition-colors"
        >
          <IoArrowBack size={20} />
        </button>
        <div className="flex-1" />
        <Button
          variant="secondary"
          size="sm"
          onClick={() => setEditModalOpen(true)}
        >
          <IoPencil size={16} /> Edit
        </Button>
        <Button
          variant="danger"
          size="sm"
          onClick={async () => {
            if (confirm("Are you sure?")) {
              await deleteTodo();
              navigate("/dashboard");
            }
          }}
          isLoading={isDeleting}
        >
          <IoTrash size={16} /> Delete
        </Button>
      </div>

      {/* Main Card */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
      >
        <Card className="p-6 space-y-6">
          {/* Title & Status */}
          <div className="flex items-start justify-between gap-4">
            <div className="flex-1">
              <div className="flex items-center gap-2 mb-2">
                <Badge
                  className={""}
                  variant={
                    todo.status === "completed"
                      ? "success"
                      : todo.status === "in_progress"
                        ? "warning"
                        : "default"
                  }
                >
                  <span className="capitalize flex items-center gap-1">
                    {todo.status === "completed" && (
                      <IoCheckmarkCircle size={14} />
                    )}
                    {todo.status.replace("_", " ")}
                  </span>
                </Badge>
                <PriorityBadge priority={todo.priority} />
              </div>
              <h1
                className={cn(
                  "text-2xl font-bold text-[var(--text-primary)]",
                  todo.status === "completed" &&
                    "line-through text-[var(--text-tertiary)]",
                )}
              >
                {todo.title}
              </h1>
            </div>
          </div>

          {/* Description */}
          {todo.description && (
            <div className="bg-[var(--bg-tertiary)] rounded-xl p-4">
              <p className="text-[var(--text-secondary)] whitespace-pre-wrap">
                {todo.description}
              </p>
            </div>
          )}

          {/* Meta */}
          <div className="flex flex-wrap gap-4 text-sm">
            {todo.due_date && (
              <div className="flex items-center gap-2 text-[var(--text-secondary)]">
                <IoCalendar size={16} className="text-pastel-blue-400" />
                <span>Due {formatDate(todo.due_date)}</span>
              </div>
            )}
            <div className="flex items-center gap-2 text-[var(--text-secondary)]">
              <IoFlag size={16} className="text-pastel-lavender-400" />
              <span className="capitalize">{todo.priority} priority</span>
            </div>
          </div>

          {/* Tags */}
          {todo.tags?.length > 0 && (
            <div className="flex items-center gap-2 flex-wrap">
              <IoPricetag size={16} className="text-[var(--text-tertiary)]" />
              {todo.tags.map((tag) => (
                <TagBadge key={tag.id} tag={tag} />
              ))}
            </div>
          )}

          {/* Created/Updated */}
          <div className="pt-4 border-t border-[var(--border-color)] text-xs text-[var(--text-tertiary)] space-y-1">
            <p>Created: {formatDate(todo.created_at)}</p>
            <p>Updated: {formatDate(todo.updated_at)}</p>
          </div>
        </Card>
      </motion.div>

      {/* Sub-tasks */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.1 }}
      >
        <Card className="p-6">
          <h2 className="text-lg font-semibold text-[var(--text-primary)] mb-4">
            Sub-tasks
          </h2>
          <SubTaskList
            todoId={todoId}
            subTasks={todo.sub_tasks || []}
            isLoading={isLoading}
          />
        </Card>
      </motion.div>

      {/* Edit Modal */}
      <TodoForm
        isOpen={editModalOpen}
        onClose={() => setEditModalOpen(false)}
        onSubmit={async (data) => {
          await updateTodo(data);
          setEditModalOpen(false);
        }}
        initialData={todo}
        isLoading={isUpdating}
      />
    </div>
  );
};

export default TodoDetail;
