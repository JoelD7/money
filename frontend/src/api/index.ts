import { login, logout, signUp } from "./auth.ts";
import { getUser } from "./money-api.ts";
import {
  createExpense,
  expensesQueryKeys,
  getExpenses,
  getPeriodStats,
} from "./expenses.ts";
import { getPeriod, getPeriods } from "./period.ts";
import { createIncome, getIncomeList } from "./income.ts";
import { createSavingGoal, getSavingGoal, getSavingGoals } from "./saving_goals.ts";
import { createSaving, getSavings } from "./savings.ts";

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
  createIncome,
  getIncomeList,
  getPeriods,
  getSavingGoals,
  createSavingGoal,
  getSavings,
  createSaving,
  getSavingGoal,
};

export default api;
