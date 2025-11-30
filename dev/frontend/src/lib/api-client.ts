import axios from "axios";

export const apiClient = axios.create({
  baseURL: "http://localhost:1323",
  headers: {
    "Content-Type": "application/json",
  },
});

apiClient.interceptors.response.use(
  (response) => {
    return response.data;
  },
  (error) => {
    return Promise.reject(error);
  }
);
