import { cn } from "../../utils/cn";
import { PRIORITY_COLORS } from "../../utils/constants";

const PriorityBadge = ({ priority, showLabel = true, size = "md" }) => {
  const colors = PRIORITY_COLORS[priority] || PRIORITY_COLORS.medium;

  const sizes = {
    sm: "px-2 py-0.5 text-xs",
    md: "px-2.5 py-1 text-sm",
  };

  const icons = {
    low: "↓",
    medium: "→",
    high: "↑",
    urgent: "!",
  };

  return (
    <span
      className={cn(
        "inline-flex items-center gap-1 rounded-full font-medium border",
        colors.bg,
        colors.text,
        colors.border,
        sizes[size],
      )}
    >
      <span className="font-bold">{icons[priority]}</span>
      {showLabel && <span className="capitalize">{priority}</span>}
    </span>
  );
};

export default PriorityBadge;
