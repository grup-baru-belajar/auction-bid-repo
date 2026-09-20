import { Cell, Pie, PieChart, ResponsiveContainer, Tooltip } from "recharts";
import type { AuctionStatusCount, AuctionStatusName } from "../../types";
import ChartCard from "./ChartCard";

const STATUS_META: Record<AuctionStatusName, { label: string; color: string }> =
  {
    ACTIVE: { label: "Active", color: "#0ea5e9" },
    ENDED: { label: "Ended", color: "#075985" },
  };

interface AuctionStatusChartProps {
  data: AuctionStatusCount[];
}

const AuctionStatusChart = ({ data }: AuctionStatusChartProps) => {
  const total = data.reduce((sum, d) => sum + d.total, 0);

  const chartData = data.map((d) => ({
    name: STATUS_META[d.status]?.label ?? d.status,
    value: d.total,
    color: STATUS_META[d.status]?.color ?? "#94a3b8",
  }));

  return (
    <ChartCard
      title="Auction Status"
      subtitle="Distribusi status seluruh auction"
      className="lg:flex-1"
    >
      {total === 0 ? (
        <div className="h-44 flex items-center justify-center text-sm text-slate-400">
          No auctions yet.
        </div>
      ) : (
        <div className="h-full flex flex-col items-center justify-center gap-8">
          <div className="relative w-48 h-48 shrink-0">
            <ResponsiveContainer width="100%" height="100%">
              <PieChart>
                <Pie
                  data={chartData}
                  dataKey="value"
                  nameKey="name"
                  innerRadius={62}
                  outerRadius={88}
                  startAngle={90}
                  endAngle={-270}
                  stroke="none"
                >
                  {chartData.map((entry) => (
                    <Cell key={entry.name} fill={entry.color} />
                  ))}
                </Pie>
                <Tooltip
                  contentStyle={{
                    background: "#0f172a",
                    border: "none",
                    borderRadius: 6,
                    fontSize: 12,
                  }}
                  itemStyle={{ color: "#fff" }}
                  formatter={(value, name) => [value, name]}
                />
              </PieChart>
            </ResponsiveContainer>
            <div className="absolute inset-0 flex flex-col items-center justify-center pointer-events-none">
              <span className="text-[28px] font-bold text-slate-900 leading-none">
                {total}
              </span>
              <span className="text-xs text-slate-500 mt-1">
                Total Auctions
              </span>
            </div>
          </div>

          <div className="w-full flex flex-col gap-2 min-w-0">
            {chartData.map((entry) => {
              const pct =
                total > 0 ? Math.round((entry.value / total) * 100) : 0;
              return (
                <div
                  key={entry.name}
                  className="flex items-center gap-2 py-2 border-b border-slate-100 last:border-0"
                >
                  <span
                    className="w-2.5 h-2.5 rounded-full shrink-0"
                    style={{ backgroundColor: entry.color }}
                  />
                  <span className="text-sm text-slate-700 flex-1">
                    {entry.name}
                  </span>
                  <span className="text-sm font-semibold text-slate-900">
                    {entry.value}
                  </span>
                  <span className="text-xs text-slate-400 w-9 text-right">
                    {pct}%
                  </span>
                </div>
              );
            })}
          </div>
        </div>
      )}
    </ChartCard>
  );
};

export default AuctionStatusChart;
