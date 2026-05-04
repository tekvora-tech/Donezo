import { Outlet } from "react-router-dom";
import { motion } from "framer-motion";
import Navbar from "./Navbar";
import Sidebar from "./Sidebar";
import ToastContainer from "../ui/Toast";

const AppLayout = () => {
  return (
    <div className="min-h-screen bg-[var(--bg-primary)]">
      <Navbar />
      <div className="flex">
        <Sidebar />
        <motion.main
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ duration: 0.3 }}
          className="flex-1 min-h-[calc(100vh-64px)] p-4 sm:p-6 lg:p-8"
        >
          <Outlet />
        </motion.main>
      </div>
      <ToastContainer />
    </div>
  );
};

export default AppLayout;
