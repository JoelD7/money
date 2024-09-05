import { useQuery } from "@tanstack/react-query";
import api from "../../api";
import { utils } from "../../utils";
import { AxiosError } from "axios";

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
    retry: (_, e: AxiosError) => {
      return e.response ? e.response.status !== 404 : true;
    },
  });
}
