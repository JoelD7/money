import { createFileRoute } from "@tanstack/react-router";
import { IncomeTable } from "../pages";
import { z } from "zod";
import { keys } from "../utils";
import { zodSearchValidator } from "@tanstack/router-zod-adapter";

const incomeSearchSchema = z.object({
  pageSize: z.number().default(10),
  startKey: z.string().optional(),
  period: z.string().default(localStorage.getItem(keys.CURRENT_PERIOD) || ""),
});

export const Route = createFileRoute("/income")({
  component: IncomeRoute,
  validateSearch: zodSearchValidator(incomeSearchSchema),
});

function IncomeRoute() {
  return <IncomeTable />;
}
