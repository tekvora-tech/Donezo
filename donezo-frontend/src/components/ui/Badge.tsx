import { cn } from "../../utils/cn";

const Badge = ({ children, variant = "default", className, ...props }) => {
  const variants = {
    default:
      "bg-pastel-blue-100 text-pastel-blue-600 dark:bg-pastel-blue-900/30 dark:text-pastel-blue-300",
    success:
      "bg-pastel-mint-100 text-pastel-mint-600 dark:bg-pastel-mint-900/30 dark:text-pastel-mint-300",
    warning:
      "bg-pastel-lemon-100 text-pastel-lemon-600 dark:bg-pastel-lemon-900/30 dark:text-pastel-lemon-300",
    danger:
      "bg-pastel-rose-100 text-pastel-rose-600 dark:bg-pastel-rose-900/30 dark:text-pastel-rose-300",
    info: "bg-pastel-lavender-100 text-pastel-lavender-600 dark:bg-pastel-lavender-900/30 dark:text-pastel-lavender-300",
    neutral: "bg-[var(--bg-tertiary)] text-[var(--text-secondary)]",
  };

  return (
    <span
      className={cn(
        "inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-xs font-medium",
        variants[variant],
        className,
      )}
      {...props}
    >
      {children}
    </span>
  );
};

export default Badge;
