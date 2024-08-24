import { CategoryExpenseSummary, Expense, Expenses } from "../types";
import { keys } from "../utils";
import { API_BASE_URL, axiosClient } from "./money-api.ts";
import { QueryFunctionContext } from "@tanstack/react-query";

export const expensesQueryKeys = {
  all: [{ scope: "expenses" }] as const,
  list: (
    categories: string[],
    pageSize: number,
    startKey: string,
    period: string,
  ) =>
    [
      {
        ...expensesQueryKeys.all[0],
        pageSize,
        startKey,
        period,
        categories,
      },
    ] as const,
};

export function getExpenses({
  queryKey,
}: QueryFunctionContext<ReturnType<(typeof expensesQueryKeys)["list"]>>) {
  const { categories, pageSize, startKey, period } = queryKey[0];

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

export function getCategoryExpenseSummary(period: string = "current") {
  return axiosClient.get<CategoryExpenseSummary[]>(
    API_BASE_URL + `/expenses/stats/period/${period}`,
    {
      withCredentials: true,
      headers: {
        Auth: `Bearer ${localStorage.getItem(keys.ACCESS_TOKEN)}`,
      },
    },
  );
}
