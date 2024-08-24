import { useQuery } from "@tanstack/react-query";
import api from "../api";
import { utils } from "../utils";

export function useGetUser() {
  return useQuery({
    queryKey: ["user"],
    queryFn: () => api.getUser(),
  });
}

export function useGetExpenses() {
  const { categories, pageSize, startKey, period } = utils.useExpensesParams();

  return useQuery({
    queryKey: api.expensesQueryKeys.list(
        categories,
        pageSize,
        startKey,
        period,
    ),
    queryFn: api.getExpenses,
  });
}

export function useGetPeriod() {
  return useQuery({
    queryKey: ["period"],
    queryFn: () => api.getPeriod(),
  });
}

export function useGetCategoryExpenseSummary(periodID: string = "current") {
  return useQuery({
    queryKey: ["categoryExpenseSummary"],
    queryFn: () => api.getCategoryExpenseSummary(periodID),
  });
}
