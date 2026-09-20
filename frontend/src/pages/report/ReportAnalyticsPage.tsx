import { useEffect } from "react";
import { Gavel, TrendingUp, Users, Wallet } from "lucide-react";
import { useAppDispatch, useAppSelector } from "../../store/hooks";
import {
  fetchAuctionSummary,
  fetchTotalBidders,
} from "../../features/reports/reportsSlice";
import StatCard from "../../components/report/StatCard";

const ReportAnalyticsPage = () => {
  const dispatch = useAppDispatch();
  const { summary, loading, error, totalBidders } = useAppSelector(
    (state) => state.reports,
  );

  useEffect(() => {
    dispatch(fetchAuctionSummary());
    dispatch(fetchTotalBidders());
  }, [dispatch]);

  return (
    <div>
      <h1 className="text-2xl font-semibold mb-4">Analytics</h1>

      {loading && (
        <p className="text-center text-gray-500 py-10">Loading Data...</p>
      )}
      {error && <p className="text-center text-red-500 py-10">{error}</p>}

      {!loading && !error && summary && (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
          <StatCard
            title="Total Auctions"
            value={summary.totalAuctions}
            unit="Auctions"
            icon={<Gavel size={18} />}
          />
          <StatCard
            title="Total Bids"
            value={summary.totalBidsAll}
            unit="Bids"
            icon={<TrendingUp size={18} />}
          />
          <StatCard
            title="Total Transactions"
            value={summary.completedAuctions}
            unit="Transactions"
            icon={<Wallet size={18} />}
          />
          <StatCard
            title="Total Bidders"
            value={totalBidders ?? "..."}
            unit="Bidders"
            icon={<Users size={18} />}
          />
        </div>
      )}
    </div>
  );
};

export default ReportAnalyticsPage;
