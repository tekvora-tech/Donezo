import { useEffect } from "react";
import { motion, AnimatePresence } from "framer-motion";
import {
  IoCheckmarkCircle,
  IoCloseCircle,
  IoInformationCircle,
  IoWarning,
  IoClose,
} from "react-icons/io5";
import { useUIStore } from "../../stores/uiStore";
import { cn } from "../../utils/cn";

const toastIcons = {
  success: IoCheckmarkCircle,
  error: IoCloseCircle,
  info: IoInformationCircle,
  warning: IoWarning,
};

const toastStyles = {
  success:
    "bg-pastel-mint-50 border-pastel-mint-300 text-pastel-mint-700 dark:bg-pastel-mint-900/20 dark:border-pastel-mint-700/50 dark:text-pastel-mint-300",
  error:
    "bg-pastel-rose-50 border-pastel-rose-300 text-pastel-rose-700 dark:bg-pastel-rose-900/20 dark:border-pastel-rose-700/50 dark:text-pastel-rose-300",
  info: "bg-pastel-blue-50 border-pastel-blue-300 text-pastel-blue-700 dark:bg-pastel-blue-900/20 dark:border-pastel-blue-700/50 dark:text-pastel-blue-300",
  warning:
    "bg-pastel-lemon-50 border-pastel-lemon-300 text-pastel-lemon-700 dark:bg-pastel-lemon-900/20 dark:border-pastel-lemon-700/50 dark:text-pastel-lemon-300",
};

const ToastItem = ({ toast }) => {
  const { removeToast } = useUIStore();
  const Icon = toastIcons[toast.type];

  useEffect(() => {
    const timer = setTimeout(() => removeToast(toast.id), 4000);
    return () => clearTimeout(timer);
  }, [toast.id, removeToast]);

  return (
    <motion.div
      layout
      initial={{ opacity: 0, x: 100, scale: 0.9 }}
      animate={{ opacity: 1, x: 0, scale: 1 }}
      exit={{ opacity: 0, x: 100, scale: 0.9 }}
      className={cn(
        "flex items-center gap-3 px-4 py-3 rounded-xl border shadow-lg min-w-[320px] max-w-md",
        toastStyles[toast.type],
      )}
    >
      <Icon size={20} className="shrink-0" />
      <p className="text-sm font-medium flex-1">{toast.message}</p>
      <button
        onClick={() => removeToast(toast.id)}
        className="p-1 rounded-lg hover:bg-black/5 dark:hover:bg-white/10 transition-colors"
      >
        <IoClose size={16} />
      </button>
    </motion.div>
  );
};

const ToastContainer = () => {
  const { toasts } = useUIStore();

  return (
    <div className="fixed top-4 right-4 z-[100] flex flex-col gap-2 pointer-events-none">
      <AnimatePresence mode="popLayout">
        {toasts.map((toast) => (
          <div key={toast.id} className="pointer-events-auto">
            <ToastItem toast={toast} />
          </div>
        ))}
      </AnimatePresence>
    </div>
  );
};

export default ToastContainer;
