import { login, logout, signUp } from "./auth.ts";
import { getUser } from "./money-api.ts";
import {
  createExpense,
  expensesQueryKeys,
  getExpenses,
  getPeriodStats,
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
  getPeriodStats,
  expensesQueryKeys,
};

export default api;
