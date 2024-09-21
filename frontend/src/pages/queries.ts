import { useQuery } from "@tanstack/react-query";
import api from "../api";
import { utils } from "../utils";
import { AxiosError } from "axios";

export function useGetUser() {
  return useQuery({
    queryKey: ["user"],
    queryFn: () => api.getUser(),
  });
}

export function useGetExpenses() {
  const { categories, pageSize, startKey, period } =
    utils.useTransactionsParams();

  return useQuery({
    queryKey: api.expensesQueryKeys.list(
      categories,
      pageSize,
      startKey,
      period,
    ),
    queryFn: api.getExpenses,
    retry: (_, e: AxiosError) => {
      return e.response ? e.response.status !== 404 : true;
    },
  });
}

export function useGetPeriod() {
  return useQuery({
    queryKey: ["period"],
    queryFn: () => api.getPeriod(),
    retry: (_, e: AxiosError) => {
      return e.response ? e.response.status !== 404 : true;
    },
  });
}

export function useGetCategoryExpenseSummary(periodID: string = "current") {
  return useQuery({
    queryKey: ["categoryExpenseSummary"],
    queryFn: () => api.getCategoryExpenseSummary(periodID),
    retry: (_, e: AxiosError) => {
      return e.response ? e.response.status !== 404 : true;
    },
  });
}

export function useGetIncome() {
  const { pageSize, startKey, period } = utils.useTransactionsParams();

  return useQuery({
    queryKey: api.incomeKeys.list(pageSize, startKey, period),
    queryFn: api.getIncomeList,
    retry: (_, e: AxiosError) => {
      return e.response ? e.response.status !== 404 : true;
    },
  });
}
