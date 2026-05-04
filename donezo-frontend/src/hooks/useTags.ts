import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { tagsApi } from "../api/tags";
import { useUIStore } from "../stores/uiStore";
import { type UpdateTagInput } from "../types/tag";

export function useTags() {
  const { addToast } = useUIStore();
  const queryClient = useQueryClient();

  const tagsQuery = useQuery({
    queryKey: ["tags"],
    queryFn: () => tagsApi.list(),
    select: (data) => data.result.data,
  });

  const createMutation = useMutation({
    mutationFn: tagsApi.create,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["tags"] });
      addToast("Tag created! 🏷️", "success");
    },
    onError: (error: Error) => addToast(error.message, "error"),
  });

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateTagInput }) =>
      tagsApi.update(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["tags"] });
      queryClient.invalidateQueries({ queryKey: ["todos"] });
      addToast("Tag updated!", "success");
    },
    onError: (error: Error) => addToast(error.message, "error"),
  });

  const deleteMutation = useMutation({
    mutationFn: tagsApi.delete,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["tags"] });
      queryClient.invalidateQueries({ queryKey: ["todos"] });
      addToast("Tag deleted", "info");
    },
    onError: (error: Error) => addToast(error.message, "error"),
  });

  return {
    tags: tagsQuery.data ?? [],
    isLoading: tagsQuery.isLoading,
    createTag: createMutation.mutateAsync,
    updateTag: updateMutation.mutateAsync,
    deleteTag: deleteMutation.mutateAsync,
    isCreating: createMutation.isPending,
    isUpdating: updateMutation.isPending,
    isDeleting: deleteMutation.isPending,
  };
}
