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

export function useGetExpenses(periodID: string) {
  // eslint-disable-next-line prefer-const
  let { categories, pageSize, startKey, period } = utils.useExpensesParams();

  if (!period){
    period = periodID;
  }

  return useQuery({
    queryKey: api.expensesQueryKeys.list(
      categories,
      pageSize,
      startKey,
      period,
    ),
    queryFn: api.getExpenses,
    enabled: periodID !== "",
    retry: (_, e: AxiosError) => {
      return e.response ? e.response.status !== 404 : true;
    },
  });
}
