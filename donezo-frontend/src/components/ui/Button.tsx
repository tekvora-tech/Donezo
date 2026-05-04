import { forwardRef } from "react";
import { motion, type HTMLMotionProps } from "framer-motion";
import { cn } from "../../utils/cn";

type ButtonVariant = "primary" | "secondary" | "ghost" | "danger" | "outline";
type ButtonSize = "sm" | "md" | "lg" | "icon";

type ButtonProps = Omit<HTMLMotionProps<"button">, "children"> & {
  children?: React.ReactNode;
  variant?: ButtonVariant;
  size?: ButtonSize;
  isLoading?: boolean;
};

const variants = {
  primary:
    "bg-gradient-to-r from-pastel-blue-400 to-pastel-lavender-400 text-white shadow-lg shadow-pastel-blue-400/25 hover:shadow-pastel-blue-400/40",
  secondary:
    "bg-[var(--bg-tertiary)] text-[var(--text-primary)] hover:bg-[var(--border-color)]",
  ghost:
    "bg-transparent text-[var(--text-secondary)] hover:bg-[var(--bg-tertiary)] hover:text-[var(--text-primary)]",
  danger:
    "bg-gradient-to-r from-pastel-rose-400 to-pastel-peach-400 text-white shadow-lg shadow-pastel-rose-400/25 hover:shadow-pastel-rose-400/40",
  outline:
    "bg-transparent border-2 border-[var(--border-color)] text-[var(--text-primary)] hover:border-pastel-blue-400 hover:text-pastel-blue-500",
} as const;

const sizes = {
  sm: "px-3 py-1.5 text-sm rounded-lg",
  md: "px-4 py-2.5 text-sm rounded-xl",
  lg: "px-6 py-3 text-base rounded-xl",
  icon: "p-2 rounded-xl",
} as const;

const Button = forwardRef<HTMLButtonElement, ButtonProps>(
  (
    {
      children,
      variant = "primary",
      size = "md",
      isLoading = false,
      disabled,
      className,
      ...rest
    },
    ref,
  ) => {
    return (
      <motion.button
        ref={ref}
        whileTap={{ scale: 0.97 }}
        whileHover={{ scale: 1.02 }}
        className={cn(
          "inline-flex items-center justify-center gap-2 font-medium transition-all duration-200",
          "disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:scale-100",
          variants[variant],
          sizes[size],
          className,
        )}
        disabled={disabled || isLoading}
        {...rest}
      >
        {isLoading && (
          <motion.div
            animate={{ rotate: 360 }}
            transition={{ repeat: Infinity, duration: 1, ease: "linear" }}
            className="w-4 h-4 border-2 border-current border-t-transparent rounded-full"
          />
        )}
        {children}
      </motion.button>
    );
  },
);

Button.displayName = "Button";

export default Button;
