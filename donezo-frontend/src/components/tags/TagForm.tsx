import { useState } from "react";
import { motion } from "framer-motion";
import { IoColorPalette, IoAdd } from "react-icons/io5";
import Input from "../ui/Input";
import Button from "../ui/Button";
import { DEFAULT_TAG_COLORS } from "../../utils/constants";
import { cn } from "../../utils/cn";

const TagForm = ({ onSubmit, isLoading }) => {
  const [name, setName] = useState("");
  const [color, setColor] = useState(DEFAULT_TAG_COLORS[0]);

  const handleSubmit = (e) => {
    e.preventDefault();
    if (!name.trim()) return;
    onSubmit({ name: name.trim(), color });
    setName("");
    setColor(DEFAULT_TAG_COLORS[0]);
  };

  return (
    <motion.form
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      onSubmit={handleSubmit}
      className="space-y-4"
    >
      <Input
        label="Tag Name"
        placeholder="e.g. Work, Personal, Urgent"
        value={name}
        onChange={(e) => setName(e.target.value)}
        required
      />

      <div className="space-y-1.5">
        <label className="text-sm font-medium text-[var(--text-secondary)] flex items-center gap-1.5">
          <IoColorPalette size={14} /> Color
        </label>
        <div className="flex flex-wrap gap-2">
          {DEFAULT_TAG_COLORS.map((c) => (
            <button
              key={c}
              type="button"
              onClick={() => setColor(c)}
              className={cn(
                "w-8 h-8 rounded-lg transition-all",
                color === c
                  ? "ring-2 ring-offset-2 ring-[var(--text-primary)] scale-110"
                  : "hover:scale-105",
              )}
              style={{ backgroundColor: c }}
            />
          ))}
        </div>
      </div>

      <Button type="submit" isLoading={isLoading} className="w-full">
        <IoAdd size={18} /> Create Tag
      </Button>
    </motion.form>
  );
};

export default TagForm;
