import {
  Category,
  PeriodStats,
  SavingGoal,
  TransactionSearchParams,
  User,
} from "../types";
import { Colors } from "../assets";
import { useLocation } from "@tanstack/react-router";
import { CURRENT_PERIOD } from "./keys.ts"; // Sets category name and color to the categoryExpenseSummary object

// Sets category name and color to the categoryExpenseSummary object
export function setAdditionalData(
  periodStats: PeriodStats | undefined,
  user: User | undefined,
): PeriodStats | undefined {
  if (!periodStats || !user || !user.categories) {
    return periodStats;
  }

  const categoryByID: Map<string, Category> = new Map<string, Category>();
  user.categories.forEach((category) => {
    categoryByID.set(category.id, category);
  });

  periodStats.category_expense_summary.forEach((ces) => {
    const category = categoryByID.get(ces.category_id);
    if (category) {
      ces.name = category.name;
      ces.color = category.color;
      return;
    }

    ces.name = "Other";
    ces.color = Colors.GRAY_DARK;
  });

  return periodStats;
}

export function useTransactionsParams(): TransactionSearchParams {
  const location = useLocation();
  const params = new URLSearchParams(location.search);

  let categories: string[] = [];
  let param = params.get("categories");
  if (param !== null && param !== "") {
    categories = param.split(",");
  }

  let pageSize: number | undefined;

  param = params.get("pageSize");
  if (param !== null) {
    pageSize = parseInt(param);
  }

  let startKey: string | undefined;

  param = params.get("startKey");
  if (param !== null) {
    startKey = param;
  }

  let period: string | undefined;

  param = params.get("period");
  if (param !== null) {
    period = param;
  } else {
    period = localStorage.getItem(CURRENT_PERIOD) || "";
  }

  let sortBy: string | undefined;
  param = params.get("sortBy");
  if (param !== null) {
    sortBy = param;
  }

  let sortOrder: string | undefined;
  param = params.get("sortOrder");
  if (param !== null) {
    sortOrder = param;
  }

  return {
    categories,
    pageSize,
    startKey,
    period,
    sortBy,
    sortOrder,
  };
}

export function estimateSavingAmount(savingGoal?: SavingGoal): number {
  if (!savingGoal) return 0;

  const monthsUntilDeadline = getMonthDifference(new Date(), new Date(savingGoal.deadline));

  if (monthsUntilDeadline === 0) {
    return (savingGoal.target - savingGoal.progress);
  }

  const result = (savingGoal.target - savingGoal.progress) / (monthsUntilDeadline);

  //Use ceil because it's preferred for the amount to be over the goal than under it.
  return Math.ceil(result * 10) / 10;
}

export function getMonthDifference(current: Date, deadline: Date): number {
  const years: number = deadline.getFullYear() - current.getFullYear()
  if (years === 0) return deadline.getMonth() - current.getMonth()

  // for the case where the deadline is in the next year
  if(deadline.getMonth() < current.getMonth()) {
    return (deadline.getMonth() - current.getMonth() + 12)
  }

  if(deadline.getMonth() < current.getMonth() && years > 1) {
    return (deadline.getMonth() - current.getMonth() + 12) + (years * 12)
  }

  return (deadline.getMonth() - current.getMonth()) + (years * 12)
}