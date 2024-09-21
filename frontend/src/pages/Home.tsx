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
import { CategoryExpenseSummary, IncomeList, Period, User } from "../types";
import { Loading } from "./Loading.tsx";
import { Error } from "./Error.tsx";
import { useState } from "react";
import {
  useGetCategoryExpenseSummary,
  useGetIncome,
  useGetPeriod,
  useGetUser,
} from "./queries.ts";
import { utils } from "../utils";

export function Home() {
  const theme = useTheme();

  const [openNewExpense, setOpenNewExpense] = useState<boolean>(false);
  const [openNewIncome, setOpenNewIncome] = useState<boolean>(false);

  const lgUp: boolean = useMediaQuery(theme.breakpoints.up("lg"));

  const getUser = useGetUser();
  const getPeriod = useGetPeriod();
  const getCategoryExpenseSummary = useGetCategoryExpenseSummary();
  const getIncome = useGetIncome();

  const user: User | undefined = getUser.data?.data;
  const period: Period | undefined = getPeriod.data?.data;
  const categoryExpenseSummary: CategoryExpenseSummary[] =
    utils.setAdditionalData(getCategoryExpenseSummary.data?.data, user);
  const incomeList: IncomeList | undefined = getIncome.data?.data;
  const totalIncome: number = incomeList? incomeList.income.reduce((acc, cur) => acc+cur.amount, 0) : 0;

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
      getCategoryExpenseSummary.isRefetchError ||
      getPeriod.isRefetchError
    );
  }

  if (getUser.isPending && user === undefined) {
    return <Loading />;
  }

  if (getUser.isError && user === undefined) {
    return <Error />;
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
          <IncomeCard income={totalIncome} />
        </Grid>

        {/*Balance*/}
        <Grid xs={12} sm={6} hidden={lgUp}>
          <BalanceCard remainder={user ? user.remainder : 0} />
        </Grid>

        {/*Expenses*/}
        <Grid xs={12} sm={6} hidden={lgUp}>
          <ExpenseCard
            expenses={categoryExpenseSummary.reduce(
              (acc, curr) => acc + curr.total,
              0,
            )}
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
                  summary={categoryExpenseSummary ? categoryExpenseSummary : []}
                  chartHeight={chartHeight}
                  isLoading={getUser.isLoading}
                  isError={getCategoryExpenseSummary.isError}
                />
              </Grid>

              {/*New expense/income buttons, Current balance and expenses*/}
              <Grid xs={12} lg={4}>
                <div>
                  <Grid container mt={"1rem"} spacing={1}>
                    {/*Balance*/}
                    <Grid xs={12} hidden={!lgUp}>
                      <BalanceCard remainder={user ? user.remainder : 0} />
                    </Grid>

                    {/*Expenses*/}
                    <Grid xs={12} hidden={!lgUp}>
                      <ExpenseCard
                        expenses={categoryExpenseSummary.reduce(
                          (acc, curr) => acc + curr.total,
                          0,
                        )}
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
            <ExpensesTable categories={user.categories} />
          )}
        </Grid>
      </Grid>

      <NewTransaction
        type={"expense"}
        open={openNewExpense}
        onClose={handleNewExpenseClose}
      />
      <NewTransaction
        type={"income"}
        open={openNewIncome}
        onClose={handleNewIncomeClose}
      />
    </Container>
  );
}
