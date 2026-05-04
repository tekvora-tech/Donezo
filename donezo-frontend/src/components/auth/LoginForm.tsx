import { useState } from "react";
import { Link } from "react-router-dom";
import { motion } from "framer-motion";
import { IoMail, IoLockClosed, IoEye, IoEyeOff } from "react-icons/io5";
import Input from "../ui/Input";
import Button from "../ui/Button";
import { useAuth } from "../../hooks/useAuth";

const LoginForm = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const { login, isLoggingIn } = useAuth();

  const handleSubmit = async (e) => {
    e.preventDefault();
    await login({ email, password });
  };

  return (
    <motion.form
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay: 0.2 }}
      onSubmit={handleSubmit}
      className="space-y-5"
    >
      <Input
        label="Email"
        type="email"
        placeholder="john@example.com"
        value={email}
        onChange={(e) => setEmail(e.target.value)}
        icon={IoMail}
        required
      />

      <div className="relative">
        <Input
          label="Password"
          type={showPassword ? "text" : "password"}
          placeholder="••••••••"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          icon={IoLockClosed}
          required
        />
        <button
          type="button"
          onClick={() => setShowPassword(!showPassword)}
          className="absolute right-3.5 top-[38px] text-[var(--text-tertiary)] hover:text-[var(--text-secondary)] transition-colors"
        >
          {showPassword ? <IoEyeOff size={18} /> : <IoEye size={18} />}
        </button>
      </div>

      <div className="flex items-center justify-between text-sm">
        <label className="flex items-center gap-2 text-[var(--text-secondary)] cursor-pointer">
          <input
            type="checkbox"
            className="rounded border-[var(--border-color)] text-pastel-blue-400 focus:ring-pastel-blue-400"
          />
          Remember me
        </label>
        <Link
          to="/forgot-password"
          className="text-pastel-blue-400 hover:text-pastel-blue-500 font-medium"
        >
          Forgot password?
        </Link>
      </div>

      <Button type="submit" isLoading={isLoggingIn} className="w-full">
        Sign In
      </Button>

      <p className="text-center text-sm text-[var(--text-secondary)]">
        Don't have an account?{" "}
        <Link
          to="/register"
          className="text-pastel-blue-400 hover:text-pastel-blue-500 font-medium"
        >
          Sign up
        </Link>
      </p>
    </motion.form>
  );
};

export default LoginForm;
