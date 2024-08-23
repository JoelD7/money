import { Typography, useMediaQuery, useTheme } from "@mui/material";
import Grid from "@mui/material/Unstable_Grid2";
import AddIcon from "@mui/icons-material/Add";
import {
  BackgroundRefetchErrorSnackbar,
  BalanceCard,
  Button,
  Container,
  ExpenseCard,
  ExpensesChart,
  ExpensesTable,
  LinearProgress,
  Navbar,
  NewExpense,
} from "../components";
import {
  CategoryExpenseSummary,
  Expense,
  Period,
  User,
} from "../types";
import { Loading } from "./Loading.tsx";
import { Error } from "./Error.tsx";
import { useState } from "react";
import {
  useGetCategoryExpenseSummary,
  useGetExpenses,
  useGetPeriod,
  useGetUser,
} from "./queries.ts";
import { utils } from "../utils";

export function Home() {
  const theme = useTheme();

  const [openNewExpense, setOpenNewExpense] = useState<boolean>(false);

  const mdUp: boolean = useMediaQuery(theme.breakpoints.up("md"));

  const getUser = useGetUser();
  const getExpenses = useGetExpenses();
  const getPeriod = useGetPeriod();
  const getCategoryExpenseSummary = useGetCategoryExpenseSummary();

  const user: User | undefined = getUser.data?.data;
  const expenses: Expense[] | undefined = getExpenses.data?.data.expenses;
  const period: Period | undefined = getPeriod.data?.data;
  const categoryExpenseSummary: CategoryExpenseSummary[] = utils.setAdditionalData(getCategoryExpenseSummary.data?.data, user);

  const chartHeight: number = 250;

  function handleClose() {
    setOpenNewExpense(false);
  }

  if (getUser.isPending && user === undefined) {
    return <Loading />;
  }

  if (getUser.isError && user === undefined) {
    return <Error />;
  }

  return (
    <Container>
      <BackgroundRefetchErrorSnackbar />
      <LinearProgress
        loading={
          getUser.isFetching || getExpenses.isFetching || getPeriod.isFetching
        }
      />
      <Navbar />

      <Grid container justifyContent={"center"} position={"relative"}>
        {/*Balance*/}
        <Grid xs={12} sm={6} hidden={mdUp}>
          <BalanceCard remainder={user ? user.remainder : 0} />
        </Grid>

        {/*Expenses*/}
        <Grid xs={12} sm={6} hidden={mdUp}>
          <ExpenseCard expenses={user ? user.expenses : 0} />
        </Grid>

        {/*Chart, Current balance and expenses*/}
        <Grid xs={12} maxWidth={"1200px"}>
          <div>
            <Grid container spacing={1}>
              {/*Chart section*/}
              <Grid xs={12} md={6}>
                <ExpensesChart
                  period={period}
                  summary={categoryExpenseSummary ? categoryExpenseSummary : []}
                  chartHeight={chartHeight}
                  isLoading={getUser.isLoading || getExpenses.isLoading}
                  isError={getExpenses.isError}
                />
              </Grid>
              {/*New expense/income buttons, Current balance and expenses*/}
              <Grid xs={12} md={6}>
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
                        onClick={() => setOpenNewExpense(true)}
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
        <Grid xs={12} maxWidth={"1200px"}>
          <Typography mt={"2rem"} variant={"h4"}>
            Latest
          </Typography>

          {expenses && user && user.categories && (
            <ExpensesTable expenses={expenses} categories={user.categories} />
          )}
        </Grid>
      </Grid>

      <NewExpense open={openNewExpense} onClose={handleClose} />
    </Container>
  );
}
