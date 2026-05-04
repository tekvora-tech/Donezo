import { forwardRef } from "react";
import type { InputHTMLAttributes } from "react";
import { motion } from "framer-motion";
import { cn } from "../../utils/cn";

type InputProps = Omit<InputHTMLAttributes<HTMLInputElement>, "size"> & {
  label?: string;
  error?: string;
  icon?: React.ComponentType<{ size?: number }>;
  containerClassName?: string;
};

const Input = forwardRef<HTMLInputElement, InputProps>(
  (
    { label, error, icon: Icon, className, containerClassName, ...props },
    ref,
  ) => {
    return (
      <div className={cn("space-y-1.5", containerClassName)}>
        {label && (
          <label className="block text-sm font-medium text-[var(--text-secondary)]">
            {label}
          </label>
        )}

        <div className="relative">
          {Icon && (
            <div className="absolute left-3.5 top-1/2 -translate-y-1/2 text-[var(--text-tertiary)]">
              <Icon size={18} />
            </div>
          )}

          <input
            ref={ref}
            className={cn(
              "w-full px-4 py-3 rounded-xl bg-[var(--bg-secondary)] border-2",
              "text-[var(--text-primary)] placeholder:text-[var(--text-tertiary)]",
              "focus:border-pastel-blue-400 focus:ring-4 focus:ring-pastel-blue-400/20",
              "outline-none transition-all duration-200",
              Icon && "pl-11",
              error &&
                "border-pastel-rose-400 focus:border-pastel-rose-400 focus:ring-pastel-rose-400/20",
              className,
            )}
            {...props}
          />
        </div>

        {error && (
          <motion.p
            initial={{ opacity: 0, y: -5 }}
            animate={{ opacity: 1, y: 0 }}
            className="text-sm text-pastel-rose-500"
          >
            {error}
          </motion.p>
        )}
      </div>
    );
  },
);

Input.displayName = "Input";

export default Input;
