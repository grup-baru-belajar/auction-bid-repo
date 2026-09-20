import { Sidebar, SidebarItem } from "../../components/common/Sidebar";
import { BarChart2, Gavel } from "lucide-react";
import { Outlet } from "react-router-dom";

const ReportPage = () => {
  return (
    <div className="flex h-screen">
      <Sidebar>
        <SidebarItem
          to="/report/analytics"
          icon={<BarChart2 size={20} />}
          text="Analytics"
        />
        <SidebarItem
          to="/report/auction"
          icon={<Gavel size={20} />}
          text="Auction"
        />
      </Sidebar>

      <main className="flex-1 overflow-y-auto p-6">
        <Outlet />
      </main>
    </div>
  );
};

export default ReportPage;
