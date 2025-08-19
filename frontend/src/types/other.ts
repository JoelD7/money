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
  title: string;
  message?: string;
};

export type TransactionSearchParams = {
  categories?: string[];
  pageSize?: number;
  startKey?: string;
  period?: string;
  sortBy?: string;
  sortOrder?: string;
  active?: boolean;
  savingGoalID?: string;
};

export type PaginationModel = {
  page: number;
  pageSize: number;
};

// IdempotencyKVP holds the encoded request body and the idempotency key values that are kept on local storage
export type IdempotencyKVP = {
  encodedRequestBody: string;
  idempotencyKey: string;
}