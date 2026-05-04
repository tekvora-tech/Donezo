import { useState, useEffect } from "react";
import { motion } from "framer-motion";
import { IoCalendar, IoFlag, IoPricetag } from "react-icons/io5";
import Input from "../ui/Input";
import Button from "../ui/Button";
import Modal from "../ui/Modal";
import TagPicker from "../tags/TagPicker";
import { PRIORITY_OPTIONS } from "../../utils/constants";
import { cn } from "../../utils/cn";

const TodoForm = ({ isOpen, onClose, onSubmit, initialData, isLoading }) => {
  const [formData, setFormData] = useState({
    title: "",
    description: "",
    priority: "medium",
    due_date: "",
    tag_ids: [],
  });
  const [showTagPicker, setShowTagPicker] = useState(false);

  useEffect(() => {
    if (initialData) {
      setFormData({
        title: initialData.title || "",
        description: initialData.description || "",
        priority: initialData.priority || "medium",
        due_date: initialData.due_date
          ? new Date(initialData.due_date).toISOString().slice(0, 16)
          : "",
        tag_ids: initialData.tags?.map((t) => t.id) || [],
      });
    } else {
      setFormData({
        title: "",
        description: "",
        priority: "medium",
        due_date: "",
        tag_ids: [],
      });
    }
  }, [initialData, isOpen]);

  const handleSubmit = (e) => {
    e.preventDefault();
    const payload = {
      ...formData,
      due_date: formData.due_date
        ? new Date(formData.due_date).toISOString()
        : undefined,
      tag_ids: formData.tag_ids.length > 0 ? formData.tag_ids : undefined,
    };
    onSubmit(payload);
  };

  const toggleTag = (tagId) => {
    setFormData((prev) => ({
      ...prev,
      tag_ids: prev.tag_ids.includes(tagId)
        ? prev.tag_ids.filter((id) => id !== tagId)
        : [...prev.tag_ids, tagId],
    }));
  };

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title={initialData ? "Edit Todo" : "New Todo"}
      size="md"
      className={""}
    >
      <form onSubmit={handleSubmit} className="space-y-4">
        <Input
          label="Title"
          placeholder="What needs to be done?"
          value={formData.title}
          onChange={(e) =>
            setFormData((prev) => ({ ...prev, title: e.target.value }))
          }
          required
        />

        <div className="space-y-1.5">
          <label className="text-sm font-medium text-[var(--text-secondary)]">
            Description
          </label>
          <textarea
            rows={3}
            placeholder="Add details..."
            value={formData.description}
            onChange={(e) =>
              setFormData((prev) => ({ ...prev, description: e.target.value }))
            }
            className="w-full px-4 py-3 rounded-xl bg-[var(--bg-secondary)] border-2 border-[var(--border-color)] text-[var(--text-primary)] placeholder:text-[var(--text-tertiary)] focus:border-pastel-blue-400 focus:ring-4 focus:ring-pastel-blue-400/20 outline-none transition-all resize-none text-sm"
          />
        </div>

        <div className="grid grid-cols-2 gap-4">
          {/* Priority */}
          <div className="space-y-1.5">
            <label className="text-sm font-medium text-[var(--text-secondary)] flex items-center gap-1.5">
              <IoFlag size={14} /> Priority
            </label>
            <div className="flex flex-wrap gap-1.5">
              {PRIORITY_OPTIONS.map((p) => (
                <button
                  key={p.value}
                  type="button"
                  onClick={() =>
                    setFormData((prev) => ({ ...prev, priority: p.value }))
                  }
                  className={cn(
                    "px-3 py-1.5 rounded-lg text-xs font-medium transition-all",
                    formData.priority === p.value
                      ? "bg-pastel-blue-400 text-white shadow-md"
                      : "bg-[var(--bg-tertiary)] text-[var(--text-secondary)] hover:bg-[var(--border-color)]",
                  )}
                >
                  {p.icon} {p.label}
                </button>
              ))}
            </div>
          </div>

          {/* Due Date */}
          <div className="space-y-1.5">
            <label className="text-sm font-medium text-[var(--text-secondary)] flex items-center gap-1.5">
              <IoCalendar size={14} /> Due Date
            </label>
            <input
              type="datetime-local"
              value={formData.due_date}
              onChange={(e) =>
                setFormData((prev) => ({ ...prev, due_date: e.target.value }))
              }
              className="w-full px-3 py-2 rounded-xl bg-[var(--bg-secondary)] border-2 border-[var(--border-color)] text-[var(--text-primary)] focus:border-pastel-blue-400 outline-none transition-all text-sm"
            />
          </div>
        </div>

        {/* Tags */}
        <div className="space-y-1.5">
          <div className="flex items-center justify-between">
            <label className="text-sm font-medium text-[var(--text-secondary)] flex items-center gap-1.5">
              <IoPricetag size={14} /> Tags
            </label>
            <button
              type="button"
              onClick={() => setShowTagPicker(!showTagPicker)}
              className="text-xs text-pastel-blue-400 hover:text-pastel-blue-500 font-medium"
            >
              {showTagPicker ? "Hide" : "Select tags"}
            </button>
          </div>

          {showTagPicker && (
            <motion.div
              initial={{ height: 0, opacity: 0 }}
              animate={{ height: "auto", opacity: 1 }}
              className="overflow-hidden"
            >
              <TagPicker selectedIds={formData.tag_ids} onToggle={toggleTag} />
            </motion.div>
          )}

          {formData.tag_ids.length > 0 && !showTagPicker && (
            <div className="flex flex-wrap gap-1.5">
              {/* Show selected tags summary */}
            </div>
          )}
        </div>

        <div className="flex gap-3 pt-2">
          <Button
            type="button"
            variant="secondary"
            onClick={onClose}
            className="flex-1"
          >
            Cancel
          </Button>
          <Button type="submit" isLoading={isLoading} className="flex-1">
            {initialData ? "Update" : "Create"} Todo
          </Button>
        </div>
      </form>
    </Modal>
  );
};

export default TodoForm;
