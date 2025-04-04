import { createFileRoute } from "@tanstack/react-router";
import { SavingGoalDetail } from "../pages";

export const Route = createFileRoute("/savings/goals/$savingGoalId")({
  component: SavingGoalDetailRoute,
});

function SavingGoalDetailRoute() {
  return <SavingGoalDetail />;
}
