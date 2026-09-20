import {
  Bar,
  BarChart,
  CartesianGrid,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";
import type { AuctionActivity } from "../../types";
import ChartCard from "./ChartCard";
import { formatShortDate } from "../common/formatters";

const PERIODS = [
  { label: "7 Days", days: 7 },
  { label: "30 Days", days: 30 },
  { label: "3 Months", days: 90 },
  { label: "1 Year", days: 365 },
];

interface AuctionActivityChartProps {
  data: AuctionActivity[];
  intervalDays: number;
  onIntervalChange: (days: number) => void;
  loading?: boolean;
}

const AuctionActivityChart = ({
  data,
  intervalDays,
  onIntervalChange,
  loading,
}: AuctionActivityChartProps) => {
  const total = data.reduce((sum, d) => sum + d.totalAuctions, 0);
  const avgPerDay = data.length > 0 ? total / data.length : 0;

  const chartData = data.map((d) => ({
    ...d,
    label: formatShortDate(d.date),
  }));

  return (
    <ChartCard
      title="Auction Activity"
      subtitle="Jumlah auction baru yang dibuat per hari"
      className="lg:flex-2"
      headerRight={
        <div className="flex items-center gap-0.5 bg-slate-100 rounded-md p-1">
          {PERIODS.map((period) => (
            <button
              key={period.days}
              onClick={() => onIntervalChange(period.days)}
              className={`px-3 py-1 rounded text-[13px] font-medium transition-colors ${
                intervalDays === period.days
                  ? "bg-white text-slate-900 shadow-sm"
                  : "text-slate-500 hover:text-slate-700"
              }`}
            >
              {period.label}
            </button>
          ))}
        </div>
      }
      summary={
        <div className="flex items-baseline gap-2">
          <span className="text-2xl font-bold text-slate-900">{total}</span>
          <span className="text-[13px] text-slate-500">
            auction · rata-rata {avgPerDay.toFixed(1)} per hari
          </span>
        </div>
      }
    >
      {loading ? (
        <div className="h-70 flex items-center justify-center text-sm text-slate-400">
          Loading...
        </div>
      ) : chartData.length === 0 ? (
        <div className="h-70 flex items-center justify-center text-sm text-slate-400">
          No data for this period.
        </div>
      ) : (
        <ResponsiveContainer width="100%" height={280}>
          <BarChart data={chartData} margin={{ left: -20, top: 8 }}>
            <CartesianGrid vertical={false} stroke="#e2e8f0" />
            <XAxis
              dataKey="label"
              tick={{ fontSize: 11, fill: "#94a3b8" }}
              axisLine={{ stroke: "#e2e8f0" }}
              tickLine={false}
            />
            <YAxis
              allowDecimals={false}
              tick={{ fontSize: 11, fill: "#94a3b8" }}
              axisLine={false}
              tickLine={false}
            />
            <Tooltip
              cursor={{ fill: "#f1f5f9" }}
              contentStyle={{
                background: "#0f172a",
                border: "none",
                borderRadius: 6,
                fontSize: 12,
              }}
              labelStyle={{ color: "#cbd5e1" }}
              itemStyle={{ color: "#fff" }}
              formatter={(value) => [`${value} auctions`, ""]}
            />
            <Bar dataKey="totalAuctions" fill="#7dd3fc" radius={[3, 3, 0, 0]} />
          </BarChart>
        </ResponsiveContainer>
      )}
    </ChartCard>
  );
};

export default AuctionActivityChart;
