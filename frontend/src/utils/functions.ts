import {
  Category,
  CategoryExpenseSummary,
  TransactionSearchParams,
  User,
} from "../types";
import { Colors } from "../assets";
import { useLocation } from "@tanstack/react-router";

// Sets category name and color to the categoryExpenseSummary object
export function setAdditionalData(
  categoryExpenseSummary: CategoryExpenseSummary[] | undefined,
  user: User | undefined,
): CategoryExpenseSummary[] {
  if (!categoryExpenseSummary || !user || !user.categories) {
    return [];
  }

  const categoryByID: Map<string, Category> = new Map<string, Category>();
  user.categories.forEach((category) => {
    categoryByID.set(category.id, category);
  });

  categoryExpenseSummary.forEach((ces) => {
    const category = categoryByID.get(ces.category_id);
    if (category) {
      ces.name = category.name;
      ces.color = category.color;
      return;
    }

    ces.name = "Other";
    ces.color = Colors.GRAY_DARK;
  });

  return categoryExpenseSummary;
}

export function useTransactionsParams(): TransactionSearchParams {
  const location = useLocation();
  const params = new URLSearchParams(location.search);

  let categories: string[] = [];
  let param = params.get("categories")
  if (param !== null && param !== ""){
    categories = param.split(",")
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
