import { useThemeStore } from "../stores/themeStore";

export function useTheme() {
  const { isDark, toggleTheme, setTheme } = useThemeStore();

  return {
    isDark,
    toggleTheme,
    setTheme,
    theme: isDark ? "dark" : "light",
  };
}
