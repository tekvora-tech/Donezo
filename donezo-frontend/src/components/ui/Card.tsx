import { cn } from "../../utils/cn";
import { motion } from "framer-motion";

const Card = ({ children, className, hover = false, ...props }) => {
  return (
    <motion.div
      initial={{ opacity: 0, y: 10 }}
      animate={{ opacity: 1, y: 0 }}
      whileHover={hover ? { y: -4, transition: { duration: 0.2 } } : undefined}
      className={cn(
        "bg-[var(--bg-secondary)] rounded-2xl border border-[var(--border-color)]",
        "shadow-sm transition-all duration-300",
        hover &&
          "hover:shadow-lg hover:shadow-pastel-blue-400/10 hover:border-pastel-blue-300/50",
        className,
      )}
      {...props}
    >
      {children}
    </motion.div>
  );
};

export default Card;
