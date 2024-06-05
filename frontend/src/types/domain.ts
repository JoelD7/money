export type User = {
  username: string;
  current_period: string;
  remainder: number;
  expenses?: number;
  categories?: Category[];
};

export type Expense = {
  expense_id: string;
  username: string;
  category_id?: string;
  category_name?: string;
  amount: number;
  name: string;
  notes?: string;
  created_date: Date;
  period: string;
  update_date: Date;
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
  start_date: Date;
  end_date: Date;
  created_date: Date;
  updated_date: Date;
};
