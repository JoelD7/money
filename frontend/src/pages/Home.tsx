import { Typography, useMediaQuery, useTheme } from "@mui/material";
import Grid from "@mui/material/Unstable_Grid2";
import AddIcon from "@mui/icons-material/Add";
import {
  BackgroundRefetchErrorSnackbar,
  BalanceCard,
  Button,
  Container,
  ErrorSnackbar,
  ExpenseCard,
  ExpensesChart,
  ExpensesTable,
  IncomeCard,
  LinearProgress,
  Navbar,
  NewTransaction,
} from "../components";
import { PeriodStats, User } from "../types";
import { Loading } from "./Loading.tsx";
import { Error } from "./Error.tsx";
import { useState } from "react";
import { useGetPeriodStats, useGetUser } from "../queries";
import { utils } from "../utils";

export function Home() {
  const theme = useTheme();

  const [openNewExpense, setOpenNewExpense] = useState<boolean>(false);
  const [openNewIncome, setOpenNewIncome] = useState<boolean>(false);
  const errSnackbar = {
    open: true,
    title: "Error fetching period stats",
  };

  const lgUp: boolean = useMediaQuery(theme.breakpoints.up("lg"));

  const getUser = useGetUser();
  const user: User | undefined = getUser.data;
  const getPeriodStats = useGetPeriodStats(user);
  const periodStats: PeriodStats | undefined = utils.setAdditionalData(
    getPeriodStats.data,
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
    return getUser.isRefetchError || getPeriodStats.isRefetchError;
  }

  if (getUser.isLoading) {
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

  function showPeriodStatsErr(): boolean{
    if (getPeriodStats.isError && getPeriodStats.error.response){
      return getPeriodStats.error.response.status !== 404
    }

    return getPeriodStats.isError
  }

  return (
    <Container>
      <Navbar />

      <BackgroundRefetchErrorSnackbar show={showRefetchErrorSnackbar()} />
      <LinearProgress loading={getUser.isFetching} />

      {showPeriodStatsErr() && (
        <ErrorSnackbar openProp={errSnackbar.open} title={errSnackbar.title} />
      )}

      <Grid
        container
        justifyContent={"center"}
        position={"relative"}
        spacing={1}
        marginTop={"20px"}
      >
        {/*Income*/}
        <Grid xs={12} sm={4} hidden={lgUp}>
          <IncomeCard
            loading={getPeriodStats.isLoading}
            income={periodStats?.total_income}
          />
        </Grid>

        {/*Balance*/}
        <Grid xs={12} sm={4} hidden={lgUp}>
          <BalanceCard
            loading={getUser.isLoading}
            remainder={user ? user.remainder : 0}
          />
        </Grid>

        {/*Expenses*/}
        <Grid xs={12} sm={4} hidden={lgUp}>
          <ExpenseCard
            loading={getPeriodStats.isLoading}
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
                  user={user}
                  chartHeight={chartHeight}
                />
              </Grid>

              {/*New expense/income buttons, Current balance and expenses*/}
              <Grid xs={12} lg={4}>
                <div>
                  <Grid container mt={"1rem"} spacing={1}>
                    {/*Income*/}
                    <Grid xs={12} hidden={!lgUp}>
                      <IncomeCard
                        loading={getPeriodStats.isLoading}
                        income={periodStats?.total_income}
                      />
                    </Grid>

                    {/*Balance*/}
                    <Grid xs={12} hidden={!lgUp}>
                      <BalanceCard
                        loading={getUser.isLoading}
                        remainder={user ? user.remainder : 0}
                      />
                    </Grid>

                    {/*Expenses*/}
                    <Grid xs={12} hidden={!lgUp}>
                      <ExpenseCard
                        loading={getPeriodStats.isLoading}
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
            <ExpensesTable period={user.current_period} categories={user.categories} />
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
