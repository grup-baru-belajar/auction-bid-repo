import { useEffect } from "react";
import { Gavel, TrendingUp, Users, Wallet } from "lucide-react";
import { useAppDispatch, useAppSelector } from "../../store/hooks";
import {
  fetchAuctionActivity,
  fetchAuctionStatus,
  fetchAuctionSummary,
  fetchTopAuctions,
  fetchTotalBidders,
  fetchTransactionOverview,
  setActivityIntervalDays,
} from "../../features/reports/reportsSlice";
import StatCard from "../../components/report/StatCard";
import AuctionActivityChart from "../../components/report/AuctionActivityChart";
import AuctionStatusChart from "../../components/report/AuctionStatusChart";
import BiddingActivityChart from "../../components/report/BiddingActivityChart";
import TransactionOverviewChart from "../../components/report/TransactionOverviewChart";
import TopAuctionsTable from "../../components/report/TopAuctionsTable";

const TRANSACTION_WEEKS = 4;

const ReportAnalyticsPage = () => {
  const dispatch = useAppDispatch();
  const {
    summary,
    loading,
    error,
    totalBidders,
    activity,
    activityIntervalDays,
    activityLoading,
    statusBreakdown,
    transactionOverview,
    topAuctions,
  } = useAppSelector((state) => state.reports);

  useEffect(() => {
    dispatch(fetchAuctionSummary());
    dispatch(fetchTotalBidders());
    dispatch(fetchAuctionStatus());
    dispatch(fetchTransactionOverview(TRANSACTION_WEEKS));
    dispatch(fetchTopAuctions());
  }, [dispatch]);

  useEffect(() => {
    dispatch(fetchAuctionActivity(activityIntervalDays));
  }, [dispatch, activityIntervalDays]);

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

      {!loading && !error && summary && (
        <div className="flex flex-col gap-4 mt-4">
          <div className="flex flex-col lg:flex-row gap-4">
            <AuctionActivityChart
              data={activity ?? []}
              intervalDays={activityIntervalDays}
              onIntervalChange={(days) =>
                dispatch(setActivityIntervalDays(days))
              }
              loading={activityLoading}
            />
            <AuctionStatusChart data={statusBreakdown ?? []} />
          </div>

          <div className="flex flex-col lg:flex-row gap-4">
            <BiddingActivityChart
              data={activity ?? []}
              intervalDays={activityIntervalDays}
              loading={activityLoading}
            />
            <TransactionOverviewChart data={transactionOverview ?? []} />
          </div>

          <TopAuctionsTable data={topAuctions ?? []} />
        </div>
      )}
    </div>
  );
};

export default ReportAnalyticsPage;
