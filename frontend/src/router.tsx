import { createBrowserRouter, Navigate } from "react-router-dom";
import LoginPage from "./pages/auth/LoginPage";
import RegisterPage from "./pages/auth/RegisterPage";
import AuctionPage from "./pages/auction/AuctionPage";
import AuctionDetailPage from "./pages/detail-auction/AuctionDetailPage";
import ReportPage from "./pages/report/ReportPage";
import Layout from "./components/layout/Layout";
import NotFoundPage from "./pages/NotFoundPage";
import ReportAnalyticsPage from "./pages/report/ReportAnalyticsPage";
import ReportAuctionPage from "./pages/report/ReportAuctionPage";

const router = createBrowserRouter([
  {
    path: "/login",
    element: <LoginPage />,
  },
  {
    path: "/register",
    element: <RegisterPage />,
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
            element: <ReportAuctionPage />,
          },
        ],
      },
    ],
  },
  {
    path: "*",
    element: <NotFoundPage />,
  },
]);

export default router;
