import { useMemo, useState } from "react";
import { ArrowDown, ArrowUp, ChevronsUpDown } from "lucide-react";
import type { TopAuction } from "../../types";
import { formatRupiah } from "../common/formatters";

type SortKey =
  | "auctionName"
  | "startingPrice"
  | "highestBid"
  | "totalBid"
  | "bidders"
  | "status";
type SortDir = "asc" | "desc";

const STATUS_META: Record<string, { label: string; className: string }> = {
  ACTIVE: { label: "Active", className: "bg-sky-100 text-sky-700" },
  ENDED: { label: "Completed", className: "bg-slate-100 text-slate-600" },
};

const COLUMNS: { key: SortKey; label: string; align?: "right" }[] = [
  { key: "auctionName", label: "Auction" },
  { key: "startingPrice", label: "Starting Price", align: "right" },
  { key: "highestBid", label: "Highest Bid", align: "right" },
  { key: "totalBid", label: "Total Bids", align: "right" },
  { key: "bidders", label: "Bidders", align: "right" },
  { key: "status", label: "Status" },
];

interface TopAuctionsTableProps {
  data: TopAuction[];
  loading?: boolean;
}

const TopAuctionsTable = ({ data, loading }: TopAuctionsTableProps) => {
  const [sortKey, setSortKey] = useState<SortKey>("totalBid");
  const [sortDir, setSortDir] = useState<SortDir>("desc");

  const sorted = useMemo(() => {
    const copy = [...data];
    copy.sort((a, b) => {
      const av = a[sortKey];
      const bv = b[sortKey];
      const cmp =
        typeof av === "string"
          ? av.localeCompare(bv as string)
          : av - (bv as number);
      return sortDir === "asc" ? cmp : -cmp;
    });
    return copy;
  }, [data, sortKey, sortDir]);

  const handleSort = (key: SortKey) => {
    if (key === sortKey) {
      setSortDir((d) => (d === "asc" ? "desc" : "asc"));
    } else {
      setSortKey(key);
      setSortDir("desc");
    }
  };

  return (
    <div className="bg-white rounded-xl border border-slate-200 shadow-sm p-6">
      <div className="flex items-center justify-between mb-4">
        <div>
          <h3 className="text-base font-semibold text-slate-900">
            Top Auctions
          </h3>
          <p className="text-[13px] text-slate-500 mt-0.5">
            Auction dengan aktivitas bid tertinggi
          </p>
        </div>
      </div>

      {loading ? (
        <p className="text-center text-sm text-slate-400 py-10">Loading...</p>
      ) : sorted.length === 0 ? (
        <p className="text-center text-sm text-slate-400 py-10">
          No auctions yet.
        </p>
      ) : (
        <div className="overflow-x-auto">
          <table className="w-full min-w-[640px] text-sm">
            <thead>
              <tr className="bg-slate-50 rounded-md">
                {COLUMNS.map((col) => (
                  <th
                    key={col.key}
                    onClick={() => handleSort(col.key)}
                    className={`px-3 py-2.5 text-xs font-medium text-slate-500 cursor-pointer select-none first:rounded-l-md last:rounded-r-md ${
                      col.align === "right" ? "text-right" : "text-left"
                    }`}
                  >
                    <span
                      className={`inline-flex items-center gap-1 ${
                        col.align === "right" ? "flex-row-reverse" : ""
                      }`}
                    >
                      {col.label}
                      {sortKey === col.key ? (
                        sortDir === "asc" ? (
                          <ArrowUp size={12} />
                        ) : (
                          <ArrowDown size={12} />
                        )
                      ) : (
                        <ChevronsUpDown size={12} className="text-slate-300" />
                      )}
                    </span>
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {sorted.map((auction) => {
                const status = STATUS_META[auction.status] ?? {
                  label: auction.status,
                  className: "bg-slate-100 text-slate-600",
                };
                return (
                  <tr
                    key={auction.id}
                    className="border-b border-slate-100 last:border-0"
                  >
                    <td className="px-3 py-3.5">
                      <div className="font-medium text-slate-900">
                        {auction.auctionName}
                      </div>
                      <div className="text-xs text-slate-500">
                        #AUC-{String(auction.id).padStart(4, "0")}
                      </div>
                    </td>
                    <td className="px-3 py-3.5 text-right text-slate-700">
                      {formatRupiah(auction.startingPrice)}
                    </td>
                    <td className="px-3 py-3.5 text-right font-semibold text-slate-900">
                      {formatRupiah(auction.highestBid)}
                    </td>
                    <td className="px-3 py-3.5 text-right text-slate-700">
                      {auction.totalBid}
                    </td>
                    <td className="px-3 py-3.5 text-right text-slate-700">
                      {auction.bidders}
                    </td>
                    <td className="px-3 py-3.5">
                      <span
                        className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${status.className}`}
                      >
                        {status.label}
                      </span>
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
};

export default TopAuctionsTable;
