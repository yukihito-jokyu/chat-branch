import { apiClient } from "@/lib/api-client";

export type SignupResponse = {
  token: string;
  user: {
    id: string;
    name: string;
  };
};

export type LoginRequest = {
  user_id: string;
};

export type LoginResponse = {
  token: string;
};

export const signupGuest = async (): Promise<SignupResponse> => {
  return apiClient.post("/api/auth/signup");
};

export const loginGuest = async (
  data: LoginRequest
): Promise<LoginResponse> => {
  return apiClient.post("/api/auth/login", data);
};
