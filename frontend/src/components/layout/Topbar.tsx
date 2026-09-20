import { useState, useRef, useEffect } from "react";
import { useAppDispatch, useAppSelector } from "../../store/hooks";
import { logout } from "../../features/auth/authSlice";
import { useNavigate, useSearchParams } from "react-router-dom";

const Topbar = () => {
  const dispatch = useAppDispatch();
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const { user, isAuthenticated } = useAppSelector((state) => state.auth);
  const isAdmin = user?.role === "ADMIN";

  const [dropdownOpen, setDropdownOpen] = useState(false);
  const dropdownRef = useRef<HTMLDivElement>(null);
  const [searchTerm, setSearchTerm] = useState(
    searchParams.get("search") ?? "",
  );

  useEffect(() => {
    const handleClickOutside = (e: MouseEvent) => {
      if (
        dropdownRef.current &&
        !dropdownRef.current.contains(e.target as Node)
      ) {
        setDropdownOpen(false);
      }
    };
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  const handleLogout = () => {
    dispatch(logout());
    navigate("/login");
  };

  const handleSearchSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    const trimmed = searchTerm.trim();
    navigate(trimmed ? `/?search=${encodeURIComponent(trimmed)}` : "/");
  };

  return (
    <header className="bg-sky-500 sticky top-0 z-30">
      <div className="flex items-center justify-between gap-8 px-6 py-3">
        <div className="flex items-center gap-3 shrink-0">
          {isAdmin && (
            <span className="text-xs font-semibold bg-sky-800 text-white px-2 py-0.5 rounded">
              ADMIN
            </span>
          )}
          <h1 className="text-xl font-semibold text-sky-50">Auction App</h1>
        </div>

        <form
          onSubmit={handleSearchSubmit}
          className="flex-1 max-w-lg flex items-center"
        >
          <input
            type="text"
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            placeholder="Search auction..."
            className="flex-1 h-10 rounded-l-md border-0 px-3 text-sm bg-white placeholder:text-gray-400 focus:outline-none focus:ring-2 focus:ring-sky-700"
          />
          <button
            type="submit"
            aria-label="Search"
            className="h-10 w-10 flex items-center justify-center rounded-r-md bg-sky-600 hover:bg-sky-700 text-white transition-colors"
          >
            <svg
              className="w-4 h-4"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              strokeWidth="2"
            >
              <circle cx="11" cy="11" r="7" />
              <line x1="21" y1="21" x2="16.65" y2="16.65" />
            </svg>
          </button>
        </form>

        <div className="shrink-0">
          {!isAuthenticated ? (
            <button
              onClick={() => navigate("/login")}
              className="px-4 py-2 bg-sky-800 hover:bg-sky-900 text-white text-sm font-semibold rounded-lg transition-colors"
            >
              Login
            </button>
          ) : (
            <div className="relative" ref={dropdownRef}>
              <button
                onClick={() => setDropdownOpen(!dropdownOpen)}
                className="flex items-center gap-2 text-sm text-sky-50 hover:text-white font-medium transition-colors"
              >
                <div className="w-8 h-8 rounded-full bg-sky-800 text-white flex items-center justify-center text-xs font-bold">
                  {user?.name?.charAt(0).toUpperCase()}
                </div>
                <svg
                  className={`w-4 h-4 transition-transform ${dropdownOpen ? "rotate-180" : ""}`}
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="2"
                >
                  <polyline points="6 9 12 15 18 9" />
                </svg>
              </button>

              {dropdownOpen && (
                <div className="absolute right-0 top-full mt-1 w-48 bg-white border border-gray-200 rounded-lg shadow-lg py-1 z-50">
                  <div className="px-4 py-2 border-b border-gray-100">
                    <p className="text-sm font-medium text-gray-800">
                      {user?.name}
                    </p>
                    <p className="text-xs text-gray-500">{user?.username}</p>
                  </div>
                  <button
                    onClick={handleLogout}
                    className="w-full flex items-center gap-2 px-4 py-2 text-sm text-red-600 hover:bg-red-50 transition-colors"
                  >
                    <svg
                      className="w-4 h-4"
                      viewBox="0 0 24 24"
                      fill="none"
                      stroke="currentColor"
                      strokeWidth="2"
                    >
                      <path d="M9 21H5a2 2 0 01-2-2V5a2 2 0 012-2h4" />
                      <polyline points="16 17 21 12 16 7" />
                      <line x1="21" y1="12" x2="9" y2="12" />
                    </svg>
                    Keluar
                  </button>
                </div>
              )}
            </div>
          )}
        </div>
      </div>
    </header>
  );
};

export default Topbar;
