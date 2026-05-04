import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { authApi } from "../api/auth";
import { usersApi } from "../api/users";
import { useAuthStore } from "../stores/authStore";
import { useUIStore } from "../stores/uiStore";

export function useAuth() {
  const { setAuth, logout: storeLogout, user } = useAuthStore();
  const { addToast } = useUIStore();
  const queryClient = useQueryClient();

  const loginMutation = useMutation({
    mutationFn: authApi.login,
    onSuccess: (data) => {
      setAuth(data.result.user, {
        access_token: data.result.access_token,
        refresh_token: data.result.refresh_token,
        expires_in: data.result.expires_in,
      });
      addToast("Welcome back! 👋", "success");
    },
    onError: (error: Error) => {
      addToast(error.message, "error");
    },
  });

  const registerMutation = useMutation({
    mutationFn: authApi.register,
    onSuccess: (data) => {
      setAuth(data.result.user, {
        access_token: data.result.access_token,
        refresh_token: data.result.refresh_token,
        expires_in: data.result.expires_in,
      });
      addToast("Account created successfully! 🎉", "success");
    },
    onError: (error: Error) => {
      addToast(error.message, "error");
    },
  });

  const logoutMutation = useMutation({
    mutationFn: authApi.logout,
    onSuccess: () => {
      storeLogout();
      queryClient.clear();
      addToast("Logged out successfully", "info");
    },
    onError: () => {
      storeLogout();
      queryClient.clear();
    },
  });

  const profileQuery = useQuery({
    queryKey: ["profile"],
    queryFn: usersApi.getProfile,
    enabled: !!user,
    select: (data) => data.result,
  });

  const updateProfileMutation = useMutation({
    mutationFn: usersApi.updateProfile,
    onSuccess: (data) => {
      useAuthStore.getState().updateUser(data.result);
      queryClient.invalidateQueries({ queryKey: ["profile"] });
      addToast("Profile updated!", "success");
    },
    onError: (error: Error) => {
      addToast(error.message, "error");
    },
  });

  return {
    user,
    isAuthenticated: !!user,
    login: loginMutation.mutateAsync,
    register: registerMutation.mutateAsync,
    logout: logoutMutation.mutateAsync,
    isLoggingIn: loginMutation.isPending,
    isRegistering: registerMutation.isPending,
    isLoggingOut: logoutMutation.isPending,
    profile: profileQuery.data,
    isLoadingProfile: profileQuery.isLoading,
    updateProfile: updateProfileMutation.mutateAsync,
    isUpdatingProfile: updateProfileMutation.isPending,
  };
}
