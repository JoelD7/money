import { API_BASE_URL, axiosClient } from "./money-api.ts";
import { Income } from "../types";
import { keys } from "../utils/index.ts";

export function createIncome(income: Income) {
  return axiosClient.post(API_BASE_URL + "/income", income, {
    withCredentials: true,
    headers: {
      Auth: `Bearer ${localStorage.getItem(keys.ACCESS_TOKEN)}`,
    },
  });
}
