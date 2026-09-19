export type Role = "ADMIN" | "USER";

export interface User {
  id: number;
  name: string;
  username: string;
  role: Role;
}

export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginResponse {
  id: number;
  name: string;
  username: string;
  role: Role;
  token: string;
}
