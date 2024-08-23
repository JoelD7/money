import {Box, Typography, useMediaQuery, useTheme} from "@mui/material";
import {CategoryExpenseSummary, Expense, Period, RechartsLabelProps, User} from "../types";
import Grid from "@mui/material/Unstable_Grid2";
import ArrowCircleUpRoundedIcon from "@mui/icons-material/ArrowCircleUpRounded";
import ArrowCircleDownRoundedIcon from "@mui/icons-material/ArrowCircleDownRounded";
import {Cell, Pie, PieChart, ResponsiveContainer, Tooltip} from "recharts";
import {
    BackgroundRefetchErrorSnackbar,
    BalanceCard,
    Button,
    Container,
    ExpenseCard,
    ExpensesTable, LinearProgress, Navbar
} from "../components";
import AddIcon from "@mui/icons-material/Add";
import ArrowCircleLeftIcon from '@mui/icons-material/ArrowCircleLeft';
import {Colors} from "../assets";
import {useGetCategoryExpenseSummary, useGetExpenses, useGetPeriod, useGetUser} from "./queries.ts";
import { utils } from "../utils";
import {Loading} from "./Loading.tsx";
import {Error} from "./Error.tsx";

export function PeriodDetail() {
    const cashFlowIconStyles = {
        '&.MuiSvgIcon-root': {
            width: "28px",
            height: "28px",
        },
    }

    const backButtonStyle = {
        '&.MuiSvgIcon-root': {
            width: "38px",
            height: "38px",
        },
    }

    const theme = useTheme()
    const smUp: boolean = useMediaQuery(theme.breakpoints.up('sm'));
    const mdUp: boolean = useMediaQuery(theme.breakpoints.up('md'));

    const getUser = useGetUser()
    const getExpenses = useGetExpenses()
    const getPeriod = useGetPeriod()
    const getCategoryExpenseSummary = useGetCategoryExpenseSummary();

    const user: User | undefined = getUser.data?.data;
    const expenses: Expense[] | undefined = getExpenses.data?.data.expenses;
    const period: Period | undefined = getPeriod.data?.data;
    const categoryExpenseSummary: CategoryExpenseSummary[] = utils.setAdditionalData(getCategoryExpenseSummary.data?.data, user);

    const totalExpenses: number = categoryExpenseSummary.reduce((acc, category) => acc + category.total, 0)

    const colorByCategory: Map<string, string> = getColorByCategory()

    function getColorByCategory(): Map<string, string> {
        const m: Map<string, string> = new Map<string, string>()

        if (!user || !user.categories) {
            return m
        }

        m.set("Other", Colors.GRAY_DARK)

        user.categories.forEach((category) => {
            m.set(category.id, category.color)
        })

        return m
    }

    const RADIAN: number = Math.PI / 180;

    function getCustomLabel({cx, cy, midAngle, innerRadius, outerRadius, percent}: RechartsLabelProps) {
        const radius = innerRadius + (outerRadius - innerRadius) * 0.5;
        const x = cx + radius * Math.cos(-midAngle * RADIAN);
        const y = cy + radius * Math.sin(-midAngle * RADIAN);

        return (
            <text x={x} y={y} fill="white" textAnchor={x > cx ? 'start' : 'end'} dominantBaseline="central">
                {`${(percent * 100).toFixed(0)}%`}
            </text>
        );
    }

    function getChartHeight(): number {
        if (smUp) {
            return 350
        }

        return 250
    }

    function getChartWidth(): number {
        if (smUp) {
            return 450
        }

        return 350
    }

    function getPeriodDates(): string {
        if (!period) {
            return ""
        }

        return `${new Intl.DateTimeFormat('en-US', {
            month: 'short',
            day: '2-digit',
        }).format(new Date(period.start_date))} - ${new Intl.DateTimeFormat('default', {
            month: 'short',
            day: '2-digit',
        }).format(new Date(period.end_date))}`
    }

    function getCategoryExpensePercentage(total: number): number {
        if (totalExpenses === 0) {
            return 0
        }

        return (total / totalExpenses) * 100
    }

    if (getUser.isPending && user === undefined) {
        return <Loading />;
    }

    if (getUser.isError && user === undefined) {
        return <Error />;
    }

    return (
        <Container>
            <BackgroundRefetchErrorSnackbar/>
            <LinearProgress loading={getUser.isFetching || getExpenses.isFetching || getPeriod.isFetching}/>
            <Navbar/>

            <Grid container spacing={1} justifyContent={"center"}>
                {/*Page title*/}
                <Grid xs={12} mt={"1rem"}>
                    <div>
                        <Grid container>
                            <Grid xs={2}>
                                <div className={"flex justify-center items-center h-full"}>
                                    {/*@ts-ignore*/}
                                    <ArrowCircleLeftIcon sx={backButtonStyle} color="darkGreen"/>
                                </div>
                            </Grid>
                            <Grid xs={10}>
                                <div>
                                    <Grid container>
                                        <Grid xs={12}>
                                            <Typography variant={"h3"}>
                                                {period && period.name}
                                            </Typography>
                                        </Grid>
                                        <Grid xs={12}>
                                            <Typography color="gray.light">
                                                {getPeriodDates()}
                                            </Typography>
                                        </Grid>
                                    </Grid>
                                </div>
                            </Grid>
                        </Grid>
                    </div>
                </Grid>

                {/*Balance and expenses*/}
                <Grid xs={12} hidden={mdUp}>
                    <Box hidden={smUp} maxWidth={"435px"} borderRadius="1rem" mt={"1rem"} p="1rem" bgcolor="white.main"
                         boxShadow={"2"}>
                        <div className={"flex w-11/12 m-auto items-center"}>
                            {/*@ts-ignore*/}
                            <ArrowCircleUpRoundedIcon sx={cashFlowIconStyles} color="darkGreen"/>
                            <Typography variant={"h6"} fontWeight={"bold"}>
                                You have:
                            </Typography>

                            <Typography variant={"h6"} color={"primary.darker"} marginLeft={"auto"}>
                                {new Intl.NumberFormat('en-US', {
                                    style: 'currency',
                                    currency: 'USD'
                                }).format(user ? user.remainder : 0)}
                            </Typography>
                        </div>
                        <div className={"flex w-11/12 m-auto items-center"}>
                            {/*@ts-ignore*/}
                            <ArrowCircleDownRoundedIcon sx={cashFlowIconStyles} color="red"/>
                            <Typography variant={"h6"} fontWeight={"bold"}>
                                You have spent:
                            </Typography>

                            <Typography variant={"h6"} color={"red.main"} marginLeft={"auto"}>
                                {new Intl.NumberFormat('en-US', {
                                    style: 'currency',
                                    currency: 'USD'
                                }).format(user && user.expenses ? user.expenses : 0)}
                            </Typography>
                        </div>
                    </Box>
                    <div hidden={!smUp}>
                        <Grid container spacing={1}>
                            <Grid xs={6}>
                                <BalanceCard remainder={user ? user.remainder : 0}/>
                            </Grid>
                            <Grid xs={6}>
                                <ExpenseCard expenses={user && user.expenses ? user.expenses : 0}/>
                            </Grid>
                        </Grid>
                    </div>
                </Grid>

                {/**New expense/income buttons*/}
                <Grid xs={12} hidden={mdUp}>
                    <div className={"flex mt-3"}>
                        <Button color={"secondary"} variant={"contained"}
                                startIcon={<AddIcon/>}>
                            New expense
                        </Button>

                        <Button sx={{marginLeft: "1rem"}} variant={"contained"}
                                startIcon={<AddIcon/>}>
                            New income
                        </Button>
                    </div>
                </Grid>

                {/*Chart, Current balance and expenses*/}
                <Grid xs={12} maxWidth={"1200px"}>
                    <div className={"mt-4"}>
                        <Typography variant={"h4"}>Breakdown</Typography>
                        <Grid container spacing={1} mt="1rem">
                            {/*Chart section*/}
                            <Grid xs={12} md={6} lg={7}>
                                <div>
                                    <Grid container bgcolor={"white.main"} borderRadius="1rem"
                                          p="1rem" boxShadow="3" alignItems={"center"}>
                                        {/*Chart*/}
                                        <Grid xs={12} sm={5} md={12} lg={5} height={getChartHeight()}>
                                            <ResponsiveContainer width="100%" height="100%">
                                                <PieChart width={getChartWidth()} height={getChartHeight()}>
                                                    <Pie data={categoryExpenseSummary}
                                                         label={getCustomLabel}
                                                         dataKey="total"
                                                         nameKey="name"
                                                         cx="50%"
                                                         cy="50%"
                                                         labelLine={false}>
                                                        {categoryExpenseSummary.map((category, index) => (
                                                            <Cell key={`cell-${index}`}
                                                                  fill={category.color}/>
                                                        ))}
                                                    </Pie>
                                                    <Tooltip/>
                                                </PieChart>
                                            </ResponsiveContainer>
                                        </Grid>

                                        {/*Total by category*/}
                                        <Grid xs={12} sm={7} md={12} lg={7}>
                                            <div className={"w-4/5 m-auto"}>
                                                <Grid container>
                                                    <Grid xs={8}>
                                                        {categoryExpenseSummary.map((category) => (
                                                            <div key={category.category_id}
                                                                 className="flex gap-1 items-center">
                                                                <div className="rounded-full w-3 h-3"
                                                                     style={{backgroundColor: colorByCategory.get(category.category_id)}}/>
                                                                <Typography
                                                                    sx={{color: colorByCategory.get(category.category_id)}}>
                                                                    {Math.round(getCategoryExpensePercentage(category.total))}%
                                                                </Typography>
                                                                <Typography color="gray.light">
                                                                    {category.name}
                                                                </Typography>
                                                            </div>
                                                        ))}
                                                    </Grid>

                                                    <Grid xs={4}>
                                                        {categoryExpenseSummary.map((category) => (
                                                            <Typography key={category.category_id} color="gray.light">
                                                                {new Intl.NumberFormat('en-US', {
                                                                    style: 'currency',
                                                                    currency: 'USD'
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

                            {/*Balance, expenses, creation buttons*/}
                            <Grid xs={12} md={6} hidden={!mdUp} lg={5}>
                                <div>
                                    <Grid container>
                                        {/*Balance and expenses*/}
                                        <Grid xs={12}>
                                            <BalanceCard remainder={user ? user.remainder : 0}/>
                                            <div className={"pt-2"}>
                                                <ExpenseCard expenses={user && user.expenses ? user.expenses : 0}/>
                                            </div>
                                        </Grid>

                                        {/**New expense/income buttons*/}
                                        <Grid xs={12}>
                                            <div className={"flex mt-3"}>
                                                <Button color={"secondary"} variant={"contained"}
                                                        startIcon={<AddIcon/>}>
                                                    New expense
                                                </Button>

                                                <Button sx={{marginLeft: "1rem"}} variant={"contained"}
                                                        startIcon={<AddIcon/>}>
                                                    New income
                                                </Button>
                                            </div>
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
                        Expenses
                    </Typography>

                    <div className="pt-3">
                        <ExpensesTable expenses={expenses} categories={user ? user.categories : []}/>
                    </div>
                </Grid>
            </Grid>

        </Container>
    );
}