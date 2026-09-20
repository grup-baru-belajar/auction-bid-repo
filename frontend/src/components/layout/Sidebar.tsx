import { NavLink } from "react-router-dom";

const Sidebar = () => {
  return (
    <aside className="w-56 bg-white border-r border-gray-200 shrink-0">
      <nav className="p-4 space-y-1">
        <NavLink
          to="/"
          className={({ isActive }) =>
            `flex items-center gap-2.5 px-3 py-2.5 rounded-lg text-sm font-medium transition-colors ${
              isActive
                ? "bg-[#1A4B69] text-white"
                : "text-gray-600 hover:bg-gray-100"
            }`
          }
        >
          <svg className="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
            <path d="M6 2L3 6v14a2 2 0 002 2h14a2 2 0 002-2V6l-3-4z" />
            <line x1="3" y1="6" x2="21" y2="6" />
            <path d="M16 10a4 4 0 01-8 0" />
          </svg>
          Auction
        </NavLink>
        <NavLink
          to="/report"
          className={({ isActive }) =>
            `flex items-center gap-2.5 px-3 py-2.5 rounded-lg text-sm font-medium transition-colors ${
              isActive
                ? "bg-[#1A4B69] text-white"
                : "text-gray-600 hover:bg-gray-100"
            }`
          }
        >
          <svg className="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
            <path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z" />
            <polyline points="14 2 14 8 20 8" />
            <line x1="16" y1="13" x2="8" y2="13" />
            <line x1="16" y1="17" x2="8" y2="17" />
          </svg>
          Report
        </NavLink>
      </nav>
    </aside>
  );
};

export default Sidebar;
