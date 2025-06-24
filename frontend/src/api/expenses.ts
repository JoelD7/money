import {
  Expense,
  Expenses,
  ExpensesSchema,
  IdempotencyKVP,
  PeriodStats,
  PeriodStatsSchema,
} from "../types";
import { keys } from "../utils";
import { API_BASE_URL, axiosClient } from "./money-api.ts";
import { QueryFunctionContext } from "@tanstack/react-query";
import { AxiosResponse } from "axios";
import { getIdempotencyKey, handleIdempotentRequest } from "./utils.ts";

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

export async function getExpenses({
  queryKey,
}: QueryFunctionContext<ReturnType<(typeof expensesQueryKeys)["list"]>>) {
  const { categories, pageSize, startKey, period, sortBy, sortOrder } = queryKey[0];

  const paramArr: string[] = [];

  if (period) {
    paramArr.push(`period=${period}`);
  }
  if (pageSize) {
    paramArr.push(`page_size=${pageSize}`);
  }

  if (categories && categories.length > 0) {
    for (let i = 0; i < categories.length; i++) {
      paramArr.push(`category=${categories[i]}`);
    }
  }

  if (startKey && startKey !== "") {
    paramArr.push(`start_key=${startKey}`);
  }

  if (sortBy && sortBy !== "") {
    paramArr.push(`sort_by=${sortBy}`);
  }

  if (sortOrder && sortOrder !== "") {
    paramArr.push(`sort_order=${sortOrder}`);
  }

  const params: string = paramArr.join("&");

  const res: AxiosResponse = await axiosClient.get<Expenses>(
    API_BASE_URL + `/expenses?${params}`,
    {
      withCredentials: true,
      headers: {
        Auth: `Bearer ${localStorage.getItem(keys.ACCESS_TOKEN)}`,
      },
    },
  );

  if (res.status === 204) {
    const emtpyResponse: Expenses = {
      expenses: [],
      next_key: "",
    };

    return emtpyResponse;
  }

  try {
    return ExpensesSchema.parse(res.data);
  } catch (e) {
    console.error("[money] - Error parsing expenses response", e);
    return Promise.reject(new Error("Invalid expenses data"));
  }
}

export function createExpense(expense: Expense) {
  let accessToken = localStorage.getItem(keys.ACCESS_TOKEN);
  if (!accessToken) {
    accessToken = ""
  }

  const idempotenceKVP: IdempotencyKVP = getIdempotencyKey(expense, accessToken, "");

  const p = axiosClient.post(API_BASE_URL + "/expenses", expense, {
    withCredentials: true,
    headers: {
      Auth: `Bearer ${accessToken}`,
      "Idempotency-Key": idempotenceKVP.idempotencyKey,
    },
  });

  return handleIdempotentRequest(p, idempotenceKVP.encodedRequestBody);
}

export async function getPeriodStats(period: string) {
  const res: AxiosResponse = await axiosClient.get<PeriodStats>(
    API_BASE_URL + `/periods/${period}/stats`,
    {
      withCredentials: true,
      headers: {
        Auth: `Bearer ${localStorage.getItem(keys.ACCESS_TOKEN)}`,
      },
    },
  );

  try {
    return PeriodStatsSchema.parse(res.data);
  } catch (e) {
    console.error("[money] - Error parsing period stats response", e);
    return Promise.reject(new Error("Invalid period stats data"));
  }
}
