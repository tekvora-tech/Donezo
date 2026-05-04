import { useTags } from "../../hooks/useTags";
import { cn } from "../../utils/cn";
import { IoCheckmark } from "react-icons/io5";
import Skeleton from "../ui/Skeleton";

const TagPicker = ({ selectedIds, onToggle }) => {
  const { tags, isLoading } = useTags();

  if (isLoading) {
    return <Skeleton count={3} className="h-8" />;
  }

  if (tags.length === 0) {
    return (
      <p className="text-sm text-[var(--text-tertiary)] py-2">
        No tags yet. Create some tags first!
      </p>
    );
  }

  return (
    <div className="flex flex-wrap gap-2 p-3 rounded-xl bg-[var(--bg-tertiary)] border border-[var(--border-color)]">
      {tags.map((tag) => {
        const isSelected = selectedIds.includes(tag.id);
        return (
          <button
            key={tag.id}
            type="button"
            onClick={() => onToggle(tag.id)}
            className={cn(
              "flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-sm font-medium transition-all",
              isSelected
                ? "ring-2 ring-offset-1"
                : "bg-[var(--bg-primary)] hover:bg-[var(--bg-secondary)]",
            )}
            style={
              isSelected
                ? {
                    backgroundColor: `${tag.color}20`,
                    color: tag.color,
                    boxShadow: `0 0 0 2px ${tag.color}40`,
                  }
                : { color: tag.color }
            }
          >
            <span
              className="w-2.5 h-2.5 rounded-full"
              style={{ backgroundColor: tag.color }}
            />
            {tag.name}
            {isSelected && <IoCheckmark size={14} />}
          </button>
        );
      })}
    </div>
  );
};

export default TagPicker;
