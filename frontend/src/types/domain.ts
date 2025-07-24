import { z } from "zod";

export const CategorySchema = z.object({
  id: z.string(),
  name: z.string(),
  budget: z.number().optional(),
  color: z.string(),
});

export type Category = z.infer<typeof CategorySchema>;

export const UserSchema = z.object({
  username: z.string(),
  current_period: z.string().optional(),
  remainder: z.number(),
  expenses: z.number().optional(),
  categories: z.array(CategorySchema).optional(),
});

export type User = z.infer<typeof UserSchema>;

export const ExpenseTypeSchema = z.enum(["regular", "recurring"]);
export type ExpenseType = z.infer<typeof ExpenseTypeSchema>;

export const ExpenseSchema = z.object({
  expense_id: z.string(),
  username: z.string(),
  category_id: z.string().optional(),
  category_name: z.string().optional(),
  amount: z.number(),
  name: z.string(),
  notes: z.string().optional(),
  type: ExpenseTypeSchema.optional(),
  created_date: z.string(),
  period_id: z.string(),
  update_date: z.string().optional(),
});

export type Expense = z.infer<typeof ExpenseSchema>;

export const CategoryExpenseSummarySchema = z.object({
  category_id: z.string(),
  total: z.number(),
  period: z.string().optional(),
  name: z.string().optional(),
  color: z.string().optional(),
});

export type CategoryExpenseSummary = z.infer<typeof CategoryExpenseSummarySchema>;

export const PeriodStatsSchema = z.object({
  period_id: z.string(),
  total_income: z.number(),
  category_expense_summary: z.array(CategoryExpenseSummarySchema),
});

export type PeriodStats = z.infer<typeof PeriodStatsSchema>;

export const ExpensesSchema = z.object({
  expenses: z.array(ExpenseSchema),
  next_key: z.string(),
});

export type Expenses = z.infer<typeof ExpensesSchema>;

export const SignUpUserSchema = z.object({
  username: z.string(),
  password: z.string(),
  fullname: z.string(),
});

export type SignUpUser = z.infer<typeof SignUpUserSchema>;

export const CredentialsSchema = z.object({
  username: z.string(),
  password: z.string(),
});

export type Credentials = z.infer<typeof CredentialsSchema>;

export type APIError = {
  message: string;
  http_code: number;
};

export const AccessTokenSchema = z.object({
  sub: z.string(),
  exp: z.number(),
  iat: z.number(),
});

export type AccessToken = z.infer<typeof AccessTokenSchema>;

export const PeriodSchema = z.object({
  username: z.string(),
  period: z.string(),
  name: z.string(),
  start_date: z.string(),
  end_date: z.string(),
  created_date: z.string(),
  updated_date: z.string(),
});

export const PeriodsSchema = z.object({
  periods: z.array(PeriodSchema),
  next_key: z.string(),
});

export type Period = z.infer<typeof PeriodSchema>;
export type PeriodList = z.infer<typeof PeriodsSchema>;

export const IncomeSchema = z.object({
  income_id: z.string(),
  amount: z.number(),
  name: z.string(),
  period_id: z.string(),
  notes: z.string().optional(),
  created_date: z.string(),
});

export type Income = z.infer<typeof IncomeSchema>;

export const IncomeListSchema = z.object({
  income: z.array(IncomeSchema),
  next_key: z.string(),
  periods: z.array(z.string()).optional(),
});

export type IncomeList = z.infer<typeof IncomeListSchema>;

export const SavingSchema = z.object({
  saving_id: z.string(),
  saving_goal_id: z.string().optional(),
  saving_goal_name: z.string().optional(),
  username: z.string(),
  period_id: z.string(),
  created_date: z.string().optional(),
  updated_date: z.string().optional(),
  amount: z.number(),
});

export const SavingsSchema = z.object({
  savings: z.array(SavingSchema),
  next_key: z.string(),
});

export type Saving = z.infer<typeof SavingSchema>;
export type SavingList = z.infer<typeof SavingsSchema>;

export const SavingGoalSchema = z.object({
  saving_goal_id: z.string(),
  name: z.string(),
  target: z.number(),
  progress: z.number(),
  deadline: z.string(),
  username: z.string(),
  is_recurring: z.boolean().optional(),
  recurring_amount: z.number().optional(),
});

export const SavingGoalsSchema = z.object({
  saving_goals: z.array(SavingGoalSchema),
  next_key: z.string(),
});

export type SavingGoal = z.infer<typeof SavingGoalSchema>;
export type SavingGoalList = z.infer<typeof SavingGoalsSchema>;
