import { cn } from "../../utils/cn";
import { motion } from "framer-motion";

const Skeleton = ({ className, count = 1 }) => {
  return (
    <>
      {Array.from({ length: count }).map((_, i) => (
        <motion.div
          key={i}
          initial={{ opacity: 0 }}
          animate={{ opacity: [0.5, 1, 0.5] }}
          transition={{ repeat: Infinity, duration: 1.5, ease: "easeInOut" }}
          className={cn(
            "bg-[var(--bg-tertiary)] rounded-xl animate-pulse",
            className,
          )}
        />
      ))}
    </>
  );
};

export default Skeleton;
