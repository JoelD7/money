import { Cell, Pie, PieChart, ResponsiveContainer, Tooltip } from "recharts";
import {
  CategoryExpenseSummary,
  Period,
  PeriodStats,
  RechartsLabelProps,
  User,
} from "../../types";
import { CircularProgress, Typography } from "@mui/material";
import Grid from "@mui/material/Unstable_Grid2";
import { Button } from "../atoms";
import { Colors } from "../../assets";
import { useGetPeriod, useGetPeriodStats } from "../../queries";
import { utils } from "../../utils";
import { ReactNode } from "react";

type ExpensesChartProps = {
  user?: User;
  chartHeight: number;
};

export function ExpensesChart({ user, chartHeight }: ExpensesChartProps) {
  const RADIAN: number = Math.PI / 180;

  const getPeriod = useGetPeriod(user);
  const period: Period | undefined = getPeriod.data;
  const getPeriodStats = useGetPeriodStats(user);
  const periodStats: PeriodStats | undefined = utils.setAdditionalData(
    getPeriodStats.data,
    user,
  );

  const totalExpenses: number = periodStats
    ? periodStats.category_expense_summary.reduce((acc, ce) => acc + ce.total, 0)
    : 0;
  const summary: CategoryExpenseSummary[] = periodStats
    ? periodStats.category_expense_summary
    : [];

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
    }).format(new Date(period.start_date))} - ${new Intl.DateTimeFormat("default", {
      month: "short",
      day: "2-digit",
    }).format(new Date(period.end_date))}`;
  }

  if (getPeriodStats.isLoading || getPeriod.isLoading) {
    return (
      <ExpensesChartContainer>
        <div className="absolute top-0 left-0 right-0 bottom-0 flex items-center justify-center z-10 rounded-xl">
          <CircularProgress size={"7rem"} />
        </div>
      </ExpensesChartContainer>
    );
  }

  if (getPeriodStats.isError && getPeriodStats.error.response?.status === 404) {
    return (
      <ExpensesChartContainer>
        <div className="flex items-center justify-center z-10 rounded-xl">
          <Chart404Error />
        </div>
      </ExpensesChartContainer>
    );
  }

  if (getPeriodStats.isError || getPeriod.isError) {
    return (
      <ExpensesChartContainer>
        <div className="absolute top-0 left-0 right-0 bottom-0 flex items-center justify-center z-10 rounded-xl">
          <ChartError />
        </div>
      </ExpensesChartContainer>
    );
  }

  return (
    <ExpensesChartContainer>
      <Grid xs={12}>
        <Typography variant="h4">{period ? period.name : ""}</Typography>
        <Typography color="gray.light">{getPeriodDates()}</Typography>
      </Grid>

      {/*Chart and category summary*/}
      <Grid xs={12} height={chartHeight}>
        <div className={"h-full"}>
          <Grid container height={"100%"}>
            {/*Chart*/}
            <Grid xs={7}>
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

            {/*Summary*/}
            <Grid xs={5}>
              <div className={"h-full"}>
                <Grid container height={"100%"} alignItems={"center"}>
                  {/*Category data*/}
                  <Grid xs={8} marginBottom={"50px"}>
                    {summary.map((category) => (
                      <div key={category.category_id} className="flex gap-1 items-center">
                        <div
                          className="rounded-full w-3 h-3"
                          style={{ backgroundColor: category.color }}
                        />
                        <Typography sx={{ color: category.color }}>
                          {Math.round((category.total * 100) / totalExpenses)}%
                        </Typography>
                        <Typography color="gray.light">{category.name}</Typography>
                      </div>
                    ))}
                  </Grid>

                  {/*Category total*/}
                  <Grid xs={4} marginBottom={"50px"}>
                    {summary.map((category) => (
                      <Typography key={category.category_id} color="gray.light">
                        {new Intl.NumberFormat("en-US", {
                          style: "currency",
                          currency: "USD",
                        }).format(category?.total)}
                      </Typography>
                    ))}
                  </Grid>
                </Grid>
              </div>
            </Grid>
          </Grid>
        </div>
      </Grid>
    </ExpensesChartContainer>
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
          Our servers seem to be having some issues. Please try again in a few minutes.
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

function Chart404Error() {
  return (
    <div className={"p-4"}>
      <Typography color={"darkGreen.main"} variant={"h4"}>
        Your expenses will show here
      </Typography>
      <div className={"pt-4 md:max-w-sm"}>
        <p style={{ color: Colors.GRAY_DARK, fontSize: "18px" }}>
          Add expenses so you can see them in a chart
        </p>
      </div>
    </div>
  );
}

type ExpensesChartContainerProps = {
  children?: ReactNode;
};

function ExpensesChartContainer({ children }: ExpensesChartContainerProps) {
  return (
    <div>
      <Grid
        container
        bgcolor={"white.main"}
        borderRadius="1rem"
        p="1rem"
        boxShadow="3"
        mt="1rem"
        style={{ position: "relative", minHeight: "450px" }}
      >
        {children}
      </Grid>
    </div>
  );
}
