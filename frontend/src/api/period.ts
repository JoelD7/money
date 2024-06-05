import { Period } from "../types";
import { keys } from "../utils/index.ts";
import { API_BASE_URL, axiosClient } from "./money-api.ts";

export function getPeriod(period: string = "current") {
  return axiosClient.get<Period>(API_BASE_URL + `/periods/${period}`, {
    withCredentials: true,
    headers: {
      Auth: `Bearer ${localStorage.getItem(keys.ACCESS_TOKEN)}`,
    },
  });
}
