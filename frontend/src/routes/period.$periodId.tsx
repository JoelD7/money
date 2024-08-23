import { createFileRoute } from "@tanstack/react-router";
import { PeriodDetail } from "../pages";

export const Route = createFileRoute("/period/$periodId")({
  component: PeriodDetailRoute,
});

function PeriodDetailRoute() {
  return <PeriodDetail />;
}
