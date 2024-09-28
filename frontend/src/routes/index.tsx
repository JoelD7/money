import { createFileRoute, redirect } from "@tanstack/react-router";
import { Home } from "../pages";
import { store } from "../store";
import { z } from "zod";
import { zodSearchValidator } from "@tanstack/router-zod-adapter";
import { keys } from "../utils";

function isAuth() {
  return store.getState().authReducer.isAuthenticated;
}

const expensesSearchSchema = z.object({
  categories: z.string().optional(),
  pageSize: z.number().default(10),
  startKey: z.string().optional(),
  period: z.string().default(localStorage.getItem(keys.CURRENT_PERIOD) || ""),
});

export const Route = createFileRoute("/")({
  beforeLoad: async ({ location }) => {
    if (!isAuth()) {
      throw redirect({
        to: "/login",
        search: {
          redirect: location.href,
        },
      });
    }
  },
  validateSearch: zodSearchValidator(expensesSearchSchema),
  component: Index,
});

function Index() {
  return <Home />;
}
