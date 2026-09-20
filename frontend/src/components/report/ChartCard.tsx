import type { ReactNode } from "react";

interface ChartCardProps {
  title: string;
  subtitle?: string;
  headerRight?: ReactNode;
  summary?: ReactNode;
  children: ReactNode;
  className?: string;
}

const ChartCard = ({
  title,
  subtitle,
  headerRight,
  summary,
  children,
  className = "",
}: ChartCardProps) => {
  return (
    <div
      className={`flex-1 bg-white rounded-xl border border-slate-200 shadow-sm p-6 flex flex-col gap-5 min-w-0 ${className}`}
    >
      <div className="flex flex-wrap items-start justify-between gap-4">
        <div>
          <h3 className="text-base font-semibold text-slate-900">{title}</h3>
          {subtitle && (
            <p className="text-[13px] text-slate-500 mt-0.5">{subtitle}</p>
          )}
        </div>
        {headerRight}
      </div>

      {summary}

      <div className="min-w-0">{children}</div>
    </div>
  );
};

export default ChartCard;
