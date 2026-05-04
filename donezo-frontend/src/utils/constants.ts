export const API_BASE_URL = import.meta.env.VITE_API_URL || "/api";

export const PRIORITY_COLORS: Record<
  string,
  { bg: string; text: string; border: string }
> = {
  low: {
    bg: "bg-pastel-mint-100",
    text: "text-pastel-mint-600",
    border: "border-pastel-mint-300",
  },
  medium: {
    bg: "bg-pastel-blue-100",
    text: "text-pastel-blue-600",
    border: "border-pastel-blue-300",
  },
  high: {
    bg: "bg-pastel-peach-100",
    text: "text-pastel-peach-600",
    border: "border-pastel-peach-300",
  },
  urgent: {
    bg: "bg-pastel-rose-100",
    text: "text-pastel-rose-600",
    border: "border-pastel-rose-300",
  },
};

export const STATUS_COLORS: Record<string, { bg: string; text: string }> = {
  pending: { bg: "bg-pastel-blue-100", text: "text-pastel-blue-600" },
  in_progress: { bg: "bg-pastel-lemon-100", text: "text-pastel-lemon-600" },
  completed: { bg: "bg-pastel-mint-100", text: "text-pastel-mint-600" },
  cancelled: { bg: "bg-slate-100", text: "text-slate-500" },
};

export const PRIORITY_OPTIONS = [
  { value: "low", label: "Low", icon: "↓" },
  { value: "medium", label: "Medium", icon: "→" },
  { value: "high", label: "High", icon: "↑" },
  { value: "urgent", label: "Urgent", icon: "!" },
] as const;

export const STATUS_OPTIONS = [
  { value: "pending", label: "Pending" },
  { value: "in_progress", label: "In Progress" },
  { value: "completed", label: "Completed" },
  { value: "cancelled", label: "Cancelled" },
] as const;

export const DEFAULT_TAG_COLORS = [
  "#60A5FA", // pastel blue
  "#4ADE80", // pastel mint
  "#C084FC", // pastel lavender
  "#FB923C", // pastel peach
  "#FB7185", // pastel rose
  "#FACC15", // pastel lemon
  "#38BDF8", // sky
  "#A78BFA", // violet
];
