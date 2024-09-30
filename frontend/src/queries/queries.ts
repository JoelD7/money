import { useQuery } from "@tanstack/react-query";
import api from "../api";
import { keys, utils } from "../utils";
import { AxiosError, AxiosResponse } from "axios";
import { User } from "../types";
import { INCOME, PERIOD, PERIOD_STATS, USER } from "./keys";

export const incomeKeys = {
  all: [{ scope: INCOME }] as const,
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

export function useGetUser() {
  return useQuery({
    queryKey: [USER],
    queryFn: () => {
      const result: Promise<AxiosResponse<User>> = api.getUser();
      result.then((res) => {
        localStorage.setItem(keys.CURRENT_PERIOD, res.data.current_period);
      });

      return result;
    },
  });
}

export function useGetPeriod(user?: User) {
  const periodID =
    user?.current_period || localStorage.getItem(keys.CURRENT_PERIOD) || "";

  return useQuery({
    queryKey: [PERIOD],
    queryFn: () => api.getPeriod(periodID),
    enabled: periodID !== "",
    retry: (_, e: AxiosError) => {
      return e.response ? e.response.status !== 404 : true;
    },
  });
}

export function useGetPeriodStats(user?: User) {
  const periodID =
    user?.current_period || localStorage.getItem(keys.CURRENT_PERIOD) || "";

  return useQuery({
    queryKey: [PERIOD_STATS, periodID],
    queryFn: () => api.getPeriodStats(periodID),
    enabled: periodID !== "",
    retry: (_, e: AxiosError) => {
      return e.response ? e.response.status !== 404 : true;
    },
  });
}

export function useGetIncome(periodID: string) {
  // eslint-disable-next-line prefer-const
  let { pageSize, startKey, period } = utils.useTransactionsParams();

  if (!period) {
    period = periodID;
  }

  return useQuery({
    queryKey: incomeKeys.list(pageSize, startKey, period),
    queryFn: api.getIncomeList,
    retry: (_, e: AxiosError) => {
      return e.response ? e.response.status !== 404 : true;
    },
  });
}
