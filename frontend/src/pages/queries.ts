import { useQuery } from "@tanstack/react-query";
import api from "../api";

export function useGetUser() {
  return useQuery({
    queryKey: ["user"],
    queryFn: () => api.getUser(),
  });
}

export function useGetExpenses() {
  return useQuery({
    queryKey: ["expenses"],
    queryFn: () => api.getExpenses(),
  });
}

export function useGetPeriod() {
  return useQuery({
    queryKey: ["period"],
    queryFn: () => api.getPeriod(),
  });
}

export function useGetCategoryExpenseSummary(periodID:string="current"){
    return useQuery({
        queryKey: ["categoryExpenseSummary"],
        queryFn: () => api.getCategoryExpenseSummary(periodID),
    });
}