import { Outlet } from "react-router-dom";
import { useAppSelector } from "../../store/hooks";
import Topbar from "./Topbar";
import Sidebar from "./Sidebar";

const Layout = () => {
  const { user } = useAppSelector((state) => state.auth);
  const isAdmin = user?.role === "ADMIN";

  return (
    <div className="min-h-screen bg-gray-100 flex flex-col">
      <Topbar />
      <div className="flex flex-1">
        {isAdmin && <Sidebar />}
        <main className="flex-1 p-6">
          <Outlet />
        </main>
      </div>
    </div>
  );
};

export default Layout;
