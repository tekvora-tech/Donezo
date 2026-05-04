import { cn } from "../../utils/cn";
import { IoClose } from "react-icons/io5";

type Tag = {
  id: string;
  name: string;
  color: string;
};

type TagBadgeProps = {
  tag: Tag;
  size?: "sm" | "md";
} & (
  | {
      removable: true;
      onRemove: (id: string) => void;
    }
  | {
      removable?: false;
      onRemove?: never;
    }
);

const TagBadge = ({
  tag,
  onRemove,
  removable = false,
  size = "md",
}: TagBadgeProps) => {
  const sizes = {
    sm: "px-2 py-0.5 text-xs gap-1",
    md: "px-2.5 py-1 text-sm gap-1.5",
  } as const;

  return (
    <span
      className={cn(
        "inline-flex items-center rounded-full font-medium",
        "bg-opacity-15 dark:bg-opacity-25",
        sizes[size],
      )}
      style={{
        backgroundColor: `${tag.color}25`,
        color: tag.color,
        border: `1px solid ${tag.color}40`,
      }}
    >
      <span
        className="w-2 h-2 rounded-full shrink-0"
        style={{ backgroundColor: tag.color }}
      />

      {tag.name}

      {removable && (
        <button
          onClick={(e) => {
            e.stopPropagation();
            onRemove(tag.id);
          }}
          className="hover:opacity-70 transition-opacity"
        >
          <IoClose size={14} />
        </button>
      )}
    </span>
  );
};

export default TagBadge;
