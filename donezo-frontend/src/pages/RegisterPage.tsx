import { Link } from "react-router-dom";
import { motion } from "framer-motion";
import { IoArrowBack } from "react-icons/io5";
import RegisterForm from "../components/auth/RegisterForm";

const RegisterPage = () => {
  return (
    <div className="min-h-screen bg-[var(--bg-primary)] flex items-center justify-center p-4">
      <div className="absolute inset-0 overflow-hidden pointer-events-none">
        <div className="absolute -top-40 -right-40 w-96 h-96 bg-pastel-mint-400/10 rounded-full blur-3xl" />
        <div className="absolute -bottom-40 -left-40 w-96 h-96 bg-pastel-peach-400/10 rounded-full blur-3xl" />
      </div>

      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        className="relative w-full max-w-md"
      >
        <Link
          to="/"
          className="inline-flex items-center gap-2 text-sm text-[var(--text-tertiary)] hover:text-[var(--text-primary)] mb-8 transition-colors"
        >
          <IoArrowBack size={16} /> Back to home
        </Link>

        <div className="bg-[var(--bg-secondary)] rounded-3xl border border-[var(--border-color)] shadow-xl p-8">
          <div className="text-center mb-8">
            <div className="w-14 h-14 rounded-2xl bg-gradient-to-br from-pastel-mint-400 to-pastel-blue-400 flex items-center justify-center mx-auto mb-4 shadow-lg shadow-pastel-mint-400/30">
              <span className="text-white font-bold text-2xl">D</span>
            </div>
            <h1 className="text-2xl font-bold text-[var(--text-primary)]">
              Create account
            </h1>
            <p className="text-sm text-[var(--text-secondary)] mt-1">
              Start organizing your tasks today
            </p>
          </div>

          <RegisterForm />
        </div>
      </motion.div>
    </div>
  );
};

export default RegisterPage;
