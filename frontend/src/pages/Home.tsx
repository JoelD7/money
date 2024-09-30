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
  IncomeCard,
  LinearProgress,
  Navbar,
  NewTransaction,
} from "../components";
import { Period, PeriodStats, User } from "../types";
import { Loading } from "./Loading.tsx";
import { Error } from "./Error.tsx";
import { useState } from "react";
import { useGetPeriod, useGetPeriodStats, useGetUser } from "./queries.ts";
import { utils } from "../utils";

export function Home() {
  const theme = useTheme();

  const [openNewExpense, setOpenNewExpense] = useState<boolean>(false);
  const [openNewIncome, setOpenNewIncome] = useState<boolean>(false);

  const lgUp: boolean = useMediaQuery(theme.breakpoints.up("lg"));

  const getUser = useGetUser();
  const user: User | undefined = getUser.data?.data;

  const getPeriod = useGetPeriod(user);
  const period: Period | undefined = getPeriod.data?.data;

  const getPeriodStats = useGetPeriodStats(user);
  const periodStats: PeriodStats | undefined = utils.setAdditionalData(
    getPeriodStats.data?.data,
    user,
  );

  const chartHeight: number = 350;

  function handleNewExpenseClose() {
    setOpenNewExpense(false);
  }

  function handleNewIncomeClose() {
    setOpenNewIncome(false);
  }

  function showRefetchErrorSnackbar() {
    return (
      getUser.isRefetchError ||
      getPeriodStats.isRefetchError ||
      getPeriod.isRefetchError
    );
  }

  if (getUser.isPending && user === undefined) {
    return <Loading />;
  }

  if (getUser.isError && user === undefined) {
    return <Error />;
  }

  function getPeriodTotalExpenses() {
    if (periodStats) {
      return periodStats.category_expense_summary.reduce(
        (acc, curr) => acc + curr.total,
        0,
      );
    }

    return 0;
  }

  return (
    <Container>
      <BackgroundRefetchErrorSnackbar show={showRefetchErrorSnackbar()} />
      <LinearProgress loading={getUser.isFetching || getPeriod.isFetching} />
      <Navbar />

      <Grid
        container
        justifyContent={"center"}
        position={"relative"}
        spacing={1}
        marginTop={"20px"}
      >
        {/*Income*/}
        <Grid xs={12} sm={6} hidden={lgUp}>
          <IncomeCard loading={true} income={0} />
        </Grid>

        {/*Balance*/}
        <Grid xs={12} sm={6} hidden={lgUp}>
          <BalanceCard
            loading={getUser.isPending}
            remainder={user ? user.remainder : 0}
          />
        </Grid>

        {/*Expenses*/}
        <Grid xs={12} sm={6} hidden={lgUp}>
          <ExpenseCard
            loading={getPeriodStats.isPending}
            expenses={getPeriodTotalExpenses()}
          />
        </Grid>

        {/*Chart, Current balance and expenses*/}
        <Grid xs={12} maxWidth={"1200px"}>
          <div>
            <Grid container spacing={1}>
              {/*Chart section*/}
              <Grid xs={12} lg={8}>
                <ExpensesChart
                  period={period}
                  summary={
                    periodStats ? periodStats.category_expense_summary : []
                  }
                  chartHeight={chartHeight}
                  isLoading={getUser.isLoading}
                  isError={getPeriodStats.isError}
                />
              </Grid>

              {/*New expense/income buttons, Current balance and expenses*/}
              <Grid xs={12} lg={4}>
                <div>
                  <Grid container mt={"1rem"} spacing={1}>
                    {/*Income*/}
                    <Grid xs={12} hidden={!lgUp}>
                      <IncomeCard
                        loading={getPeriodStats.isPending}
                        income={periodStats?.total_income}
                      />
                    </Grid>

                    {/*Balance*/}
                    <Grid xs={12} hidden={!lgUp}>
                      <BalanceCard
                        loading={getUser.isPending}
                        remainder={user ? user.remainder : 0}
                      />
                    </Grid>

                    {/*Expenses*/}
                    <Grid xs={12} hidden={!lgUp}>
                      <ExpenseCard
                        loading={getPeriodStats.isPending}
                        expenses={getPeriodTotalExpenses()}
                      />
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
                        onClick={() => setOpenNewIncome(true)}
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

          {user && user.categories && (
            <ExpensesTable
              period={user.current_period}
              categories={user.categories}
            />
          )}
        </Grid>
      </Grid>

      <NewTransaction
        type={"expense"}
        user={user}
        open={openNewExpense}
        onClose={handleNewExpenseClose}
      />
      <NewTransaction
        type={"income"}
        user={user}
        open={openNewIncome}
        onClose={handleNewIncomeClose}
      />
    </Container>
  );
}
