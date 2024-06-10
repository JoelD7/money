import { Expense, Expenses } from "../types";
import { keys } from "../utils";
import { API_BASE_URL, axiosClient } from "./money-api.ts";

export function getExpenses(
  period: string = "current",
  categories: string[] = [],
  startKey: string = "",
  pageSize: number = 10,
) {
  let params: string = `period=${period}&page_size=${pageSize}`;
  if (categories.length > 0) {
    for (let i = 0; i < categories.length; i++) {
      params += `&category=${categories[i]}`;
    }
  }

  if (startKey !== "") {
    params += `&start_key=${startKey}`;
  }

  return axiosClient.get<Expenses>(API_BASE_URL + `/expenses?${params}`, {
    withCredentials: true,
    headers: {
      Auth: `Bearer ${localStorage.getItem(keys.ACCESS_TOKEN)}`,
    },
  });
}

export function createExpense(expense: Expense) {
  return axiosClient.post(API_BASE_URL + "/expenses", expense, {
    withCredentials: true,
    headers: {
      Auth: `Bearer ${localStorage.getItem(keys.ACCESS_TOKEN)}`,
    },
  });
}
