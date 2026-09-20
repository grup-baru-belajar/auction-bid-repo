interface NavbarProps {
  webName: string;
  userName: string;
}

const Navbar = ({ webName, userName }: NavbarProps) => {
  return (
    <nav className="h-16 bg-sky-500 border-b shadow-sm flex items-center justify-between px-6">
      <h1 className="text-lg font-semibold text-gray-200">{webName}</h1>

      <div className="flex items-center gap-3">
        <img
          src="https://ui-avatars.com/api/?background=c7d2fe&color=3730a3&bold=true"
          alt=""
          className="w-9 h-9 rounded-md"
        />
        <span className="text-sm font-medium text-gray-200">
          Hello, {userName}
        </span>
      </div>
    </nav>
  );
};

export default Navbar;