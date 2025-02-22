import { AxiosResponse } from "axios";
import { API_BASE_URL, axiosClient } from "./money-api.ts";
import { SavingGoalList } from "../types";
import { keys } from "../utils/index.ts";
import { SavingGoalsSchema } from "../types/domain.ts";

export async function getSavingGoals(startKey: string = "", pageSize: number = 10) {
  const params = [];

  if (startKey) {
    params.push(`start_key=${startKey}`);
  }

  if (pageSize) {
    params.push(`page_size=${pageSize}`);
  }

  let url = API_BASE_URL + `/savings/goals`;
  if (params.length > 0) {
    url += `?${params.join("&")}`;
  }

  const res: AxiosResponse = await axiosClient.get<SavingGoalList>(url, {
    withCredentials: true,
    headers: {
      Auth: `Bearer ${localStorage.getItem(keys.ACCESS_TOKEN)}`,
    },
  });

  try {
    return SavingGoalsSchema.parse(res.data);
  } catch (e) {
    console.error("[money] - Error parsing GET saving goals response", e);
    return Promise.reject(new Error("Invalid saving goals data"));
  }
}
