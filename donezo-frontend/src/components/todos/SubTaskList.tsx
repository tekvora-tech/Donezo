import { useState } from "react";
import { motion, AnimatePresence } from "framer-motion";
import {
  IoAdd,
  IoCheckmarkCircle,
  IoEllipseOutline,
  IoTrash,
} from "react-icons/io5";
import { useTodoDetail } from "../../hooks/useTodos";
import { cn } from "../../utils/cn";
import Button from "../../components/ui/Button";

const SubTaskList = ({ todoId, subTasks, isLoading }) => {
  const [newTitle, setNewTitle] = useState("");
  const { createSubTask, updateSubTask, deleteSubTask, isCreatingSubTask } =
    useTodoDetail(todoId);

  const handleAdd = async (e) => {
    e.preventDefault();
    if (!newTitle.trim()) return;
    await createSubTask({ title: newTitle.trim() });
    setNewTitle("");
  };

  const toggleComplete = async (subTask) => {
    await updateSubTask({
      subTaskId: subTask.id,
      data: { is_completed: !subTask.is_completed },
    });
  };

  const progress = subTasks.length
    ? Math.round(
        (subTasks.filter((s) => s.is_completed).length / subTasks.length) * 100,
      )
    : 0;

  return (
    <div className="space-y-4">
      {/* Progress */}
      {subTasks.length > 0 && (
        <div className="space-y-1.5">
          <div className="flex items-center justify-between text-sm">
            <span className="text-[var(--text-secondary)] font-medium">
              Progress
            </span>
            <span className="text-pastel-blue-400 font-bold">{progress}%</span>
          </div>
          <div className="h-2 bg-[var(--bg-tertiary)] rounded-full overflow-hidden">
            <motion.div
              initial={{ width: 0 }}
              animate={{ width: `${progress}%` }}
              transition={{ duration: 0.5 }}
              className={cn(
                "h-full rounded-full",
                progress === 100
                  ? "bg-pastel-mint-400"
                  : "bg-gradient-to-r from-pastel-blue-400 to-pastel-lavender-400",
              )}
            />
          </div>
        </div>
      )}

      {/* Add new */}
      <form onSubmit={handleAdd} className="flex gap-2">
        <input
          type="text"
          placeholder="Add a sub-task..."
          value={newTitle}
          onChange={(e) => setNewTitle(e.target.value)}
          className="flex-1 px-4 py-2.5 rounded-xl bg-[var(--bg-secondary)] border border-[var(--border-color)] text-[var(--text-primary)] placeholder:text-[var(--text-tertiary)] focus:border-pastel-blue-400 focus:ring-2 focus:ring-pastel-blue-400/20 outline-none transition-all text-sm"
        />
        <Button
          type="submit"
          variant="primary"
          size="icon"
          isLoading={isCreatingSubTask}
        >
          <IoAdd size={18} />
        </Button>
      </form>

      {/* List */}
      <div className="space-y-1">
        <AnimatePresence mode="popLayout">
          {subTasks.map((subTask) => (
            <motion.div
              key={subTask.id}
              layout
              initial={{ opacity: 0, x: -20 }}
              animate={{ opacity: 1, x: 0 }}
              exit={{ opacity: 0, x: 20 }}
              className="group flex items-center gap-3 p-3 rounded-xl hover:bg-[var(--bg-tertiary)] transition-colors"
            >
              <button
                onClick={() => toggleComplete(subTask)}
                className={cn(
                  "shrink-0 transition-colors",
                  subTask.is_completed
                    ? "text-pastel-mint-400"
                    : "text-[var(--text-tertiary)] hover:text-pastel-blue-400",
                )}
              >
                {subTask.is_completed ? (
                  <IoCheckmarkCircle size={20} />
                ) : (
                  <IoEllipseOutline size={20} />
                )}
              </button>

              <span
                className={cn(
                  "flex-1 text-sm transition-all",
                  subTask.is_completed
                    ? "line-through text-[var(--text-tertiary)]"
                    : "text-[var(--text-primary)]",
                )}
              >
                {subTask.title}
              </span>

              <button
                onClick={() => deleteSubTask(subTask.id)}
                className="opacity-0 group-hover:opacity-100 p-1.5 rounded-lg hover:bg-pastel-rose-50 dark:hover:bg-pastel-rose-900/20 text-[var(--text-tertiary)] hover:text-pastel-rose-400 transition-all"
              >
                <IoTrash size={14} />
              </button>
            </motion.div>
          ))}
        </AnimatePresence>

        {subTasks.length === 0 && !isLoading && (
          <div className="text-center py-8 text-[var(--text-tertiary)] text-sm">
            No sub-tasks yet. Add one above!
          </div>
        )}
      </div>
    </div>
  );
};

export default SubTaskList;
