export type RechartsLabelProps = {
  cx: number;
  cy: number;
  midAngle: number;
  innerRadius: number;
  outerRadius: number;
  percent: number;
  index: number;
};

export type InputError = {
  username?: string;
  password?: string;
};

export type AccessTokenResponse = {
  accessToken: string;
};

export type CategoryExpense = {
  category: string;
  color: string;
  value: number;
};

export type SnackAlert = {
  open: boolean;
  type: "success" | "error";
  message: string;
};

export type ExpensesSearchParams = {
  categories: string[];
  pageSize: number;
  startKey: string;
  period: string;
};
