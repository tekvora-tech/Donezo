import { useParams } from "react-router-dom";
import TodoDetail from "../components/todos/TodoDetail";

const TodoDetailPage = () => {
  const { id } = useParams();

  return <TodoDetail todoId={id} />;
};

export default TodoDetailPage;
