import { AxiosResponse } from "axios";
import { API_BASE_URL, axiosClient } from "./money-api.ts";
import { SavingGoalList } from "../types";
import { keys } from "../utils/index.ts";
import { SavingGoalsSchema } from "../types/domain.ts";

export async function getSavingGoals() {
  const res: AxiosResponse = await axiosClient.get<SavingGoalList>(
    API_BASE_URL + "/savings/goals",
    {
      withCredentials: true,
      headers: {
        Auth: `Bearer ${localStorage.getItem(keys.ACCESS_TOKEN)}`,
      },
    },
  );

  try {
    return SavingGoalsSchema.parse(res.data);
  } catch (e) {
    console.error("[money] - Error parsing GET saving goals response", e);
    return Promise.reject(new Error("Invalid saving goals data"));
  }
}
