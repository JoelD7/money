import { API_BASE_URL, axiosClient } from "./money-api.ts";
import { IdempotencyKVP, Income, IncomeList, IncomeListSchema } from "../types";
import { keys } from "../utils/index.ts";
import { QueryFunctionContext } from "@tanstack/react-query";
import { incomeKeys } from "../queries";
import { AxiosResponse } from "axios";
import { getIdempotencyKey, handleIdempotentRequest } from "./utils.ts";

export function createIncome(income: Income) {
  let accessToken = localStorage.getItem(keys.ACCESS_TOKEN);
  if (!accessToken) {
    accessToken = "";
  }

  const idempotenceKVP: IdempotencyKVP = getIdempotencyKey(income, accessToken, "");
  const p = axiosClient.post(API_BASE_URL + "/income", income, {
    withCredentials: true,
    headers: {
      Auth: `Bearer ${localStorage.getItem(keys.ACCESS_TOKEN)}`,
      "Idempotency-Key": idempotenceKVP.idempotencyKey,
    },
  });

  return handleIdempotentRequest(p, idempotenceKVP.encodedRequestBody);
}

export async function getIncomeList({
                                        queryKey,
                                    }: QueryFunctionContext<ReturnType<(typeof incomeKeys)["list"]>>) {
    const {pageSize, startKey, period, sortOrder, sortBy} = queryKey[0];

    const paramArr: string[] = [];

    if (period) {
        paramArr.push(`period=${period}`);
    }
    if (pageSize) {
        paramArr.push(`page_size=${pageSize}`);
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

    const res: AxiosResponse = await axiosClient.get<IncomeList>(API_BASE_URL + `/income?${params}`, {
        withCredentials: true,
        headers: {
            Auth: `Bearer ${localStorage.getItem(keys.ACCESS_TOKEN)}`,
        },
    });

    try {
        return IncomeListSchema.parse(res.data);
    } catch (e) {
        console.error("[money] - Error parsing GET income response", e)
        return Promise.reject(new Error("Invalid income data"))
    }
}
