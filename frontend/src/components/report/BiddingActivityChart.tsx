import {
  Area,
  AreaChart,
  CartesianGrid,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";
import type { AuctionActivity } from "../../types";
import ChartCard from "./ChartCard";
import { formatShortDate } from "../common/formatters";

interface BiddingActivityChartProps {
  data: AuctionActivity[];
  intervalDays: number;
  loading?: boolean;
}

const BiddingActivityChart = ({
  data,
  intervalDays,
  loading,
}: BiddingActivityChartProps) => {
  const total = data.reduce((sum, d) => sum + d.totalBids, 0);
  const peak = data.reduce(
    (max, d) => (d.totalBids > max.totalBids ? d : max),
    { date: "", totalBids: 0, totalAuctions: 0 } as AuctionActivity,
  );

  const chartData = data.map((d) => ({
    ...d,
    label: formatShortDate(d.date),
  }));

  return (
    <ChartCard
      title="Bidding Activity"
      subtitle={`Jumlah bid masuk per hari · ${intervalDays} hari terakhir`}
      summary={
        <div className="flex items-baseline gap-2">
          <span className="text-2xl font-bold text-slate-900">
            {total.toLocaleString("id-ID")}
          </span>
          <span className="text-[13px] text-slate-500">
            {peak.totalBids > 0
              ? `bid · puncak ${peak.totalBids} bid pada ${formatShortDate(peak.date)}`
              : "bid"}
          </span>
        </div>
      }
    >
      {loading ? (
        <div className="h-52 flex items-center justify-center text-sm text-slate-400">
          Loading...
        </div>
      ) : chartData.length === 0 ? (
        <div className="h-52 flex items-center justify-center text-sm text-slate-400">
          No bids in this period.
        </div>
      ) : (
        <ResponsiveContainer width="100%" height={208}>
          <AreaChart data={chartData} margin={{ left: -20, top: 8 }}>
            <defs>
              <linearGradient id="biddingFill" x1="0" y1="0" x2="0" y2="1">
                <stop offset="0%" stopColor="#0ea5e9" stopOpacity={0.15} />
                <stop offset="100%" stopColor="#0ea5e9" stopOpacity={0} />
              </linearGradient>
            </defs>
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
              contentStyle={{
                background: "#0f172a",
                border: "none",
                borderRadius: 6,
                fontSize: 12,
              }}
              labelStyle={{ color: "#cbd5e1" }}
              itemStyle={{ color: "#fff" }}
              formatter={(value) => [`${value} bids`, ""]}
            />
            <Area
              type="monotone"
              dataKey="totalBids"
              stroke="#0284c7"
              strokeWidth={2}
              fill="url(#biddingFill)"
              dot={{ r: 3, fill: "#fff", stroke: "#0284c7", strokeWidth: 2 }}
            />
          </AreaChart>
        </ResponsiveContainer>
      )}
    </ChartCard>
  );
};

export default BiddingActivityChart;
