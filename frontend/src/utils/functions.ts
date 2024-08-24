import {
  Category,
  CategoryExpenseSummary,
  ExpensesSearchParams,
  User,
} from "../types";
import { Colors } from "../assets";
import { useLocation } from "@tanstack/react-router";

// Sets category name and color to the categoryExpenseSummary object
export function setAdditionalData(
  categoryExpenseSummary: CategoryExpenseSummary[] | undefined,
  user: User | undefined,
): CategoryExpenseSummary[]{
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

export function useExpensesParams(): ExpensesSearchParams {
  const location = useLocation();
  const params = new URLSearchParams(location.search);

  const categories: string[] = params.get("categories")?.split(",") || [];
  const pageSize: number = parseInt(params.get("pageSize") || "10");
  const startKey: string = params.get("startKey") || "";
  const period: string = params.get("period") || "current";

  return { pageSize, categories, startKey, period };
}
