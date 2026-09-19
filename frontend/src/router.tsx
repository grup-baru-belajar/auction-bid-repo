import { createBrowserRouter } from "react-router-dom";
import LoginPage from "./pages/auth/LoginPage";
import RegisterPage from "./pages/auth/RegisterPage";
import AuctionPage from "./pages/auction/AuctionPage";
import AuctionDetailPage from "./pages/detail-auction/AuctionDetailPage";
import ReportPage from "./pages/report/ReportPage";

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
    path: "/",
    element: <AuctionPage />,
  },
  {
    path: "/auction/:id",
    element: <AuctionDetailPage />,
  },
  {
    path: "/report",
    element: <ReportPage />,
  },
]);

export default router;
