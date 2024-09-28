import { useQuery } from "@tanstack/react-query";
import api from "../api";
import { keys } from "../utils";
import { AxiosError, AxiosResponse } from "axios";
import { User } from "../types";

export function useGetUser() {
  return useQuery({
    queryKey: ["user"],
    queryFn: () => {
      const result: Promise<AxiosResponse<User>> = api.getUser();
      result.then((res) => {
        localStorage.setItem(keys.CURRENT_PERIOD, res.data.current_period);
      });

      return result;
    },
  });
}

export function useGetPeriod(user? :User) {
  const periodID =
      user?.current_period || localStorage.getItem(keys.CURRENT_PERIOD) || "";

  return useQuery({
    queryKey: ["period"],
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
    queryKey: ["periodStats", periodID],
    queryFn: () => api.getPeriodStats(periodID),
    enabled: periodID !== "",
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
