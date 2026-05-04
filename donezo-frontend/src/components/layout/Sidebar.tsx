import { Link, useLocation } from "react-router-dom";
import { motion, AnimatePresence } from "framer-motion";
import {
  IoGrid,
  IoCheckmarkCircle,
  IoCalendar,
  IoFlag,
  IoClose,
  IoAdd,
} from "react-icons/io5";
import { useUIStore } from "../../stores/uiStore";
import { useTags } from "../../hooks/useTags";
import { cn } from "../../utils/cn";

const navItems = [
  { path: "/dashboard", label: "Dashboard", icon: IoGrid },
  {
    path: "/dashboard?status=completed",
    label: "Completed",
    icon: IoCheckmarkCircle,
  },
  { path: "/dashboard?priority=urgent", label: "Urgent", icon: IoFlag },
  { path: "/dashboard?view=calendar", label: "Calendar", icon: IoCalendar },
];

const Sidebar = () => {
  const location = useLocation();
  const { sidebarOpen, setSidebarOpen } = useUIStore();
  const { tags, isLoading } = useTags();

  const isActive = (path) => {
    if (path.includes("?")) {
      return (
        location.pathname === path.split("?")[0] &&
        location.search === "?" + path.split("?")[1]
      );
    }
    return location.pathname === path;
  };

  const sidebarContent = (
    <div className="flex flex-col h-full">
      {/* Create New Button */}
      <div className="p-4">
        <button className="w-full btn-primary flex items-center justify-center gap-2">
          <IoAdd size={18} />
          <span>New Todo</span>
        </button>
      </div>

      {/* Main Nav */}
      <nav className="flex-1 px-3 space-y-1">
        <p className="px-3 text-xs font-semibold text-[var(--text-tertiary)] uppercase tracking-wider mb-2">
          Menu
        </p>
        {navItems.map((item) => {
          const Icon = item.icon;
          const active = isActive(item.path);
          return (
            <Link
              key={item.path}
              to={item.path}
              onClick={() => setSidebarOpen(false)}
              className={cn(
                "flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-medium transition-all duration-200",
                active
                  ? "bg-gradient-to-r from-pastel-blue-400/10 to-pastel-lavender-400/10 text-pastel-blue-500 dark:text-pastel-blue-300"
                  : "text-[var(--text-secondary)] hover:bg-[var(--bg-tertiary)] hover:text-[var(--text-primary)]",
              )}
            >
              <Icon
                size={18}
                className={active ? "text-pastel-blue-400" : ""}
              />
              {item.label}
            </Link>
          );
        })}

        {/* Tags Section */}
        <div className="mt-6">
          <p className="px-3 text-xs font-semibold text-[var(--text-tertiary)] uppercase tracking-wider mb-2">
            Tags
          </p>
          {isLoading ? (
            <div className="space-y-2 px-3">
              {[1, 2, 3].map((i) => (
                <div
                  key={i}
                  className="h-8 bg-[var(--bg-tertiary)] rounded-lg animate-pulse"
                />
              ))}
            </div>
          ) : (
            tags.map((tag) => (
              <Link
                key={tag.id}
                to={`/dashboard?tag_id=${tag.id}`}
                onClick={() => setSidebarOpen(false)}
                className="flex items-center gap-3 px-3 py-2 rounded-xl text-sm text-[var(--text-secondary)] hover:bg-[var(--bg-tertiary)] hover:text-[var(--text-primary)] transition-all duration-200"
              >
                <span
                  className="w-3 h-3 rounded-full shrink-0"
                  style={{ backgroundColor: tag.color }}
                />
                <span className="truncate">{tag.name}</span>
              </Link>
            ))
          )}
        </div>
      </nav>
    </div>
  );

  return (
    <>
      {/* Mobile overlay */}
      <AnimatePresence>
        {sidebarOpen && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            onClick={() => setSidebarOpen(false)}
            className="fixed inset-0 bg-black/30 backdrop-blur-sm z-40 lg:hidden"
          />
        )}
      </AnimatePresence>

      {/* Mobile sidebar */}
      <AnimatePresence>
        {sidebarOpen && (
          <motion.aside
            initial={{ x: -280 }}
            animate={{ x: 0 }}
            exit={{ x: -280 }}
            transition={{ type: "spring", damping: 25, stiffness: 300 }}
            className="fixed left-0 top-0 bottom-0 w-[280px] bg-[var(--bg-primary)] border-r border-[var(--border-color)] z-50 lg:hidden"
          >
            <div className="flex items-center justify-between p-4 border-b border-[var(--border-color)]">
              <span className="text-lg font-bold gradient-text">Donezo</span>
              <button
                onClick={() => setSidebarOpen(false)}
                className="p-2 rounded-lg hover:bg-[var(--bg-tertiary)]"
              >
                <IoClose size={20} />
              </button>
            </div>
            {sidebarContent}
          </motion.aside>
        )}
      </AnimatePresence>

      {/* Desktop sidebar */}
      <aside className="hidden lg:block w-64 shrink-0 border-r border-[var(--border-color)] bg-[var(--bg-secondary)]/50 min-h-[calc(100vh-64px)]">
        {sidebarContent}
      </aside>
    </>
  );
};

export default Sidebar;
