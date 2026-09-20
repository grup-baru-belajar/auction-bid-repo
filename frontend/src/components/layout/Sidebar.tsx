import { NavLink } from "react-router-dom";
import { useAppSelector } from "../../store/hooks";

interface SidebarProps {
  /** Whether the off-canvas drawer is open. Ignored at `lg` and up, where the sidebar is always visible. */
  open: boolean;
  onClose: () => void;
}

const Sidebar = ({ open, onClose }: SidebarProps) => {
  const { user } = useAppSelector((state) => state.auth);
  const isAdmin = user?.role === "ADMIN";

  // Figma: admin sidebar is sky/900 (dark), regular users get a lighter sky/500
  const bgClass = isAdmin ? "bg-sky-900" : "bg-sky-500";
  const activeClass = "bg-slate-50 text-slate-900";
  const inactiveClass = isAdmin
    ? "text-sky-200/70 hover:bg-sky-800 hover:text-white"
    : "text-white/80 hover:bg-sky-600 hover:text-white";

  const linkClass = ({ isActive }: { isActive: boolean }) =>
    `flex items-center gap-2.5 px-3 py-2.5 rounded-lg text-sm font-medium transition-colors ${
      isActive ? activeClass : inactiveClass
    }`;

  return (
    <>
      {open && (
        <div
          className="fixed inset-0 z-30 bg-black/40 lg:hidden"
          onClick={onClose}
          aria-hidden="true"
        />
      )}

      <aside
        className={`fixed inset-y-0 left-0 z-40 w-56 ${bgClass} shrink-0 transform transition-transform duration-200 ease-in-out lg:static lg:translate-x-0 ${
          open ? "translate-x-0" : "-translate-x-full"
        }`}
      >
        <nav className="p-4 space-y-1">
          <div className="flex items-center gap-2 px-3 py-2 mb-4 pb-4 border-b border-white/10">
            <span className="flex h-7 w-7 shrink-0 items-center justify-center rounded bg-[#2051E5] text-sm font-bold text-white">
              a.
            </span>
            <span className="whitespace-nowrap text-sm font-medium text-slate-100">
              Auction<span className="font-light">App</span>
            </span>
          </div>

          <NavLink
            to={isAdmin ? "/report/auction" : "/"}
            className={linkClass}
            onClick={onClose}
          >
            <svg
              className="w-4 h-4"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              strokeWidth="2"
            >
              <path d="M6 2L3 6v14a2 2 0 002 2h14a2 2 0 002-2V6l-3-4z" />
              <line x1="3" y1="6" x2="21" y2="6" />
              <path d="M16 10a4 4 0 01-8 0" />
            </svg>
            Auction
          </NavLink>

          {isAdmin && (
            <NavLink
              to="/report/analytics"
              className={linkClass}
              onClick={onClose}
            >
              <svg
                className="w-4 h-4"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                strokeWidth="2"
              >
                <path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z" />
                <polyline points="14 2 14 8 20 8" />
                <line x1="16" y1="13" x2="8" y2="13" />
                <line x1="16" y1="17" x2="8" y2="17" />
              </svg>
              Report
            </NavLink>
          )}
        </nav>
      </aside>
    </>
  );
};

export default Sidebar;
