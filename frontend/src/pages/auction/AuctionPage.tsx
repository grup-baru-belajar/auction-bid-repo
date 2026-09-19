import { useState, useEffect } from "react";
import { mockGetAuctions } from "./GetAuctionListMockApi";
import type { Auction, Pagination as PaginationType } from "../../types";
import AuctionCard from "../../components/auction/AuctionCard";
import Pagination from "../../components/common/Pagination";

const LIMIT = 10;

const AuctionPage = () => {
  const [auctions, setAuctions] = useState<Auction[]>([]);
  const [pagination, setPagination] = useState<PaginationType | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const [currentPage, setCurrentPage] = useState(1);
  const [isCompletedFilter, setIsCompletedFilter] = useState<
    boolean | undefined
  >(undefined);

  useEffect(() => {
    let isCancelled = false;

    async function loadAuctions() {
      setLoading(true);
      setError(null);

      try {
        const response = await mockGetAuctions({
          page: currentPage,
          limit: LIMIT,
          isCompleted: isCompletedFilter,
        });

        if (!isCancelled) {
          setAuctions(response.data);
          setPagination(response.pagination);
        }
      } catch (err) {
        if (!isCancelled) setError("Gagal memuat data auction");
      } finally {
        if (!isCancelled) setLoading(false);
      }
    }

    loadAuctions();
    return () => {
      isCancelled = true;
    };
  }, [currentPage, isCompletedFilter]);

  const handleFilterChange = (value: string) => {
    if (value === "all") setIsCompletedFilter(undefined);
    else setIsCompletedFilter(value === "completed");
    setCurrentPage(1);
  };

  const startItem =
    pagination && auctions.length > 0 ? (currentPage - 1) * LIMIT + 1 : 0;
  const endItem = pagination ? (currentPage - 1) * LIMIT + auctions.length : 0;

  return (
    <div className="min-h-screen bg-gray-100 p-6">
      <div className="max-w-6xl mx-auto">
        <div className="flex justify-between items-center mb-4">
          <p className="text-sm text-gray-500">
            Menampilkan <span className="text-black">{startItem}</span>-
            <span className="text-black">{endItem}</span> dari total{" "}
            <span className="text-black">{pagination?.total ?? 0}</span> Auction
          </p>

          <select
            onChange={(e) => handleFilterChange(e.target.value)}
            className="border border-gray-300 rounded-lg px-4 py-2 text-sm bg-white"
          >
            <option value="all">All</option>
            <option value="active">Active</option>
            <option value="completed">Completed</option>
          </select>
        </div>

        {loading && (
          <p className="text-center text-gray-500 py-10">Loading Data...</p>
        )}
        {error && <p className="text-center text-red-500 py-10">{error}</p>}

        {!loading && !error && auctions.length === 0 && (
          <p className="text-center text-gray-500 py-10">No Auction Found.</p>
        )}

        {!loading && !error && auctions.length > 0 && (
          <>
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
              {auctions.map((auction) => (
                <AuctionCard key={auction.id} auction={auction} />
              ))}
            </div>

            <Pagination
              currentPage={currentPage}
              totalPages={pagination?.totalPages ?? 1}
              onPageChange={setCurrentPage}
            />
          </>
        )}
      </div>
    </div>
  );
};

export default AuctionPage;
