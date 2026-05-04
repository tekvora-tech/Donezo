import { Link } from "react-router-dom";
import { motion } from "framer-motion";
import {
  IoCheckmarkCircle,
  IoFlash,
  IoShield,
  IoPeople,
  IoArrowForward,
  IoSparkles,
} from "react-icons/io5";
import { useTheme } from "../hooks/useTheme";
import { cn } from "../utils/cn";

const features = [
  {
    icon: IoCheckmarkCircle,
    title: "Smart Task Management",
    desc: "Organize tasks with priorities, due dates, and sub-tasks effortlessly.",
    color: "from-pastel-mint-400 to-pastel-blue-400",
    bg: "bg-pastel-mint-50 dark:bg-pastel-mint-900/10",
  },
  {
    icon: IoFlash,
    title: "Lightning Fast",
    desc: "Built with modern tech for instant loading and smooth interactions.",
    color: "from-pastel-lemon-400 to-pastel-peach-400",
    bg: "bg-pastel-lemon-50 dark:bg-pastel-lemon-900/10",
  },
  {
    icon: IoShield,
    title: "Secure & Private",
    desc: "Your data is protected with JWT authentication and encrypted storage.",
    color: "from-pastel-blue-400 to-pastel-lavender-400",
    bg: "bg-pastel-blue-50 dark:bg-pastel-blue-900/10",
  },
  {
    icon: IoPeople,
    title: "Team Ready",
    desc: "Built for individuals now, ready for team collaboration in the future.",
    color: "from-pastel-lavender-400 to-pastel-rose-400",
    bg: "bg-pastel-lavender-50 dark:bg-pastel-lavender-900/10",
  },
];

const LandingPage = () => {
  const { isDark } = useTheme();

  return (
    <div className="min-h-screen bg-[var(--bg-primary)]">
      {/* Hero */}
      <section className="relative overflow-hidden pt-20 pb-32 px-4">
        {/* Background decorations */}
        <div className="absolute inset-0 overflow-hidden pointer-events-none">
          <div className="absolute -top-40 -right-40 w-96 h-96 bg-pastel-blue-400/10 rounded-full blur-3xl" />
          <div className="absolute -bottom-40 -left-40 w-96 h-96 bg-pastel-lavender-400/10 rounded-full blur-3xl" />
          <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[600px] h-[600px] bg-pastel-mint-400/5 rounded-full blur-3xl" />
        </div>

        <div className="relative max-w-5xl mx-auto text-center">
          <motion.div
            initial={{ opacity: 0, y: 30 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.6 }}
          >
            <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-pastel-blue-50 dark:bg-pastel-blue-900/20 text-pastel-blue-500 text-sm font-medium mb-8">
              <IoSparkles size={16} />
              Your personal productivity companion
            </div>

            <h1 className="text-5xl sm:text-6xl lg:text-7xl font-bold leading-tight mb-6">
              Get things <span className="gradient-text">done</span>
              <br />
              with elegance
            </h1>

            <p className="text-lg sm:text-xl text-[var(--text-secondary)] max-w-2xl mx-auto mb-10 leading-relaxed">
              Donezo helps you organize, track, and complete your daily tasks
              with a beautiful, intuitive interface designed for modern
              productivity.
            </p>

            <div className="flex flex-wrap items-center justify-center gap-4">
              <Link
                to="/register"
                className="btn-primary text-base px-8 py-4 flex items-center gap-2"
              >
                Get Started Free
                <IoArrowForward size={18} />
              </Link>
              <Link to="/login" className="btn-secondary text-base px-8 py-4">
                Sign In
              </Link>
            </div>
          </motion.div>

          {/* App Preview */}
          <motion.div
            initial={{ opacity: 0, y: 60 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.3, duration: 0.8 }}
            className="mt-16 relative"
          >
            <div className="relative rounded-2xl overflow-hidden shadow-2xl shadow-pastel-blue-400/10 border border-[var(--border-color)] bg-[var(--bg-secondary)] max-w-4xl mx-auto">
              <div className="p-4 border-b border-[var(--border-color)] flex items-center gap-2">
                <div className="w-3 h-3 rounded-full bg-pastel-rose-400" />
                <div className="w-3 h-3 rounded-full bg-pastel-lemon-400" />
                <div className="w-3 h-3 rounded-full bg-pastel-mint-400" />
              </div>
              <div className="p-8 space-y-4">
                {[1, 2, 3].map((i) => (
                  <div
                    key={i}
                    className="flex items-center gap-4 p-4 rounded-xl bg-[var(--bg-primary)] border border-[var(--border-color)]"
                  >
                    <div
                      className={cn(
                        "w-5 h-5 rounded-full border-2",
                        i === 1
                          ? "bg-pastel-mint-400 border-pastel-mint-400"
                          : "border-[var(--border-color)]",
                      )}
                    />
                    <div className="flex-1 space-y-2">
                      <div
                        className={cn(
                          "h-4 rounded w-3/4",
                          i === 1
                            ? "bg-[var(--text-tertiary)] opacity-30"
                            : "bg-[var(--text-primary)]",
                        )}
                      />
                      <div className="h-3 rounded w-1/2 bg-[var(--text-tertiary)] opacity-20" />
                    </div>
                    <div className="flex gap-1.5">
                      <div className="w-16 h-6 rounded-full bg-pastel-blue-100 dark:bg-pastel-blue-900/30" />
                      <div className="w-12 h-6 rounded-full bg-pastel-mint-100 dark:bg-pastel-mint-900/30" />
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </motion.div>
        </div>
      </section>

      {/* Features */}
      <section className="py-24 px-4 bg-[var(--bg-secondary)]">
        <div className="max-w-6xl mx-auto">
          <motion.div
            initial={{ opacity: 0 }}
            whileInView={{ opacity: 1 }}
            viewport={{ once: true }}
            className="text-center mb-16"
          >
            <h2 className="text-3xl sm:text-4xl font-bold text-[var(--text-primary)] mb-4">
              Everything you need
            </h2>
            <p className="text-[var(--text-secondary)] text-lg max-w-xl mx-auto">
              Powerful features wrapped in a delightful, pastel-colored
              experience.
            </p>
          </motion.div>

          <div className="grid sm:grid-cols-2 lg:grid-cols-4 gap-6">
            {features.map((feature, i) => {
              const Icon = feature.icon;
              return (
                <motion.div
                  key={feature.title}
                  initial={{ opacity: 0, y: 30 }}
                  whileInView={{ opacity: 1, y: 0 }}
                  viewport={{ once: true }}
                  transition={{ delay: i * 0.1 }}
                  whileHover={{ y: -8 }}
                  className={cn(
                    "p-6 rounded-2xl border border-[var(--border-color)] transition-all duration-300 hover:shadow-xl",
                    feature.bg,
                  )}
                >
                  <div
                    className={cn(
                      "w-12 h-12 rounded-xl bg-gradient-to-br flex items-center justify-center text-white mb-4 shadow-lg",
                      feature.color,
                    )}
                  >
                    <Icon size={24} />
                  </div>
                  <h3 className="text-lg font-semibold text-[var(--text-primary)] mb-2">
                    {feature.title}
                  </h3>
                  <p className="text-sm text-[var(--text-secondary)] leading-relaxed">
                    {feature.desc}
                  </p>
                </motion.div>
              );
            })}
          </div>
        </div>
      </section>

      {/* CTA */}
      <section className="py-24 px-4">
        <div className="max-w-3xl mx-auto text-center">
          <motion.div
            initial={{ opacity: 0, scale: 0.95 }}
            whileInView={{ opacity: 1, scale: 1 }}
            viewport={{ once: true }}
            className="p-12 rounded-3xl bg-gradient-to-br from-pastel-blue-400/10 via-pastel-lavender-400/10 to-pastel-mint-400/10 border border-[var(--border-color)]"
          >
            <h2 className="text-3xl font-bold text-[var(--text-primary)] mb-4">
              Ready to boost your productivity?
            </h2>
            <p className="text-[var(--text-secondary)] mb-8 max-w-lg mx-auto">
              Join thousands of users who trust Donezo to manage their daily
              tasks efficiently.
            </p>
            <Link
              to="/register"
              className="btn-primary text-base px-8 py-4 inline-flex items-center gap-2"
            >
              Start for Free
              <IoArrowForward size={18} />
            </Link>
          </motion.div>
        </div>
      </section>

      {/* Footer */}
      <footer className="py-8 px-4 border-t border-[var(--border-color)] text-center text-sm text-[var(--text-tertiary)]">
        <p>© 2026 Donezo. Built with 💙 for productivity enthusiasts.</p>
      </footer>
    </div>
  );
};

export default LandingPage;
