import { Cell, Pie, PieChart, ResponsiveContainer, Tooltip } from "recharts";
import {
  CategoryExpenseSummary,
  Period,
  RechartsLabelProps,
} from "../../types";
import { CircularProgress, Typography } from "@mui/material";
import Grid from "@mui/material/Unstable_Grid2";
import { Button } from "../atoms";
import { Colors } from "../../assets";
import { useNavigate } from "@tanstack/react-router";

type ExpensesChartProps = {
  summary: CategoryExpenseSummary[];
  period?: Period;
  chartHeight: number;
  isLoading: boolean;
  isError: boolean;
};

export function ExpensesChart({
  summary,
  chartHeight,
  period,
  isLoading,
  isError,
}: ExpensesChartProps) {
  const RADIAN: number = Math.PI / 180;

  const navigate = useNavigate();

  function getCustomLabel({
    cx,
    cy,
    midAngle,
    innerRadius,
    outerRadius,
    percent,
  }: RechartsLabelProps) {
    const radius = innerRadius + (outerRadius - innerRadius) * 0.5;
    const x = cx + radius * Math.cos(-midAngle * RADIAN);
    const y = cy + radius * Math.sin(-midAngle * RADIAN);

    return (
      <text
        x={x}
        y={y}
        fill="white"
        textAnchor={x > cx ? "start" : "end"}
        dominantBaseline="central"
      >
        {`${(percent * 100).toFixed(0)}%`}
      </text>
    );
  }

  function getPeriodDates(): string {
    if (!period) {
      return "";
    }

    return `${new Intl.DateTimeFormat("en-US", {
      month: "short",
      day: "2-digit",
    }).format(new Date(period.start_date))} - ${new Intl.DateTimeFormat(
      "default",
      {
        month: "short",
        day: "2-digit",
      },
    ).format(new Date(period.end_date))}`;
  }

  function getOpacity(): number {
    return isLoading || isError ? 0 : 1;
  }

  function viewPeriodDetails() {
    const route = `/period/${period ? period.name : "current"}`;

    navigate({ to: route })
      .then(() => {})
      .catch((err) => {
        console.error("Error navigating to /login", err);
      });
  }

  return (
    <div>
      <Grid
        container
        bgcolor={"white.main"}
        borderRadius="1rem"
        p="1rem"
        boxShadow="3"
        mt="1rem"
        style={{ position: "relative" }}
      >
        {isError && (
          <div className="absolute top-0 left-0 right-0 bottom-0 flex items-center justify-center z-10 rounded-xl">
            <ChartError />
          </div>
        )}

        {isLoading && (
          <div className="absolute top-0 left-0 right-0 bottom-0 flex items-center justify-center z-10 rounded-xl">
            <CircularProgress size={"7rem"} />
          </div>
        )}

        <Grid xs={12} style={{ opacity: getOpacity() }}>
          <Typography variant="h4">{period ? period.name : ""}</Typography>
          <Typography color="gray.light">{getPeriodDates()}</Typography>
        </Grid>

        <Grid xs={12} height={chartHeight} style={{ opacity: getOpacity() }}>
          <ResponsiveContainer width="100%" height="100%">
            <PieChart>
              <Pie
                data={summary}
                label={getCustomLabel}
                dataKey="total"
                nameKey="name"
                cx="50%"
                cy="50%"
                labelLine={false}
                fill="#8884d8"
              >
                {summary.map((category, index) => (
                  <Cell key={`cell-${index}`} fill={category.color} />
                ))}
              </Pie>
              <Tooltip
                formatter={(value, name) => [
                  value.toLocaleString("en-US", {
                    minimumFractionDigits: 2,
                    maximumFractionDigits: 2,
                  }),
                  name,
                ]}
              />
            </PieChart>
          </ResponsiveContainer>
        </Grid>

        <Grid xs={12} style={{ opacity: getOpacity() }}>
          <Grid container width="100%" className="justify-between">
            <Grid xs={6}>
              {summary.map((ce) => (
                <div
                  key={`${ce.category_id}`}
                  className="flex gap-1 items-center"
                >
                  <div
                    className="rounded-full w-3 h-3"
                    style={{ backgroundColor: ce.color }}
                  />
                  <Typography color="gray.light">{ce.name}</Typography>
                </div>
              ))}
            </Grid>
            <Grid xs={6}>
              <Grid container className="items-end h-full">
                <Button
                  variant="outlined"
                  sx={{
                    textTransform: "capitalize",
                    borderRadius: "1rem",
                    height: "fit-content",
                  }}
                  onClick={() => viewPeriodDetails()}
                >
                  View details
                </Button>
              </Grid>
            </Grid>
          </Grid>
        </Grid>
      </Grid>
    </div>
  );
}

function ChartError() {
  return (
    <div className={"p-4"}>
      <Typography color={"darkGreen.main"} variant={"h4"}>
        Whoops...
      </Typography>
      <Typography variant={"h5"} color={"gray.darker"}>
        Couldn't load chart...
      </Typography>
      <div className={"pt-4 md:max-w-sm"}>
        <p style={{ color: Colors.GRAY_DARK, fontSize: "18px" }}>
          Our servers seem to be having some issues. Please try again in a few
          minutes.
        </p>
      </div>

      <Button
        variant={"contained"}
        sx={{
          marginTop: "10px",
          fontSize: "16px",
        }}
        onClick={() => window.location.reload()}
      >
        Reload
      </Button>
    </div>
  );
}
