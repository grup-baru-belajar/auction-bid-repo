import apiClient from "./apiClient";
import type { ApiResponse, LoginRequest, LoginResponse } from "../../types";

export const authApi = {
  login: (data: LoginRequest) =>
    apiClient.post<ApiResponse<LoginResponse>>("/login", data),
};
