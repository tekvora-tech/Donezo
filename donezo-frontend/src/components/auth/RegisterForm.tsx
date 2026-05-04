import { useState } from "react";
import { Link } from "react-router-dom";
import { motion } from "framer-motion";
import {
  IoMail,
  IoLockClosed,
  IoPerson,
  IoEye,
  IoEyeOff,
} from "react-icons/io5";
import Input from "../ui/Input";
import Button from "../ui/Button";
import { useAuth } from "../../hooks/useAuth";

const RegisterForm = () => {
  const [formData, setFormData] = useState({
    full_name: "",
    email: "",
    password: "",
    confirmPassword: "",
  });
  const [showPassword, setShowPassword] = useState(false);
  const [error, setError] = useState("");
  const { register, isRegistering } = useAuth();

  const handleChange = (e) => {
    setFormData((prev) => ({ ...prev, [e.target.name]: e.target.value }));
    setError("");
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (formData.password !== formData.confirmPassword) {
      setError("Passwords do not match");
      return;
    }
    if (formData.password.length < 8) {
      setError("Password must be at least 8 characters");
      return;
    }
    await register({
      email: formData.email,
      password: formData.password,
      full_name: formData.full_name,
    });
  };

  return (
    <motion.form
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay: 0.2 }}
      onSubmit={handleSubmit}
      className="space-y-4"
    >
      <Input
        label="Full Name"
        name="full_name"
        placeholder="John Doe"
        value={formData.full_name}
        onChange={handleChange}
        icon={IoPerson}
        required
      />

      <Input
        label="Email"
        name="email"
        type="email"
        placeholder="john@example.com"
        value={formData.email}
        onChange={handleChange}
        icon={IoMail}
        required
      />

      <div className="relative">
        <Input
          label="Password"
          name="password"
          type={showPassword ? "text" : "password"}
          placeholder="Min 8 characters"
          value={formData.password}
          onChange={handleChange}
          icon={IoLockClosed}
          required
        />
        <button
          type="button"
          onClick={() => setShowPassword(!showPassword)}
          className="absolute right-3.5 top-[38px] text-[var(--text-tertiary)] hover:text-[var(--text-secondary)]"
        >
          {showPassword ? <IoEyeOff size={18} /> : <IoEye size={18} />}
        </button>
      </div>

      <Input
        label="Confirm Password"
        name="confirmPassword"
        type="password"
        placeholder="••••••••"
        value={formData.confirmPassword}
        onChange={handleChange}
        icon={IoLockClosed}
        required
        error={error}
      />

      <Button type="submit" isLoading={isRegistering} className="w-full mt-2">
        Create Account
      </Button>

      <p className="text-center text-sm text-[var(--text-secondary)]">
        Already have an account?{" "}
        <Link
          to="/login"
          className="text-pastel-blue-400 hover:text-pastel-blue-500 font-medium"
        >
          Sign in
        </Link>
      </p>
    </motion.form>
  );
};

export default RegisterForm;
