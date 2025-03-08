import { useInfiniteQuery, useQuery } from "@tanstack/react-query";
import { queryRetryFn } from "./index.ts";
import api from "../api";
import { SavingGoalList } from "../types";

export const savingGoalKeys = {
  all: [{ scope: "savingGoals" }],
  list: (pageSize?: number, startKey?: string, sortOrder?: string, sortBy?: string) => {
    return [
      {
        ...savingGoalKeys.all[0],
        pageSize,
        startKey,
        sortOrder,
        sortBy,
      },
    ];
  },
};

export function useGetSavingGoals(
  startKey?: string,
  pageSize?: number,
  sortOrder?: string,
  sortBy?: string,
) {
  return useQuery({
    queryKey: savingGoalKeys.list(pageSize, startKey, sortOrder, sortBy),
    queryFn: () => api.getSavingGoals(startKey, pageSize, sortOrder, sortBy),
    retry: queryRetryFn,
  });
}

export function useGetSavingGoalsInfinite() {
  return useInfiniteQuery({
    // Infinite queries must use a different key to regular queries because data is stored differently.
    queryKey: ["savingGoals", "infinite"],
    initialPageParam: "",
    getNextPageParam: (lastPage: SavingGoalList) => {
      return lastPage.next_key !== "" ? lastPage.next_key : null;
    },
    queryFn: ({ pageParam }) => api.getSavingGoals(pageParam, 10, "", ""),
    retry: queryRetryFn,
  });
}

export function useGetSavingGoal(id: string) {
  return useQuery({
    queryKey: ["savingGoal", id],
    queryFn: () => api.getSavingGoal(id),
    retry: queryRetryFn,
  });
}
