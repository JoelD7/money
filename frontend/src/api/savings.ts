import { buildQueryParams } from "./utils.ts";
import { API_BASE_URL, axiosClient } from "./money-api.ts";
import { AxiosResponse } from "axios";
import { Saving, SavingList, SavingsSchema } from "../types";
import { keys } from "../utils/index.ts";

export async function getSavings(
  startKey: string = "",
  pageSize: number = 10,
  sortOrder: string,
  sortBy: string,
) {
  const params = buildQueryParams(startKey, pageSize, sortOrder, sortBy);

  let url = API_BASE_URL + `/savings`;
  if (params.length > 0) {
    url += `?${params.join("&")}`;
  }

  const res: AxiosResponse = await axiosClient.get<SavingList>(url, {
    withCredentials: true,
    headers: {
      Auth: `Bearer ${localStorage.getItem(keys.ACCESS_TOKEN)}`,
    },
  });

  try {
    return SavingsSchema.parse(res.data);
  } catch (e) {
    console.error("[money] - Error parsing GET savings response", e);
    return Promise.reject(new Error("Invalid savings data"));
  }
}

export async function createSaving(saving: Saving) {
  return axiosClient.post(API_BASE_URL + "/savings", saving, {
    withCredentials: true,
    headers: {
      Auth: `Bearer ${localStorage.getItem(keys.ACCESS_TOKEN)}`,
    },
  });
}
