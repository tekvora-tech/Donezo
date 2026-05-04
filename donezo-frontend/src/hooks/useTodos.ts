import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { todosApi } from "../api/todos";
import { useUIStore } from "../stores/uiStore";
import { type TodoFilters, type UpdateTodoInput } from "../types/todo";

export function useTodos(filters: TodoFilters = {}) {
  const { addToast } = useUIStore();
  const queryClient = useQueryClient();

  const todosQuery = useQuery({
    queryKey: ["todos", filters],
    queryFn: () => todosApi.list(filters),
    select: (data) => data.result,
  });

  const createMutation = useMutation({
    mutationFn: todosApi.create,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["todos"] });
      addToast("Todo created! ✨", "success");
    },
    onError: (error: Error) => addToast(error.message, "error"),
  });

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateTodoInput }) =>
      todosApi.update(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["todos"] });
      addToast("Todo updated! ✨", "success");
    },
    onError: (error: Error) => addToast(error.message, "error"),
  });

  const deleteMutation = useMutation({
    mutationFn: todosApi.delete,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["todos"] });
      addToast("Todo deleted", "info");
    },
    onError: (error: Error) => addToast(error.message, "error"),
  });

  return {
    todos: todosQuery.data?.data ?? [],
    metadata: todosQuery.data?.metadata,
    isLoading: todosQuery.isLoading,
    isFetching: todosQuery.isFetching,
    error: todosQuery.error,
    createTodo: createMutation.mutateAsync,
    updateTodo: updateMutation.mutateAsync,
    deleteTodo: deleteMutation.mutateAsync,
    isCreating: createMutation.isPending,
    isUpdating: updateMutation.isPending,
    isDeleting: deleteMutation.isPending,
  };
}

export function useTodoDetail(id: string) {
  const { addToast } = useUIStore();
  const queryClient = useQueryClient();

  const todoQuery = useQuery({
    queryKey: ["todo", id],
    queryFn: () => todosApi.getById(id),
    select: (data) => data.result,
    enabled: !!id,
  });

  const updateMutation = useMutation({
    mutationFn: (data: UpdateTodoInput) => todosApi.update(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["todo", id] });
      queryClient.invalidateQueries({ queryKey: ["todos"] });
      addToast("Todo updated! ✨", "success");
    },
    onError: (error: Error) => addToast(error.message, "error"),
  });

  const deleteMutation = useMutation({
    mutationFn: () => todosApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["todos"] });
      addToast("Todo deleted", "info");
    },
    onError: (error: Error) => addToast(error.message, "error"),
  });

  // Sub-tasks
  const createSubTaskMutation = useMutation({
    mutationFn: (data: { title: string }) => todosApi.createSubTask(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["todo", id] });
      addToast("Sub-task added", "success");
    },
    onError: (error: Error) => addToast(error.message, "error"),
  });

  const updateSubTaskMutation = useMutation({
    mutationFn: ({
      subTaskId,
      data,
    }: {
      subTaskId: string;
      data: { title?: string; is_completed?: boolean };
    }) => todosApi.updateSubTask(id, subTaskId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["todo", id] });
    },
    onError: (error: Error) => addToast(error.message, "error"),
  });

  const deleteSubTaskMutation = useMutation({
    mutationFn: (subTaskId: string) => todosApi.deleteSubTask(id, subTaskId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["todo", id] });
      addToast("Sub-task removed", "info");
    },
    onError: (error: Error) => addToast(error.message, "error"),
  });

  return {
    todo: todoQuery.data,
    isLoading: todoQuery.isLoading,
    updateTodo: updateMutation.mutateAsync,
    deleteTodo: deleteMutation.mutateAsync,
    isUpdating: updateMutation.isPending,
    isDeleting: deleteMutation.isPending,
    createSubTask: createSubTaskMutation.mutateAsync,
    updateSubTask: updateSubTaskMutation.mutateAsync,
    deleteSubTask: deleteSubTaskMutation.mutateAsync,
    isCreatingSubTask: createSubTaskMutation.isPending,
  };
}
