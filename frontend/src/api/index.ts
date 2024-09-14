import { login, logout, signUp } from "./auth.ts";
import { getUser } from "./money-api.ts";
import {
  createExpense,
  expensesQueryKeys,
  getCategoryExpenseSummary,
  getExpenses,
} from "./expenses.ts";
import { getPeriod } from "./period.ts";
import { createIncome } from "./income.ts";

const api = {
  signUp,
  getUser,
  login,
  logout,
  getExpenses,
  getPeriod,
  createExpense,
  getCategoryExpenseSummary,
  expensesQueryKeys,
  createIncome,
};

export default api;
