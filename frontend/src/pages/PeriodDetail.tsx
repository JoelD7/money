import {Box, Typography, useMediaQuery, useTheme} from "@mui/material";
import {Expense, RechartsLabelProps} from "../types";
import Grid from "@mui/material/Unstable_Grid2";
import ArrowCircleUpRoundedIcon from "@mui/icons-material/ArrowCircleUpRounded";
import ArrowCircleDownRoundedIcon from "@mui/icons-material/ArrowCircleDownRounded";
import {Cell, Pie, PieChart, ResponsiveContainer, Tooltip} from "recharts";
import {BalanceCard, Button, ExpenseCard, ExpensesTable} from "../components";
import AddIcon from "@mui/icons-material/Add";
import ArrowCircleLeftIcon from '@mui/icons-material/ArrowCircleLeft';
import json2mq from "json2mq";
import {Colors} from "../assets";

type ExpenseCategorySummary = {
    categoryID: string
    categoryName: string
    percentage: number
    total: number
}

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

    const user = {
        full_name: "Joel",
        username: "test@gmail.com",
        remainder: 14456.21,
        expenses: 7381.27,
        categories: [
            {
                id: "CTGGyouAaIPPWKzxpyxHACS",
                name: "Entertainment",
                color: "#ff8733",
                value: 55
            },
            {
                id: "CTGcSuhjzVmu3WrHLKD5fhS",
                name: "Health",
                color: "#00b85e",
                value: 15
            },
            {
                id: "CTGrR7fO4ndmI0IthJ7Wg8f",
                name: "Utilities",
                color: "#009eb8",
                value: 30
            },
            {
                id: "CTGiBScOP3V16LYBjdIStP9",
                name: "Shopping",
                color: "#8c34eb",
                value: 30
            }
        ],
        current_period: "2023-5"
    }

    const period = {
        "username": "test@gmail.com",
        "period": "asdf",
        "name": "December",
        "start_date": "2023-11-26T00:00:00Z",
        "end_date": "2023-12-24T00:00:00Z",
    }

    const expenses: Expense[] = [
        {
            expenseID: "EX5DK8d8LTywTKC8r87vdS",
            username: "test@gmail.com",
            categoryID: "CTGiBScOP3V16LYBjdIStP9",
            categoryName: "Shopping",
            amount: 12.99,
            name: "Blue pair of socks",
            notes: "Ipsum mollit est pariatur esse ex. Aliqua laborum laboris laboris laboris. Laboris pectum",
            createdDate: new Date("2023-10-27T23:42:54.980596532Z"),
            period: "2023-7",
            updateDate: new Date("0001-01-01T00:00:00Z")
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
            updateDate: new Date("0001-01-01T00:00:00Z")
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
            updateDate: new Date("0001-01-01T00:00:00Z")
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
            updateDate: new Date("0001-01-01T00:00:00Z")
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
            updateDate: new Date("0001-01-01T00:00:00Z")
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
            updateDate: new Date("0001-01-01T00:00:00Z")
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
            updateDate: new Date("0001-01-01T00:00:00Z")
        },
        {
            expenseID: "EXP123",
            username: "test@gmail.com",
            amount: 893,
            name: "Jordan shopping",
            createdDate: new Date("0001-01-01T00:00:00Z"),
            period: "2023-5",
            updateDate: new Date("0001-01-01T00:00:00Z")
        },
        {
            expenseID: "EXP456",
            username: "test@gmail.com",
            amount: 112,
            name: "Uber drive",
            createdDate: new Date("0001-01-01T00:00:00Z"),
            period: "2023-5",
            updateDate: new Date("0001-01-01T00:00:00Z")
        },
        {
            expenseID: "EXP789",
            username: "test@gmail.com",
            amount: 525,
            name: "Lunch",
            createdDate: new Date("0001-01-01T00:00:00Z"),
            period: "2023-5",
            updateDate: new Date("0001-01-01T00:00:00Z")
        }
    ]

    const expenseSummaryByCategory: ExpenseCategorySummary[] = getExpenseSummaryByCategory()
    const colorByCategory: Map<string, string> = getColorByCategory()

    function getExpenseSummaryByCategory(): ExpenseCategorySummary[] {
        let summaryByCategory: Map<string, ExpenseCategorySummary> = new Map<string, ExpenseCategorySummary>()
        summaryByCategory.set("Other", {
            total: 0,
            percentage: 0,
            categoryID: "Other",
            categoryName: "Other",
        })

        let categoryData: ExpenseCategorySummary | undefined

        expenses.forEach((expense) => {
            if (!expense.categoryID) {
                categoryData = summaryByCategory.get("Other")

                if (categoryData) {
                    categoryData.total += expense.amount
                    categoryData.percentage = getPercentageFromExpenses(categoryData.total)
                    summaryByCategory.set("Other", categoryData)
                }

                return
            }

            categoryData = summaryByCategory.get(expense.categoryID)

            if (categoryData) {
                categoryData.total += expense.amount
                categoryData.percentage = getPercentageFromExpenses(categoryData.total)
                summaryByCategory.set(expense.categoryID, categoryData)
                return
            }

            summaryByCategory.set(expense.categoryID, {
                percentage: getPercentageFromExpenses(expense.amount),
                total: expense.amount,
                categoryID: expense.categoryID ? expense.categoryID : "",
                categoryName: expense.categoryName ? expense.categoryName : "",
            })
        })

        return Array.from(summaryByCategory.values())
    }

    function getColorByCategory(): Map<string, string> {
        let m: Map<string, string> = new Map<string, string>()
        m.set("Other", Colors.GRAY_DARK)

        user.categories.forEach((category) => {
            m.set(category.id, category.color)
        })

        return m
    }

    function getPercentageFromExpenses(total: number): number {
        return (total / user.expenses) * 100
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
        return `${new Intl.DateTimeFormat('en-US', {
            month: 'short',
            day: '2-digit',
        }).format(new Date(period.start_date))} - ${new Intl.DateTimeFormat('default', {
            month: 'short',
            day: '2-digit',
        }).format(new Date(period.end_date))}`
    }

    return (
        <>
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
                                                {period.name}
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
                                }).format(user.remainder)}
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
                                }).format(user.expenses)}
                            </Typography>
                        </div>
                    </Box>
                    <div hidden={!smUp}>
                        <Grid container spacing={1}>
                            <Grid xs={6}>
                                <BalanceCard remainder={user.remainder}/>
                            </Grid>
                            <Grid xs={6}>
                                <ExpenseCard expenses={user.expenses}/>
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
                                                    <Pie data={expenseSummaryByCategory}
                                                         label={getCustomLabel}
                                                         dataKey="total"
                                                         nameKey="categoryName"
                                                         cx="50%"
                                                         cy="50%"
                                                         labelLine={false}>
                                                        {expenseSummaryByCategory.map((summary, index) => (
                                                            <Cell key={`cell-${index}`}
                                                                  fill={colorByCategory.get(summary.categoryID)}/>
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
                                                        {expenseSummaryByCategory.map((summary) => (
                                                            <div key={summary.categoryID}
                                                                 className="flex gap-1 items-center">
                                                                <div className="rounded-full w-3 h-3"
                                                                     style={{backgroundColor: colorByCategory.get(summary.categoryID)}}/>
                                                                <Typography
                                                                    sx={{color: colorByCategory.get(summary.categoryID)}}>
                                                                    {Math.ceil(summary?.percentage)}%
                                                                </Typography>
                                                                <Typography color="gray.light">
                                                                    {summary?.categoryName}
                                                                </Typography>
                                                            </div>
                                                        ))}
                                                    </Grid>

                                                    <Grid xs={4}>
                                                        {expenseSummaryByCategory.map((summary) => (
                                                            <Typography key={summary.categoryID} color="gray.light">
                                                                {new Intl.NumberFormat('en-US', {
                                                                    style: 'currency',
                                                                    currency: 'USD'
                                                                }).format(summary?.total)}
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
                                            <BalanceCard remainder={user.remainder}/>
                                            <div className={"pt-2"}>
                                                <ExpenseCard expenses={user.expenses}/>
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
                        <ExpensesTable expenses={expenses}/>
                    </div>
                </Grid>
            </Grid>

        </>
    );
}