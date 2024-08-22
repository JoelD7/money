import { login, logout, signUp } from "./auth.ts";
import { getUser } from "./money-api.ts";
import {
  createExpense,
  getCategoryExpenseSummary,
  getExpenses,
} from "./expenses.ts";
import { getPeriod } from "./period.ts";

const api = {
  signUp,
  getUser,
  login,
  logout,
  getExpenses,
  getPeriod,
  createExpense,
  getCategoryExpenseSummary,
};

export default api;
