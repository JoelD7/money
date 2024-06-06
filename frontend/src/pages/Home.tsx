import { Typography, useMediaQuery, useTheme } from "@mui/material";
import Grid from "@mui/material/Unstable_Grid2";
import AddIcon from "@mui/icons-material/Add";
import {
  BalanceCard,
  Button,
  ExpenseCard,
  ExpensesChart,
  ExpensesTable,
  Navbar,
} from "../components";
import { Expense, User } from "../types";
import json2mq from "json2mq";
import { useQuery } from "@tanstack/react-query";
import api from "../api";
import { Loading } from "./Loading.tsx";
import { Error } from "./Error.tsx";
import { Colors } from "../assets";

type CategoryExpense = {
  category: string;
  color: string;
  value: number;
};

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

  const user: User | undefined = getUser.data?.data;
  const expenses: Expense[] | undefined = getExpenses.data?.data.expenses;
  const colorsByCategory: Map<string, string> = getColorsByCategory();
  const categoryExpense: CategoryExpense[] = getCategoryExpense();

  const xlCustom = useMediaQuery(
    json2mq({
      maxWidth: 2300,
    }),
  );

  const chartHeight: number = 250;

  function getColorsByCategory(): Map<string, string> {
    const colorsByCategory: Map<string, string> = new Map<string, string>();
    if (user && user.categories) {
      user.categories.forEach((category) => {
        colorsByCategory.set(category.name, category.color);
      });
    }

    return colorsByCategory;
  }

  function getCategoryExpense(): CategoryExpense[] {
    if (!user || !expenses) {
      return [];
    }

    const categoryExpense: CategoryExpense[] = [];
    const totalExpenseByCategory: Map<string, number> = new Map<
      string,
      number
    >();

    expenses.forEach((expense) => {
      let category = expense.category_name;
      if (!category) {
        category = "Other";
      }

      const curTotal: number | undefined = totalExpenseByCategory.get(category);

      if (curTotal) {
        totalExpenseByCategory.set(category, curTotal + expense.amount);
      } else {
        totalExpenseByCategory.set(category, expense.amount);
      }
    });

    totalExpenseByCategory.forEach((value, key) => {
      let color = colorsByCategory.get(key);
      if (!color) {
        color = Colors.GRAY_DARK;
      }

      categoryExpense.push({ category: key, value: value, color: color });
    });

    return categoryExpense;
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
                <ExpensesChart
                  categoryExpense={categoryExpense}
                  chartHeight={chartHeight}
                  isLoading={getUser.isLoading || getExpenses.isLoading}
                  isError={getExpenses.isError}
                />
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
