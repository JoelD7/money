import { AxiosResponse } from "axios";
import { API_BASE_URL, axiosClient } from "./money-api.ts";
import { IdempotencyKVP, SavingGoal, SavingGoalList } from "../types";
import { keys } from "../utils/index.ts";
import { SavingGoalsSchema } from "../types/domain.ts";
import { buildQueryParams, getIdempotencyKey, handleIdempotentRequest } from "./utils.ts";

export async function getSavingGoals(
  startKey: string = "",
  pageSize: number = 10,
  sortOrder: string = "",
  sortBy: string = "",
) {
  const params = buildQueryParams({ startKey, pageSize, sortOrder, sortBy });
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

export async function createSavingGoal(savingGoal: SavingGoal) {
  let accessToken = localStorage.getItem(keys.ACCESS_TOKEN);
  if (!accessToken) {
    accessToken = ""
  }

  const idempotenceKVP: IdempotencyKVP = getIdempotencyKey(savingGoal, accessToken, "");
  const p = axiosClient.post(API_BASE_URL + "/savings/goals", savingGoal, {
    withCredentials: true,
    headers: {
      Auth: `Bearer ${localStorage.getItem(keys.ACCESS_TOKEN)}`,
      "Idempotency-Key": idempotenceKVP.idempotencyKey,
    },
  });

  return handleIdempotentRequest(p, idempotenceKVP.encodedRequestBody);
}

export async function getSavingGoal(id: string) {
  const res: AxiosResponse = await axiosClient.get<SavingGoal>(
    API_BASE_URL + `/savings/goals/${id}`,
    {
      withCredentials: true,
      headers: {
        Auth: `Bearer ${localStorage.getItem(keys.ACCESS_TOKEN)}`,
      },
    },
  );

  return res.data;
}

export function updateSavingGoal(savingGoal: SavingGoal) {
  return axiosClient.put(
    API_BASE_URL + `/savings/goals/${savingGoal.saving_goal_id}`,
    savingGoal,
    {
      withCredentials: true,
      headers: {
        Auth: `Bearer ${localStorage.getItem(keys.ACCESS_TOKEN)}`,
      },
    },
  );
}

export function deleteSavingGoal(id: string) {
  return axiosClient.delete(API_BASE_URL + `/savings/goals/${id}`, {
    withCredentials: true,
    headers: {
      Auth: `Bearer ${localStorage.getItem(keys.ACCESS_TOKEN)}`,
    },
  });
}