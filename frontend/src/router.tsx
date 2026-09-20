import { createBrowserRouter, Navigate } from "react-router-dom";
import LoginPage from "./pages/auth/LoginPage";
import AuctionPage from "./pages/auction/AuctionPage";
import AuctionDetailPage from "./pages/detail-auction/AuctionDetailPage";
import ReportPage from "./pages/report/ReportPage";
import Layout from "./components/layout/Layout";
import NotFoundPage from "./pages/NotFoundPage";
import ReportAnalyticsPage from "./pages/report/ReportAnalyticsPage";

const router = createBrowserRouter([
  {
    path: "/login",
    element: <LoginPage />,
  },
  {
    path: "/auction/:id",
    element: <AuctionDetailPage />,
  },
  {
    element: <Layout />,
    children: [
      {
        path: "/",
        element: <AuctionPage />,
      },
      {
        path: "/report",
        element: <ReportPage />,
        children: [
          {
            index: true,
            element: <Navigate to="/report/analytics" replace />,
          },
          {
            path: "analytics",
            element: <ReportAnalyticsPage />,
          },
          {
            path: "auction",
            element: <AuctionPage />,
          },
        ],
      },
    ],
  },
  {
    path: "*",
    element: <NotFoundPage />,
  },
  {
    path: "/report",
    element: <ReportPage />,
    children: [
      {
        index: true,
        element: <Navigate to="/report/analytics" replace />,
      },
      {
        path: "analytics",
        element: <ReportAnalyticsPage />,
      },
      {
        path: "auction",
        element: <ReportAuctionPage />,
      },
    ],
  },
]);

export default router;
