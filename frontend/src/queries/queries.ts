import { useInfiniteQuery, useQuery } from "@tanstack/react-query";
import api from "../api";
import { keys, utils } from "../utils";
import { PeriodList, User } from "../types";
import { INCOME, PERIOD, PERIOD_STATS, PERIODS, USER } from "./keys";
import { defaultStaleTime, queryRetryFn } from "./common.ts";

export const QUERY_RETRIES = 2;

export const incomeKeys = {
  all: [{ scope: INCOME }] as const,
  list: (
    pageSize?: number,
    startKey?: string,
    period?: string,
    sortOrder?: string,
    sortBy?: string,
  ) =>
    [
      {
        ...incomeKeys.all[0],
        pageSize,
        startKey,
        period,
        sortOrder,
        sortBy,
      },
    ] as const,
};

export const expensesQueryKeys = {
  all: [{ scope: "expenses" }] as const,
  list: (
      categories?: string[],
      pageSize?: number,
      startKey?: string,
      period?: string,
      sortBy?: string,
      sortOrder?: string,
  ) =>
      [
        {
          ...expensesQueryKeys.all[0],
          pageSize,
          startKey,
          period,
          categories,
          sortBy,
          sortOrder,
        },
      ] as const,
};

export function useGetUser() {
  return useQuery({
    queryKey: [USER],
    queryFn: () => {
      const result: Promise<User> = api.getUser();
      result.then((res) => {
        localStorage.setItem(keys.CURRENT_PERIOD, res.current_period);
      });

      return result;
    },
    staleTime: defaultStaleTime,
  });
}

export function useGetPeriod(user?: User) {
  const periodID =
    user?.current_period || localStorage.getItem(keys.CURRENT_PERIOD) || "";

  return useQuery({
    queryKey: [PERIOD],
    queryFn: () => api.getPeriod(periodID),
    enabled: periodID !== "",
    staleTime: defaultStaleTime,
    retry: queryRetryFn,
  });
}

export function useGetPeriods(startKey: string = "", pageSize: number = 10) {
  return useQuery({
    queryKey: [PERIODS, startKey, pageSize],
    queryFn: () => api.getPeriods(startKey, pageSize),
    staleTime: defaultStaleTime,
    retry: queryRetryFn,
  });
}

export function useGetPeriodsInfinite() {
  return useInfiniteQuery({
    queryKey: [PERIODS],
    initialPageParam: "",
    staleTime: defaultStaleTime,
    getNextPageParam: (lastPage: PeriodList) => {
      return lastPage.next_key !== "" ? lastPage.next_key : null;
    },
    queryFn: ({ pageParam }) => api.getPeriods(pageParam),
    retry: queryRetryFn,
  });
}

export function useGetPeriodStats(user?: User) {
  const periodID =
    user?.current_period || localStorage.getItem(keys.CURRENT_PERIOD) || "";

  return useQuery({
    queryKey: [PERIOD_STATS, periodID],
    queryFn: () => api.getPeriodStats(periodID),
    enabled: periodID !== "",
    staleTime: defaultStaleTime,
    retry: queryRetryFn,
  });
}

export function useGetIncome() {
  // eslint-disable-next-line prefer-const
  let { pageSize, startKey, period, sortOrder, sortBy } = utils.useTransactionsParams();

  if (!period) {
    period = localStorage.getItem(keys.CURRENT_PERIOD) || "";
  }

  return useQuery({
    queryKey: incomeKeys.list(pageSize, startKey, period, sortOrder, sortBy),
    queryFn: api.getIncomeList,
    staleTime: defaultStaleTime,
    retry: queryRetryFn,
  });
}

export function useGetExpenses(periodID: string) {
  // eslint-disable-next-line prefer-const
  let { categories, pageSize, startKey, period, sortOrder, sortBy } =
    utils.useTransactionsParams();

  if (!period) {
    period = periodID;
  }

  return useQuery({
    queryKey: expensesQueryKeys.list(
      categories,
      pageSize,
      startKey,
      period,
      sortBy,
      sortOrder,
    ),
    queryFn: api.getExpenses,
    enabled: periodID !== "",
    staleTime: defaultStaleTime,
    retry: queryRetryFn,
  });
}

export function useGetSavings(
  startKey: string = "",
  pageSize: number = 10,
  sortOrder: string,
  sortBy: string,
  savingGoalID?: string,
) {
  return useQuery({
    queryKey: api.savingsKeys.list(pageSize, startKey, sortOrder, sortBy, savingGoalID),
    queryFn: api.getSavings,
    staleTime: defaultStaleTime,
    retry: queryRetryFn,
  });
}
