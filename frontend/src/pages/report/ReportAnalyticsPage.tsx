import { useEffect } from "react";
import { useAppDispatch, useAppSelector } from "../../store/hooks";
import { fetchAuctionSummary } from "../../features/reports/reportsSlice";

const ReportAnalyticsPage = () => {
  const dispatch = useAppDispatch();
  const { summary, loading, error } = useAppSelector((state) => state.reports);

  useEffect(() => {
    dispatch(fetchAuctionSummary());
  }, [dispatch]);

  return (
    <div>
      <h1 className="text-2xl font-semibold mb-4">Analytics</h1>

      {loading && (
        <p className="text-center text-gray-500 py-10">Loading Data...</p>
      )}
      {error && <p className="text-center text-red-500 py-10">{error}</p>}

      {!loading && !error && summary && (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          <div className="bg-white border border-gray-200 rounded-lg p-4">
            <p className="text-sm text-gray-500">Total Auctions</p>
            <p className="text-2xl font-bold text-[#1A4B69]">
              {summary.totalAuctions}
            </p>
          </div>
          <div className="bg-white border border-gray-200 rounded-lg p-4">
            <p className="text-sm text-gray-500">Ongoing Auctions</p>
            <p className="text-2xl font-bold text-[#1A4B69]">
              {summary.ongoingAuctions}
            </p>
          </div>
          <div className="bg-white border border-gray-200 rounded-lg p-4">
            <p className="text-sm text-gray-500">Completed Auctions</p>
            <p className="text-2xl font-bold text-[#1A4B69]">
              {summary.completedAuctions}
            </p>
          </div>
          <div className="bg-white border border-gray-200 rounded-lg p-4">
            <p className="text-sm text-gray-500">Bids (Ongoing)</p>
            <p className="text-2xl font-bold text-[#1A4B69]">
              {summary.totalBidsOngoing}
            </p>
          </div>
          <div className="bg-white border border-gray-200 rounded-lg p-4">
            <p className="text-sm text-gray-500">Bids (Completed)</p>
            <p className="text-2xl font-bold text-[#1A4B69]">
              {summary.totalBidsCompleted}
            </p>
          </div>
          <div className="bg-white border border-gray-200 rounded-lg p-4">
            <p className="text-sm text-gray-500">Bids (All)</p>
            <p className="text-2xl font-bold text-[#1A4B69]">
              {summary.totalBidsAll}
            </p>
          </div>
        </div>
      )}
    </div>
  );
};

export default ReportAnalyticsPage;