import { Category, ExpensesSearchParams, PeriodStats, User } from "../types";
import { Colors } from "../assets";
import { useLocation } from "@tanstack/react-router"; // Sets category name and color to the categoryExpenseSummary object

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

export function useExpensesParams(): ExpensesSearchParams {
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
  }

  return {
    categories,
    pageSize,
    startKey,
    period,
  };
}
