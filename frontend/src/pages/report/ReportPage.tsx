import { Outlet } from "react-router-dom";

const ReportPage = () => {
  return (
    <div className="flex h-screen">
      <main className="flex-1 overflow-y-auto p-6">
        <Outlet />
      </main>
    </div>
  );
};

export default ReportPage;
