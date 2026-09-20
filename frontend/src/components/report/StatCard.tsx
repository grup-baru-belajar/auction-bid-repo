import type { ReactNode } from "react";
import { ArrowDown, ArrowUp } from "lucide-react";

interface StatCardTrend {
  /** Percentage change vs. the comparison period, e.g. 12.5 or -4.2. */
  percent: number;
  /** Defaults to "from last month". */
  label?: string;
}

interface StatCardProps {
  title: string;
  value: number | string;
  unit?: string;
  prefix?: string;
  icon: ReactNode;
  trend?: StatCardTrend;
}

const StatCard = ({
  title,
  value,
  unit,
  prefix,
  icon,
  trend,
}: StatCardProps) => {
  const isPositive = (trend?.percent ?? 0) >= 0;

  return (
    <div className="flex-1 bg-white rounded-xl border border-slate-200 shadow-sm p-5 flex flex-col gap-3 min-w-0">
      <div className="flex items-center justify-between">
        <span className="text-xs text-gray-500 font-medium">{title}</span>
        <div className="w-9 h-9 rounded-xl bg-blue-50 flex items-center justify-center text-blue-500 shrink-0">
          {icon}
        </div>
      </div>

      <div className="flex items-baseline gap-1 min-w-0">
        {prefix && (
          <span className="text-lg font-semibold text-gray-500 shrink-0">
            {prefix}
          </span>
        )}
        <p className="text-2xl font-bold text-gray-900 truncate">{value}</p>
        {unit && (
          <span className="text-sm font-medium text-gray-400 shrink-0">
            {unit}
          </span>
        )}
      </div>

      {trend && (
        <div className="flex items-center gap-2">
          <span
            className={`inline-flex items-center gap-0.5 px-1.5 py-0.5 rounded-md text-xs font-semibold ${
              isPositive
                ? "bg-green-50 text-green-600"
                : "bg-red-50 text-red-600"
            }`}
          >
            {isPositive ? <ArrowUp size={12} /> : <ArrowDown size={12} />}
            {Math.abs(trend.percent)}%
          </span>
          <span className="text-xs text-gray-400">
            {trend.label ?? "from last month"}
          </span>
        </div>
      )}
    </div>
  );
};

export default StatCard;
