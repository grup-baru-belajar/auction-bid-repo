// src/components/auction/AuctionCard.tsx
import { Link } from "react-router-dom";
import type { Auction } from "../../types";

interface AuctionCardProps {
  auction: Auction;
}
function formatDateRange(createdAt: string, endTime: string): string {
  const options: Intl.DateTimeFormatOptions = {
    day: "numeric",
    month: "short",
    year: "numeric",
  };
  const start = new Date(createdAt).toLocaleDateString("id-ID", options);
  const end = new Date(endTime).toLocaleDateString("id-ID", options);
  return `${start} - ${end}`;
}

function formatRupiah(value: number | string): string {
  const num = typeof value === "string" ? Number(value) : value;
  return `Rp${num.toLocaleString("id-ID")}`;
}

function getInitial(name: string): string {
  return name.charAt(0).toUpperCase();
}

const AuctionCard = ({ auction }: AuctionCardProps) => {
  const hasBidWinner = auction.bidWinner !== null;

  return (
    <Link to={`/auction/${auction.id}`}>
      <div className="rounded-lg border border-gray-200 overflow-hidden shadow-sm hover:shadow-md transition-shadow bg-white flex flex-col">
        {/* Gambar + status dot */}
        <div className="relative px-6 pt-6">
          <img
            src={auction.imageLink}
            alt={auction.auctionName}
            className="w-full h-56 object-cover rounded-lg"
          />
        </div>

        {/* Konten teks */}
        <div className="pt-1 px-6 flex flex-col flex-1">
          <div className="p-3">
            <h3 className="font-semibold text-gray-900 line-clamp-1">
              {auction.auctionName}
            </h3>

            <p className="text-xs text-gray-400 mt-1">
              {formatDateRange(auction.createdAt, auction.endTime)}
            </p>

            <p className="text-sm text-gray-500 mt-2.5">
              from{" "}
              <span className="font-semibold text-gray-800">
                {formatRupiah(auction.startingPrice)}
              </span>
            </p>
          </div>

          {/* Footer kondisional */}
          <div className="bg-[#CBD5E1]  mb-6 border-gray-100 flex items-center gap-2 rounded-lg">
            {hasBidWinner ? (
              <>
                <div className="flex p-3">
                  <div className="mr-1 w-7 h-7 rounded-full bg-blue-100 text-blue-700 text-xs font-semibold flex items-center justify-center flex-shrink-0">
                    {getInitial(auction.bidWinner!.name)}
                  </div>
                  <div className="min-w-0">
                    <p className="text-[11px] text-[#475569] leading-tight">
                      Ditawar tertinggi oleh
                    </p>
                    <p className="text-xs font-medium text-black truncate">
                      {auction.bidWinner!.name}
                    </p>
                  </div>
                </div>
              </>
            ) : (
              <div className="min-w-0 p-3">
                <p className="text-[11px] text-[#475569] leading-tight">
                  Harga saat ini
                </p>
                <p className="text-xs font-medium text-black">
                  {formatRupiah(auction.lastPrice)}
                </p>
              </div>
            )}
          </div>
        </div>
      </div>
    </Link>
  );
};

export default AuctionCard;
