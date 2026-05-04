import { Link, useNavigate } from "react-router-dom";
import { motion, AnimatePresence } from "framer-motion";
import {
  IoMoon,
  IoSunny,
  IoMenu,
  IoLogOut,
  IoPerson,
  IoChevronDown,
} from "react-icons/io5";
import { useTheme } from "../../hooks/useTheme";
import { useAuth } from "../../hooks/useAuth";
import { useUIStore } from "../../stores/uiStore";
import { useState, useRef, useEffect } from "react";
import { cn } from "../../utils/cn";

const Navbar = () => {
  const { isDark, toggleTheme } = useTheme();
  const { user, logout, isLoggingOut } = useAuth();
  const { toggleSidebar } = useUIStore();
  const navigate = useNavigate();
  const [dropdownOpen, setDropdownOpen] = useState(false);
  const dropdownRef = useRef(null);

  useEffect(() => {
    const handleClickOutside = (e) => {
      if (dropdownRef.current && !dropdownRef.current.contains(e.target)) {
        setDropdownOpen(false);
      }
    };
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  const handleLogout = async () => {
    await logout();
    navigate("/login");
  };

  return (
    <motion.nav
      initial={{ y: -20, opacity: 0 }}
      animate={{ y: 0, opacity: 1 }}
      className="sticky top-0 z-40 glass-card border-b border-[var(--border-color)]"
    >
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between h-16">
          {/* Left */}
          <div className="flex items-center gap-3">
            <button
              onClick={toggleSidebar}
              className="p-2 rounded-xl hover:bg-[var(--bg-tertiary)] text-[var(--text-secondary)] lg:hidden transition-colors"
            >
              <IoMenu size={22} />
            </button>
            <Link to="/dashboard" className="flex items-center gap-2.5">
              <div className="w-9 h-9 rounded-xl bg-gradient-to-br from-pastel-blue-400 to-pastel-lavender-400 flex items-center justify-center shadow-lg shadow-pastel-blue-400/30">
                <span className="text-white font-bold text-lg">D</span>
              </div>
              <span className="text-xl font-bold gradient-text hidden sm:block">
                Donezo
              </span>
            </Link>
          </div>

          {/* Right */}
          <div className="flex items-center gap-2">
            {/* Theme Toggle */}
            <motion.button
              whileTap={{ scale: 0.9 }}
              onClick={toggleTheme}
              className="p-2.5 rounded-xl hover:bg-[var(--bg-tertiary)] text-[var(--text-secondary)] transition-colors relative"
            >
              <motion.div
                key={isDark ? "dark" : "light"}
                initial={{ rotate: -90, opacity: 0 }}
                animate={{ rotate: 0, opacity: 1 }}
                exit={{ rotate: 90, opacity: 0 }}
                transition={{ duration: 0.2 }}
              >
                {isDark ? <IoSunny size={20} /> : <IoMoon size={20} />}
              </motion.div>
            </motion.button>

            {/* User Dropdown */}
            {user && (
              <div className="relative" ref={dropdownRef}>
                <button
                  onClick={() => setDropdownOpen(!dropdownOpen)}
                  className="flex items-center gap-2.5 p-1.5 pr-3 rounded-xl hover:bg-[var(--bg-tertiary)] transition-colors"
                >
                  <div className="w-8 h-8 rounded-full bg-gradient-to-br from-pastel-mint-400 to-pastel-blue-400 flex items-center justify-center text-white text-sm font-bold">
                    {user.full_name?.charAt(0).toUpperCase() || "U"}
                  </div>
                  <span className="text-sm font-medium text-[var(--text-primary)] hidden sm:block">
                    {user.full_name}
                  </span>
                  <IoChevronDown
                    size={14}
                    className={cn(
                      "text-[var(--text-tertiary)] transition-transform",
                      dropdownOpen && "rotate-180",
                    )}
                  />
                </button>

                <AnimatePresence>
                  {dropdownOpen && (
                    <motion.div
                      initial={{ opacity: 0, y: 8, scale: 0.95 }}
                      animate={{ opacity: 1, y: 0, scale: 1 }}
                      exit={{ opacity: 0, y: 8, scale: 0.95 }}
                      transition={{ duration: 0.15 }}
                      className="absolute right-0 mt-2 w-48 bg-[var(--bg-primary)] rounded-xl border border-[var(--border-color)] shadow-xl py-1.5"
                    >
                      <Link
                        to="/profile"
                        onClick={() => setDropdownOpen(false)}
                        className="flex items-center gap-2.5 px-4 py-2.5 text-sm text-[var(--text-secondary)] hover:bg-[var(--bg-tertiary)] hover:text-[var(--text-primary)] transition-colors"
                      >
                        <IoPerson size={16} />
                        Profile
                      </Link>
                      <div className="mx-3 my-1.5 h-px bg-[var(--border-color)]" />
                      <button
                        onClick={handleLogout}
                        disabled={isLoggingOut}
                        className="w-full flex items-center gap-2.5 px-4 py-2.5 text-sm text-pastel-rose-500 hover:bg-pastel-rose-50 dark:hover:bg-pastel-rose-900/20 transition-colors"
                      >
                        <IoLogOut size={16} />
                        {isLoggingOut ? "Logging out..." : "Logout"}
                      </button>
                    </motion.div>
                  )}
                </AnimatePresence>
              </div>
            )}
          </div>
        </div>
      </div>
    </motion.nav>
  );
};

export default Navbar;
