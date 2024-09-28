export type User = {
  username: string;
  current_period: string;
  remainder: number;
  expenses?: number;
  categories?: Category[];
};

export type ExpenseType = "regular" | "recurring";

export type Expense = {
  expense_id: string;
  username: string;
  category_id?: string;
  category_name?: string;
  amount: number;
  name: string;
  notes?: string;
  type: ExpenseType;
  created_date: string;
  period: string;
  update_date?: string;
};

export type PeriodStats = {
  period: string;
  total_income: number;
  category_expense_summary: CategoryExpenseSummary[];
};

export type CategoryExpenseSummary = {
  category_id: string;
  name: string;
  total: number;
  period: string;
  color: string;
};

export type Expenses = {
  expenses: Expense[];
  next_key: string;
};

export type Category = {
  id: string;
  name: string;
  budget: number;
  color: string;
};

export type SignUpUser = {
  username: string;
  password: string;
  fullname: string;
};

export type Credentials = {
  username: string;
  password: string;
};

export type APIError = {
  message: string;
  http_code: number;
};

export type AccessToken = {
  sub: string;
  exp: number;
  iat: number;
};

export type Period = {
  username: string;
  period: string;
  name: string;
  start_date: string;
  end_date: string;
  created_date: string;
  updated_date: string;
};

export type IncomeList = {
  income: Income[]
  next_key: string;
}

export type Income = {
  amount: number;
  name: string;
  period: string;
  notes: string;
  created_date: string;
};
