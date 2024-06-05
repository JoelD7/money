import { Typography, useMediaQuery, useTheme } from "@mui/material";
import Grid from "@mui/material/Unstable_Grid2";
import AddIcon from "@mui/icons-material/Add";
import {
  BalanceCard,
  Button,
  ExpenseCard,
  ExpensesTable,
  Navbar,
} from "../components";
import { Cell, Pie, PieChart, ResponsiveContainer, Tooltip } from "recharts";
import { Expense, Period, RechartsLabelProps, User } from "../types";
import json2mq from "json2mq";
import { useQuery } from "@tanstack/react-query";
import api from "../api";
import { Loading } from "./Loading.tsx";
import { Error } from "./Error.tsx";

export function Home() {
  const theme = useTheme();

  const mdUp: boolean = useMediaQuery(theme.breakpoints.up("md"));

  const getUser = useQuery({
    queryKey: ["user"],
    queryFn: () => api.getUser(),
    refetchOnWindowFocus: false,
  });

  const getExpenses = useQuery({
    queryKey: ["expenses"],
    queryFn: () => api.getExpenses(),
    refetchOnWindowFocus: false,
  });

  const getPeriod = useQuery({
    queryKey: ["period"],
    queryFn: () => api.getPeriod(),
    refetchOnWindowFocus: false,
  });

  const user: User | undefined = getUser.data?.data;
  const expenses: Expense[] | undefined = getExpenses.data?.data.expenses;
  const period: Period | undefined = getPeriod.data?.data;

  const xlCustom = useMediaQuery(
    json2mq({
      maxWidth: 2300,
    }),
  );

  const chartHeight: number = 250;
  const RADIAN: number = Math.PI / 180;

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

  if (getUser.isPending) {
    return <Loading />;
  }

  if (getUser.isError) {
    return <Error />;
  }

  return (
    <>
      <Navbar />
      <Grid container spacing={1} justifyContent={"center"}>
        {/*Balance*/}
        <Grid xs={12} sm={6} hidden={mdUp}>
          <BalanceCard remainder={user ? user.remainder : 0} />
        </Grid>

        {/*Expenses*/}
        <Grid xs={12} sm={6} hidden={mdUp}>
          <ExpenseCard expenses={user ? user.expenses : 0} />
        </Grid>

        {/*Chart, Current balance and expenses*/}
        <Grid xs={12} maxWidth={"880px"}>
          <div>
            <Grid container spacing={1}>
              {/*Chart section*/}
              <Grid xs={12} md={6} maxWidth={"430px"}>
                <div>
                  <Grid
                    container
                    bgcolor={"white.main"}
                    borderRadius="1rem"
                    p="1rem"
                    boxShadow="3"
                    mt="1rem"
                  >
                    <Grid xs={12}>
                      {period && (
                        <>
                          <Typography variant="h4">{period.name}</Typography>
                          <Typography color="gray.light">
                            {getPeriodDates()}
                          </Typography>
                        </>
                      )}
                    </Grid>
                    {/*Chart*/}
                    <Grid xs={12} height={chartHeight}>
                      <ResponsiveContainer width="100%" height="100%">
                        <PieChart width={350} height={chartHeight}>
                          <Pie
                            data={user ? user.categories : []}
                            label={getCustomLabel}
                            dataKey="value"
                            nameKey="name"
                            cx="50%"
                            cy="50%"
                            labelLine={false}
                            fill="#8884d8"
                          >
                            {user &&
                              user.categories &&
                              user.categories.map((category, index) => (
                                <Cell
                                  key={`cell-${index}`}
                                  fill={category.color}
                                />
                              ))}
                          </Pie>
                          <Tooltip />
                        </PieChart>
                      </ResponsiveContainer>
                    </Grid>

                    {/*Chart legend*/}
                    <Grid xs={12}>
                      <Grid container width="100%" className="justify-between">
                        {/*Categories*/}
                        <Grid xs={6}>
                          {user &&
                            user.categories &&
                            user.categories.map((category) => (
                              <div
                                key={`${category.id}`}
                                className="flex gap-1 items-center"
                              >
                                <div
                                  className="rounded-full w-3 h-3"
                                  style={{ backgroundColor: category.color }}
                                />
                                <Typography color="gray.light">
                                  {category.name}
                                </Typography>
                              </div>
                            ))}
                        </Grid>
                        {/*Details button*/}
                        <Grid xs={6}>
                          <Grid container className="items-end h-full">
                            <Button
                              variant="outlined"
                              sx={{
                                textTransform: "capitalize",
                                borderRadius: "1rem",
                                height: "fit-content",
                              }}
                            >
                              View details
                            </Button>
                          </Grid>
                        </Grid>
                      </Grid>
                    </Grid>
                  </Grid>
                </div>
              </Grid>

              {/*New expense/income buttons, Current balance and expenses*/}
              <Grid xs={12} md={6} maxWidth={"430px"}>
                <div>
                  <Grid container mt={"1rem"} spacing={1}>
                    {/*Balance*/}
                    <Grid xs={12} hidden={!mdUp}>
                      <BalanceCard remainder={user ? user.remainder : 0} />
                    </Grid>

                    {/*Expenses*/}
                    <Grid xs={12} hidden={!mdUp}>
                      <ExpenseCard expenses={user ? user.expenses : 0} />
                    </Grid>

                    {/**New expense/income buttons*/}
                    <Grid xs={12}>
                      <Button
                        color={"secondary"}
                        variant={"contained"}
                        startIcon={<AddIcon />}
                      >
                        New expense
                      </Button>

                      <Button
                        sx={{ marginLeft: "1rem" }}
                        variant={"contained"}
                        startIcon={<AddIcon />}
                      >
                        New income
                      </Button>
                    </Grid>
                  </Grid>
                </div>
              </Grid>
            </Grid>
          </div>
        </Grid>

        {/*Latest table*/}
        <Grid xs={12} maxWidth={xlCustom ? "1200px" : "none"}>
          <Typography mt={"2rem"} variant={"h4"}>
            Latest
          </Typography>

          {expenses && <ExpensesTable expenses={expenses} />}
        </Grid>
      </Grid>
    </>
  );
}
