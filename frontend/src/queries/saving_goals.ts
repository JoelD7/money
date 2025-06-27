import { useInfiniteQuery, useQuery } from "@tanstack/react-query";
import { queryRetryFn } from "./index.ts";
import api from "../api";
import { SavingGoalList } from "../types";

export const savingGoalKeys = {
  all: [{ scope: "savingGoals" }],
  single: (id: string) => [...savingGoalKeys.all, id],
  infinite: () => [...savingGoalKeys.all, "infinite"],
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
    staleTime: 2 * 60 * 1000,
  });
}

export function useGetSavingGoalsInfinite() {
  return useInfiniteQuery({
    // Infinite queries must use a different key to regular queries because data is stored differently.
    queryKey: savingGoalKeys.infinite(),
    initialPageParam: "",
    getNextPageParam: (lastPage: SavingGoalList) => {
      return lastPage.next_key !== "" ? lastPage.next_key : null;
    },
    queryFn: ({ pageParam }) => api.getSavingGoals(pageParam, 10, "", ""),
    retry: queryRetryFn,
    staleTime: 2 * 60 * 1000,
  });
}

export function useGetSavingGoal(id: string) {
  return useQuery({
    queryKey: savingGoalKeys.single(id),
    queryFn: () => api.getSavingGoal(id),
    retry: queryRetryFn,
    staleTime: 2 * 60 * 1000,
  });
}
