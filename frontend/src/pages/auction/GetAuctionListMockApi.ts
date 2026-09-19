import type { Auction, GetAuctionsParams, Pagination } from "../../types";

// ============================================
// Data dummy
// ============================================
const productNames = [
  "iPhone 15 Pro",
  "Samsung Galaxy S24",
  "MacBook Air M2",
  "Sony WH-1000XM5",
  "Nintendo Switch OLED",
  "iPad Pro 12.9",
  "Canon EOS R6",
  "PlayStation 5",
  "Dell XPS 13",
  "Apple Watch Ultra",
];

function generateMockAuctions(count: number): Auction[] {
  const auctions: Auction[] = [];

  for (let i = 1; i <= count; i++) {
    const startingPrice = (Math.floor(Math.random() * 20) + 1) * 1_000_000;
    const hasBid = Math.random() > 0.3;
    const isCompleted = i % 7 === 0;

    auctions.push({
      id: i,
      auctionName: `${productNames[i % productNames.length]} #${i}`,
      description: `Deskripsi singkat untuk barang lelang nomor ${i}.`,
      imageLink: `https://picsum.photos/seed/${i}/400/300`,
      startingPrice,
      lastPrice: hasBid
        ? startingPrice + Math.floor(Math.random() * 5_000_000)
        : startingPrice,
      createdAt: new Date(Date.now() - i * 3600_000).toISOString(),
      endTime: new Date(Date.now() + (30 - i) * 3600_000).toISOString(),
      isCompleted,
      bidWinner: hasBid ? { id: 100 + i, name: `User ${100 + i}` } : null,
    });
  }

  return auctions;
}

const mockAuctions: Auction[] = generateMockAuctions(30);

// ============================================
// GET /auctions
// ============================================
interface MockAuctionsResponse {
  data: Auction[];
  pagination: Pagination;
}

export function mockGetAuctions(
  params?: GetAuctionsParams,
): Promise<MockAuctionsResponse> {
  const page = params?.page ?? 1;
  const limit = params?.limit ?? 10;

  let filtered = mockAuctions;
  if (params?.isCompleted !== undefined) {
    filtered = mockAuctions.filter((a) => a.isCompleted === params.isCompleted);
  }

  const total = filtered.length;
  const totalPages = Math.ceil(total / limit);
  const startIndex = (page - 1) * limit;
  const paginatedData = filtered.slice(startIndex, startIndex + limit);

  const pagination: Pagination = { page, limit, total, totalPages };

  return new Promise((resolve) => {
    setTimeout(() => resolve({ data: paginatedData, pagination }), 400);
  });
}
