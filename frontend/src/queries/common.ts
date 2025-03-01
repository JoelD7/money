import { AxiosError } from "axios";
import { QUERY_RETRIES } from "./queries.ts";

export function queryRetryFn(failureCount: number, e: AxiosError) {
  if (failureCount > QUERY_RETRIES) {
    return false;
  }

  if (e.response) {
    // Only retry server errors
    return e.response.status >= 500;
  }

  return false;
}
