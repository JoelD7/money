import { login, logout, signUp } from "./auth.ts";
import { getUser } from "./money-api.ts";
import { createExpense, getExpenses, getPeriodStats } from "./expenses.ts";
import { getPeriod, getPeriods } from "./period.ts";
import { createIncome, getIncomeList } from "./income.ts";
import {
  createSavingGoal,
  deleteSavingGoal,
  getSavingGoal,
  getSavingGoals,
  updateSavingGoal,
} from "./saving_goals.ts";
import { createSaving, getSavings, savingsKeys } from "./savings.ts";

const api = {
  signUp,
  getUser,
  login,
  logout,
  getExpenses,
  getPeriod,
  createExpense,
  getPeriodStats,
  savingsKeys,
  createIncome,
  getIncomeList,
  getPeriods,
  getSavingGoals,
  createSavingGoal,
  getSavings,
  createSaving,
  getSavingGoal,
  updateSavingGoal,
  deleteSavingGoal,
};

export default api;
