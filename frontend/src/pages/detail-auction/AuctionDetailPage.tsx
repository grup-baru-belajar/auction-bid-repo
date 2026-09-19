import { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import type { AuctionDetail } from "../../types";
import { auctionApi, bidApi } from "../../services/api";
import { useAppSelector } from "../../store/hooks";
import { toast } from "react-hot-toast";
import PersonImage from "../../assets/person.png"


const formatRupiah = (value: number) =>
  "Rp " + new Intl.NumberFormat("id-ID").format(value);

const formatDate = (iso: string) =>
  new Date(iso).toLocaleDateString("en-GB", {
    day: "2-digit",
    month: "short",
    year: "numeric",
  });


function AuctionCard({ auction }: { auction: AuctionDetail }) {
  return (
    <div className="bg-white rounded-2xl shadow-md overflow-hidden w-full">
      <div className="relative h-44">
        <img
          src={auction.imageLink}
          alt={auction.auctionName}
          className="w-full h-full object-cover"
          onError={(e) => {
            (e.currentTarget as HTMLImageElement).src =
              "https://placehold.co/400x220/7c3aed/ffffff?text=Auction";
          }}
        />
        <div className="absolute inset-0 bg-gradient-to-t from-purple-700/80 via-pink-500/40 to-transparent" />
        <div className="absolute top-3 left-3 w-6 h-6 rounded-full bg-white/40 backdrop-blur-sm" />
        <p className="absolute bottom-10 left-3 right-3 text-white text-sm font-bold leading-snug">
          {auction.auctionName}
        </p>
        <p className="absolute bottom-3 left-3 right-3 text-white/80 text-xs leading-tight">
          {auction.description.slice(0, 60)}…
        </p>
      </div>

      <div className="p-4 space-y-2">
        <h2 className="text-base font-bold text-gray-800 leading-tight">
          {auction.auctionName}
        </h2>
        <p className="text-xs text-gray-500 mt-2">{formatDate(auction.endTime)}</p>
        <p className="text-sm font-semibold text-gray-700 mt-2">
          Start from {formatRupiah(auction.startingPrice)}
        </p>

        {
          auction.isCompleted && auction.bidWinner && (
            <div className="flex items-center gap-2 mt-3 bg-gray-50 rounded-xl px-3 py-2">
              <img
                src={PersonImage}
                alt="Winner"
                className="w-8 h-8 rounded-full object-cover ring-2 ring-white"
                onError={(e) => {
                  (e.currentTarget as HTMLImageElement).src =
                    "https://placehold.co/32x32/7c3aed/ffffff?text=P";
                }}
              />
              <div>
                <p className="text-[10px] text-gray-400 leading-none">
                  Winner
                </p>
                <p className="text-xs font-semibold text-gray-700 leading-tight">
                  {auction.bidWinner?.name}
                </p>
              </div>
            </div>
          )
        }


      </div>
    </div>
  );
}


function StatCard({
  label,
  value,
  unit,
  icon,
}: {
  label: string;
  value: number;
  unit: string;
  icon: React.ReactNode;
}) {
  return (
    <div className="flex-1 bg-white rounded-2xl border border-gray-100 shadow-sm p-5 space-y-1.5 min-w-0">
      <div className="flex items-center justify-between">
        <span className="text-xs text-gray-500 font-medium">{label}</span>
        <div className="w-9 h-9 rounded-xl bg-blue-50 flex items-center justify-center text-blue-500">
          {icon}
        </div>
      </div>
      <p className="text-3xl font-bold text-gray-900">
        {new Intl.NumberFormat("id-ID").format(value)}
        <span className="text-sm font-normal text-gray-400 ml-1.5">{unit}</span>
      </p>
    </div>
  );
}


function TopBidderTable({ bids }: { bids: AuctionDetail["topBids"] }) {
  return (
    <div className="bg-white rounded-2xl border border-gray-100 shadow-sm p-5">
      <div className="flex items-start justify-between mb-4">
        <div>
          <h3 className="text-sm font-bold text-gray-800">Top Bidder</h3>
          <p className="text-xs text-gray-400 mt-2">Highest Bidder right now</p>
        </div>
        <button
          id="view-all-bidders-btn"
          className="text-xs text-blue-500 hover:text-blue-700 font-medium transition-colors whitespace-nowrap"
        >
          See All Bidder →
        </button>
      </div>

      <div className="hidden lg:block overflow-x-auto">
        <table className="w-full text-sm">
          <thead>
            <tr className="bg-gray-50">
              <th className="text-left text-xs text-gray-500 font-medium px-4 py-2.5 rounded-l-lg">
                Bidder
              </th>
              <th className="text-center text-xs text-gray-500 font-medium px-4 py-2.5 rounded-r-lg">
                Highest Bid
              </th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-50">
            {bids.map((bid, index) => (
              <tr key={bid.id} className="hover:bg-gray-50/60 transition-colors">
                <td className="px-4 py-3 text-gray-700">
                  <div className="flex items-center gap-3">
                    <span className="w-6 h-6 rounded-full bg-blue-100 text-blue-600 text-xs font-bold flex items-center justify-center shrink-0">
                      {index + 1}
                    </span>
                    <p className="text-sm font-medium text-gray-700">{bid.userName}</p>
                  </div>
                </td>
                <td className="px-4 py-3 text-center font-semibold text-gray-800">
                  {formatRupiah(bid.bidPrice)}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <div className="flex flex-col gap-2 lg:hidden">
        {bids.map((bid, index) => (
          <div
            key={bid.id}
            className="flex items-center justify-between bg-gray-50 rounded-xl px-4 py-3"
          >
            <div className="flex items-center gap-3">
              <span className="w-6 h-6 rounded-full bg-blue-100 text-blue-600 text-xs font-bold flex items-center justify-center shrink-0">
                {index + 1}
              </span>
              <p className="text-sm font-medium text-gray-700">{bid.userName}</p>
            </div>
            <div className="text-right">
              <p className="text-sm font-semibold text-gray-800">
                {formatRupiah(bid.bidPrice)}
              </p>

            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

const AuctionDetailPage = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { isAuthenticated, user} = useAppSelector((state) => state.auth);
  const isAdmin = user?.role === "ADMIN";

  const [auction, setAuction] = useState<AuctionDetail | null>(null);
  const [loading, setLoading] = useState(true);
  const [bidInput, setBidInput] = useState("");
  const [bidError, setBidError] = useState("");
  const [bidLoading, setBidLoading] = useState(false);

  useEffect(() => {
    if (!id) return;
    let isMounted = true; 
    const fetchAuctionDetail = async () => {
      setLoading(true); 
      try {
        const res = await auctionApi.getAuctionDetail(Number(id));
        if (isMounted) {
          setAuction(res.data.data);
        }
      } catch {
        if (isMounted) {
          setAuction(null);
        }
      } finally {
        if (isMounted) {
          setLoading(false);
        }
      }
    };
    fetchAuctionDetail();
    return () => {
      isMounted = false; 
    };
  }, [id]);

  const handlePlaceBid = async () => {
    const price = Number(bidInput.replace(/\D/g, ""));
    if (!price || price <= 0) {
      setBidError("Input is not valid");
      return;
    }
    if (auction && price <= auction.lastPrice) {
      setBidError(`Bid value must above ${formatRupiah(auction.lastPrice)}`);
      return;
    }
    setBidError("");
    setBidLoading(true);
    try {
      await bidApi.createBid({
        auctionId: Number(id),
        bidPrice: price,
      });
      setBidInput("");
      // Refresh data setelah bid berhasil
      const res = await auctionApi.getAuctionDetail(Number(id));
      setAuction(res.data.data);
      toast.success("Bid Placed");
    } catch (err: unknown) {
      const msg =
        (err as { response?: { data?: { message?: string } } })?.response?.data
          ?.message ?? "Failed to place bid. Please try again.";
      setBidError(msg);
    } finally {
      setBidLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="flex flex-col items-center gap-3">
          <div className="w-10 h-10 border-4 border-blue-500 border-t-transparent rounded-full animate-spin" />
          <p className="text-sm text-gray-500">Loading</p>
        </div>
      </div>
    );
  }

  if (!auction) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center space-y-4">
          <p className="text-gray-600">Auction not found.</p>
          <button
            onClick={() => navigate("/")}
            className="px-4 py-2 bg-blue-600 text-white rounded-lg text-sm hover:bg-blue-700 transition-colors"
          >
            Back
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <header className="bg-white border-b border-gray-100 sticky top-0 z-10">
        <div className="max-w-6xl mx-auto px-6 py-3 flex items-center gap-3">
          <button
            id="back-btn"
            onClick={() => navigate("/")}
            className="flex items-center gap-1.5 text-sm text-gray-500 hover:text-gray-800 transition-colors"
          >
            <svg className="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5">
              <path d="M19 12H5M12 5l-7 7 7 7" strokeLinecap="round" strokeLinejoin="round" />
            </svg>
            Back
          </button>
          <span className="text-gray-300">/</span>
          <span className="text-sm font-semibold text-gray-700 truncate">
            {auction.auctionName}
          </span>
        </div>
      </header>

      <main className="max-w-6xl mx-auto px-6 py-8">
        <div className="flex flex-col lg:flex-row gap-6">

          <aside className="w-full lg:w-72 shrink-0">
            <AuctionCard auction={auction} />
          </aside>

          <div className="flex-1 flex flex-col gap-5 min-w-0">

            <div className="flex flex-col sm:flex-row gap-4">
              <StatCard
                label="Total Bids"
                value={auction.totalBids}
                unit="Bids"
                icon={
                  <svg className="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5">
                    <polyline points="22 7 13.5 15.5 8.5 10.5 2 17" />
                    <polyline points="16 7 22 7 22 13" />
                  </svg>
                }
              />
              <StatCard
                label="Total Bidders"
                value={auction.totalBidders}
                unit="Bidders"
                icon={
                  <svg className="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5">
                    <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2" />
                    <circle cx="9" cy="7" r="4" />
                    <path d="M23 21v-2a4 4 0 0 0-3-3.87" />
                    <path d="M16 3.13a4 4 0 0 1 0 7.75" />
                  </svg>
                }
              />
            </div>

            <TopBidderTable bids={auction.topBids} />

            {!auction.isCompleted && (
              <div className="space-y-1.5">
                {isAdmin ? (
                  <div className="flex items-center gap-3 bg-amber-50 p-4 rounded-xl border border-amber-200 shadow-sm">
                    <div className="w-10 h-10 rounded-full bg-amber-100 flex items-center justify-center text-amber-600 shrink-0">
                      <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                      </svg>
                    </div>
                    <div>
                      <p className="text-sm font-semibold text-amber-900">Admin Restricted</p>
                      <p className="text-xs text-amber-700 mt-0.5">Admin accounts are not allowed to place bids on auctions.</p>
                    </div>
                  </div>
                ) : isAuthenticated ? (
                  <div className="flex items-center gap-3">
                    <div className="flex-1 flex items-center border border-gray-200 rounded-xl overflow-hidden focus-within:ring-2 focus-within:ring-blue-500 focus-within:border-transparent transition-all bg-white shadow-sm">
                      <span className="pl-4 pr-2 text-sm font-medium text-gray-400 select-none">
                        Rp
                      </span>
                      <input
                        id="bid-amount-input"
                        type="text"
                        inputMode="numeric"
                        placeholder="Input your bid nominal"
                        value={bidInput}
                        disabled={bidLoading}
                        onChange={(e) => {
                          setBidInput(e.target.value);
                          setBidError("");
                        }}
                        className="flex-1 py-3 pr-4 text-sm text-gray-800 outline-none placeholder-gray-300 bg-transparent disabled:opacity-50"
                        aria-label="Nominal bid"
                      />
                    </div>
                    <button
                      id="place-bid-btn"
                      onClick={handlePlaceBid}
                      disabled={bidLoading}
                      className="px-6 py-3 bg-[#1A4B69] hover:bg-[#12364c] active:scale-95 text-white text-sm font-semibold rounded-xl transition-all shadow-sm whitespace-nowrap disabled:opacity-60 disabled:cursor-not-allowed flex items-center gap-2"
                    >
                      {bidLoading && (
                        <span className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
                      )}
                      {bidLoading ? "Placing..." : "Place Bid"}
                    </button>
                  </div>
                ) : (
                  <div className="flex flex-col sm:flex-row justify-between items-center bg-blue-50/50 p-4 rounded-xl border border-blue-100 gap-4 shadow-sm">
                    <div className="flex items-center gap-3">
                      <div className="w-10 h-10 rounded-full bg-blue-100 flex items-center justify-center text-blue-600 shrink-0">
                        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                        </svg>
                      </div>
                      <div>
                        <p className="text-sm font-semibold text-blue-900">Login Required</p>
                        <p className="text-xs text-blue-700 mt-0.5">You must be logged in to place a bid on this item.</p>
                      </div>
                    </div>
                    <button
                      onClick={() => navigate("/login")}
                      className="w-full sm:w-auto px-6 py-2.5 bg-[#1A4B69] hover:bg-[#12364c] text-white text-sm font-semibold rounded-lg transition-colors whitespace-nowrap"
                    >
                      Login to Bid
                    </button>
                  </div>
                )}

                {bidError && (
                  <p className="text-xs text-red-500 pl-1">{bidError}</p>
                )}
              </div>
            )}
          </div>
        </div>
      </main>
    </div>
  );
};

export default AuctionDetailPage;
