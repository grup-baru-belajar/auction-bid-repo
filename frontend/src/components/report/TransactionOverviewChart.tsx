import {
  Bar,
  BarChart,
  CartesianGrid,
  Cell,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";
import type { TransactionWeek } from "../../types";
import ChartCard from "./ChartCard";
import {
  formatRupiah,
  formatRupiahCompact,
  formatWeekRange,
} from "../common/formatters";

interface TransactionOverviewChartProps {
  data: TransactionWeek[];
  loading?: boolean;
}

const TransactionOverviewChart = ({
  data,
  loading,
}: TransactionOverviewChartProps) => {
  const total = data.reduce((sum, d) => sum + d.total, 0);

  const chartData = data.map((d, i) => ({
    ...d,
    label: formatWeekRange(d.weekStart),
    isLatest: i === data.length - 1,
  }));

  return (
    <ChartCard
      title="Transaction Overview"
      subtitle="Total nilai transaksi per minggu (auction yang telah selesai)"
      summary={
        <div className="flex items-baseline gap-2">
          <span className="text-2xl font-bold text-slate-900">
            {formatRupiah(total)}
          </span>
          <span className="text-[13px] text-slate-500">total transaksi</span>
        </div>
      }
    >
      {loading ? (
        <div className="h-52 flex items-center justify-center text-sm text-slate-400">
          Loading...
        </div>
      ) : chartData.length === 0 ? (
        <div className="h-52 flex items-center justify-center text-sm text-slate-400">
          No completed transactions yet.
        </div>
      ) : (
        <ResponsiveContainer width="100%" height={208}>
          <BarChart data={chartData} margin={{ left: -10, top: 20 }}>
            <CartesianGrid vertical={false} stroke="#e2e8f0" />
            <XAxis
              dataKey="label"
              tick={{ fontSize: 11, fill: "#94a3b8" }}
              axisLine={{ stroke: "#e2e8f0" }}
              tickLine={false}
            />
            <YAxis
              tickFormatter={(v: number) => formatRupiahCompact(v)}
              tick={{ fontSize: 11, fill: "#94a3b8" }}
              axisLine={false}
              tickLine={false}
              width={70}
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
              formatter={(value) => [formatRupiah(Number(value)), ""]}
            />
            <Bar dataKey="total" radius={[4, 4, 0, 0]}>
              {chartData.map((entry) => (
                <Cell
                  key={entry.weekStart}
                  fill={entry.isLatest ? "#075985" : "#0ea5e9"}
                />
              ))}
            </Bar>
          </BarChart>
        </ResponsiveContainer>
      )}
    </ChartCard>
  );
};

export default TransactionOverviewChart;
