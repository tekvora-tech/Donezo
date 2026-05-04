import { useState } from "react";
import { motion } from "framer-motion";
import { IoPerson, IoMail, IoCamera, IoSave } from "react-icons/io5";
import Card from "../components/ui/Card";
import Input from "../components/ui/Input";
import Button from "../components/ui/Button";
import { useAuth } from "../hooks/useAuth";
import { cn } from "../utils/cn";

const ProfilePage = () => {
  const { user, profile, isLoadingProfile, updateProfile, isUpdatingProfile } =
    useAuth();
  const [formData, setFormData] = useState({
    full_name: profile?.full_name || "",
    avatar_url: profile?.avatar_url || "",
  });

  const handleSubmit = async (e) => {
    e.preventDefault();
    await updateProfile(formData);
  };

  if (isLoadingProfile) {
    return (
      <div className="max-w-2xl mx-auto space-y-4">
        {[1, 2].map((i) => (
          <div
            key={i}
            className="h-32 bg-[var(--bg-tertiary)] rounded-2xl animate-pulse"
          />
        ))}
      </div>
    );
  }

  return (
    <div className="max-w-2xl mx-auto space-y-6">
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
      >
        <h1 className="text-2xl font-bold text-[var(--text-primary)] mb-6">
          Profile Settings
        </h1>

        {/* Avatar Section */}
        <Card className="p-6 mb-6">
          <div className="flex items-center gap-6">
            <div className="relative">
              <div className="w-24 h-24 rounded-2xl bg-gradient-to-br from-pastel-blue-400 to-pastel-lavender-400 flex items-center justify-center text-white text-3xl font-bold shadow-lg">
                {profile?.full_name?.charAt(0).toUpperCase() || "U"}
              </div>
              <button className="absolute -bottom-2 -right-2 p-2 rounded-xl bg-[var(--bg-primary)] border border-[var(--border-color)] shadow-md text-[var(--text-secondary)] hover:text-pastel-blue-400 transition-colors">
                <IoCamera size={16} />
              </button>
            </div>
            <div>
              <h2 className="text-xl font-semibold text-[var(--text-primary)]">
                {profile?.full_name}
              </h2>
              <p className="text-[var(--text-secondary)]">{profile?.email}</p>
              <p className="text-xs text-[var(--text-tertiary)] mt-1">
                Member since{" "}
                {new Date(profile?.created_at).toLocaleDateString()}
              </p>
            </div>
          </div>
        </Card>

        {/* Edit Form */}
        <Card className="p-6">
          <form onSubmit={handleSubmit} className="space-y-5">
            <Input
              label="Full Name"
              value={formData.full_name}
              onChange={(e) =>
                setFormData((prev) => ({ ...prev, full_name: e.target.value }))
              }
              icon={IoPerson}
            />

            <Input
              label="Email"
              value={profile?.email || ""}
              icon={IoMail}
              disabled
              className="opacity-60"
            />

            <div className="pt-2">
              <Button
                type="submit"
                isLoading={isUpdatingProfile}
                className="w-full sm:w-auto"
              >
                <IoSave size={18} /> Save Changes
              </Button>
            </div>
          </form>
        </Card>
      </motion.div>
    </div>
  );
};

export default ProfilePage;
