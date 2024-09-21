import { API_BASE_URL, axiosClient } from "./money-api.ts";
import { Income, IncomeList } from "../types";
import { keys } from "../utils/index.ts";
import { QueryFunctionContext } from "@tanstack/react-query";

export const incomeKeys = {
  all: [{ scope: "income" }] as const,
  list: (pageSize?: number, startKey?: string, period?: string) =>
    [
      {
        ...incomeKeys.all[0],
        pageSize,
        startKey,
        period,
      },
    ] as const,
};

export function createIncome(income: Income) {
  return axiosClient.post(API_BASE_URL + "/income", income, {
    withCredentials: true,
    headers: {
      Auth: `Bearer ${localStorage.getItem(keys.ACCESS_TOKEN)}`,
    },
  });
}

export function getIncomeList({
  queryKey,
}: QueryFunctionContext<ReturnType<(typeof incomeKeys)["list"]>>) {
  const { pageSize, startKey, period } = queryKey[0];

  const paramArr: string[] = []

  if (period) {
    paramArr.push(`period=${period}`);
  }
  if (pageSize) {
    paramArr.push(`page_size=${pageSize}`);
  }

  if (startKey && startKey !== "") {
    paramArr.push(`start_key=${startKey}`);
  }

  const params: string = paramArr.join("&");

  return axiosClient.get<IncomeList>(API_BASE_URL + `/income?${params}`, {
    withCredentials: true,
    headers: {
      Auth: `Bearer ${localStorage.getItem(keys.ACCESS_TOKEN)}`,
    },
  });
}
