import { Typography, useMediaQuery, useTheme } from "@mui/material";
import Grid from "@mui/material/Unstable_Grid2";
import AddIcon from "@mui/icons-material/Add";
import { BalanceCard, Button, ExpenseCard, ExpensesTable } from "../components";
import { Cell, Pie, PieChart, ResponsiveContainer, Tooltip } from "recharts";
import { Expense, RechartsLabelProps, User } from "../types";
import json2mq from "json2mq";
import { useQuery } from "@tanstack/react-query";
import { api } from "../api";
import { AxiosError } from "axios";

export function Home() {
  const theme = useTheme();

  const mdUp: boolean = useMediaQuery(theme.breakpoints.up("md"));

  const getUserQuery = useQuery({
    queryKey: ["user"],
    queryFn: () => api.getUser(),
    retry: (_, error) => {
      const err = error as AxiosError;
      let shouldRetry: boolean = false;

      if (err.response?.status === 401) {
        shouldRetry = handleQueryRetry();
      }

      return shouldRetry;
    },
    refetchOnWindowFocus: false,
  });

  function handleQueryRetry(): boolean {
    let shouldRetry: boolean = false;

    api
        .refreshToken()
        .then(() => {
          shouldRetry = true;
          return;
        })
        .catch((error) => {
          console.error("Error refreshing token", error);

          const myErr = error as AxiosError;
          if (!myErr.response?.status) {
            //TODO: logout
            shouldRetry = false;
            return;
          }

          if (
              myErr.response?.status >= 400 &&
              myErr.response?.status <= 500
          ) {
            //TODO: logout
            shouldRetry = false;
            return;
          }

          shouldRetry = false;
          return;
        });

    return shouldRetry;
  }

  const user: User | undefined = getUserQuery.data?.data;

  const period = {
    username: "test@gmail.com",
    period: "asdf",
    name: "December",
    start_date: "2023-11-26T00:00:00Z",
    end_date: "2023-12-24T00:00:00Z",
  };

  const expenses: Expense[] = [
    {
      expenseID: "EX5DK8d8LTywTKC8r87vdS",
      username: "test@gmail.com",
      categoryID: "CTGiBScOP3V16LYBjdIStP9",
      categoryName: "Shopping",
      amount: 12.99,
      name: "Blue pair of socks",
      notes:
        "Ipsum mollit est pariatur esse ex. Aliqua laborum laboris laboris laboris. Laboris pectum",
      createdDate: new Date("2023-10-27T23:42:54.980596532Z"),
      period: "2023-7",
      updateDate: new Date("0001-01-01T00:00:00Z"),
    },
    {
      expenseID: "EXBLsynfE2QSAX8awfWptn",
      username: "test@gmail.com",
      categoryID: "CTGcSuhjzVmu3WrHLKD5fhS",
      categoryName: "Health",
      amount: 1000,
      name: "Protector solar",
      createdDate: new Date("2023-10-14T19:55:45.261990038Z"),
      period: "2023-7",
      updateDate: new Date("0001-01-01T00:00:00Z"),
    },
    {
      expenseID: "EXD5G8OdwlKC81tH9ZE3eO",
      username: "test@gmail.com",
      categoryID: "CTGiBScOP3V16LYBjdIStP9",
      categoryName: "Shopping",
      amount: 1898.11,
      name: "Vacuum Cleaner",
      createdDate: new Date("2023-10-18T22:41:56.024322091Z"),
      period: "2023-7",
      updateDate: new Date("0001-01-01T00:00:00Z"),
    },
    {
      expenseID: "EXF5Mg3fpxct3v0BI91XYB",
      username: "test@gmail.com",
      categoryID: "CTGiBScOP3V16LYBjdIStP9",
      categoryName: "Shopping",
      amount: 1202.17,
      name: "Microwave",
      createdDate: new Date("2023-10-18T22:41:46.946640398Z"),
      period: "2023-7",
      updateDate: new Date("0001-01-01T00:00:00Z"),
    },
    {
      expenseID: "EXHrwzQezXK6nXyclUHbVH",
      username: "test@gmail.com",
      categoryID: "CTGGyouAaIPPWKzxpyxHACS",
      categoryName: "Entertainment",
      amount: 955,
      name: "Plza Juan Baron",
      notes: "Lorem ipsum note to fill space",
      createdDate: new Date("2023-10-14T19:52:11.552327532Z"),
      period: "2023-7",
      updateDate: new Date("0001-01-01T00:00:00Z"),
    },
    {
      expenseID: "EXIGBwc0sBWeyL9hy8jVuI",
      username: "test@gmail.com",
      categoryID: "CTGiBScOP3V16LYBjdIStP9",
      categoryName: "Shopping",
      amount: 620,
      name: "Correa amarilla",
      createdDate: new Date("2023-10-18T22:37:04.230522146Z"),
      period: "2023-7",
      updateDate: new Date("0001-01-01T00:00:00Z"),
    },
    {
      expenseID: "EXIfxidmlBJtq97xjQZfNh",
      username: "test@gmail.com",
      categoryID: "CTGiBScOP3V16LYBjdIStP9",
      categoryName: "Shopping",
      amount: 123,
      name: "Correa azul",
      createdDate: new Date("2023-10-18T22:37:15.57296052Z"),
      period: "2023-7",
      updateDate: new Date("0001-01-01T00:00:00Z"),
    },
    {
      expenseID: "EXP123",
      username: "test@gmail.com",
      amount: 893,
      name: "Jordan shopping",
      createdDate: new Date("0001-01-01T00:00:00Z"),
      period: "2023-5",
      updateDate: new Date("0001-01-01T00:00:00Z"),
    },
    {
      expenseID: "EXP456",
      username: "test@gmail.com",
      amount: 112,
      name: "Uber drive",
      createdDate: new Date("0001-01-01T00:00:00Z"),
      period: "2023-5",
      updateDate: new Date("0001-01-01T00:00:00Z"),
    },
    {
      expenseID: "EXP789",
      username: "test@gmail.com",
      amount: 525,
      name: "Lunch",
      createdDate: new Date("0001-01-01T00:00:00Z"),
      period: "2023-5",
      updateDate: new Date("0001-01-01T00:00:00Z"),
    },
  ];

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

  return (
    <>
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
                      <Typography variant="h4">{period.name}</Typography>
                      <Typography color="gray.light">
                        {getPeriodDates()}
                      </Typography>
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
                                key={category.categoryID}
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

          <ExpensesTable expenses={expenses} />
        </Grid>
      </Grid>
    </>
  );
}
