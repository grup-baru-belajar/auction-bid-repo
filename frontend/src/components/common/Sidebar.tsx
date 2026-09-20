import type { ReactNode } from "react";
import { Link, useLocation } from "react-router-dom";

interface SidebarProps {
  children: ReactNode;
}

const Sidebar = ({ children }: SidebarProps) => {
  return (
    <aside className="h-screen">
      <nav className="h-full flex flex-col bg-sky-900 border-r shadow-sm">
        <div className="p-4 pb-2 flex items-center">
          <img
            src="https://img.logoipsum.com/243.svg"
            className="w-32"
            alt=""
          />
        </div>

        {/* AuctionApp badge */}
        <div className="flex items-center gap-2 px-3 py-2 mx-3 mb-4 pb-4 border-b border-sky-800">
          <span className="flex h-7 w-7 shrink-0 items-center justify-center rounded bg-[#2051E5] text-md font-bold text-white">
            a.
          </span>
          <span className="whitespace-nowrap text-md font-medium text-slate-300">
            Auction<span className="font-extralight">App</span>
          </span>
        </div>

        <ul className="flex-1 px-3">{children}</ul>
      </nav>
    </aside>
  );
};

interface SidebarItemProps {
  to: string;
  icon: ReactNode;
  text: string;
}

const SidebarItem = ({ to, icon, text }: SidebarItemProps) => {
  const location = useLocation();
  const active = location.pathname === to;

  return (
    <li>
      <Link
        to={to}
        className={`
          relative flex items-center py-2 px-3 my-1
          font-medium rounded-md cursor-pointer
          transition-colors group
          ${
            active
              ? "bg-slate-100 text-gray-900"
              : "hover:bg-indigo-50 text-gray-400"
          }
      `}
      >
        {icon}
        <span className="ml-3">{text}</span>
      </Link>
    </li>
  );
};

export { Sidebar, SidebarItem };
