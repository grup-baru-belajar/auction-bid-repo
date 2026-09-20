import { createBrowserRouter } from "react-router-dom";
import LoginPage from "./pages/auth/LoginPage";
import RegisterPage from "./pages/auth/RegisterPage";
import AuctionPage from "./pages/auction/AuctionPage";
import AuctionDetailPage from "./pages/detail-auction/AuctionDetailPage";
import ReportPage from "./pages/report/ReportPage";
import Layout from "./components/layout/Layout";
import NotFoundPage from "./pages/NotFoundPage";

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
      },
    ],
  },
  {
    path: "*",
    element: <NotFoundPage />,
  },
]);

export default router;
