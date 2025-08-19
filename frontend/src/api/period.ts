import { Period, PeriodSchema, PeriodsSchema } from "../types";
import { keys } from "../utils/index.ts";
import { API_BASE_URL, axiosClient } from "./money-api.ts";
import { AxiosResponse } from "axios";
import { buildQueryParams } from "./utils.ts";

export async function getPeriod(period: string) {
  const res: AxiosResponse = await axiosClient.get<Period>(
    API_BASE_URL + `/periods/${period}`,
    {
      withCredentials: true,
      headers: {
        Auth: `Bearer ${localStorage.getItem(keys.ACCESS_TOKEN)}`,
      },
    },
  );

  try {
    return PeriodSchema.parse(res.data);
  } catch (e) {
    console.error("[money] - Error parsing GET period response", e);
    return Promise.reject(new Error("Invalid period data"));
  }
}

export async function getPeriods(startKey: string, pageSize: number = 10, active: boolean = false) {
  const params = buildQueryParams({ startKey, pageSize, active });
  let url = API_BASE_URL + "/periods";

  if (params.length > 0) {
    url += `?${params.join("&")}`;
  }

  const res: AxiosResponse = await axiosClient.get<Period[]>(url, {
    withCredentials: true,
    headers: {
      Auth: `Bearer ${localStorage.getItem(keys.ACCESS_TOKEN)}`,
    },
  });

  try {
    return PeriodsSchema.parse(res.data);
  } catch (e) {
    console.error("[money] - Error parsing GET periods response", e);
    return Promise.reject(e);
  }
}

// TODO: Update the current period in localstorage after creating a new one
export function createPeriod() {}
