import {IconButton, Typography, useMediaQuery, useTheme} from "@mui/material";
import {RechartsLabelProps} from "../types";
import Grid from "@mui/material/Unstable_Grid2";
import ArrowCircleUpRoundedIcon from "@mui/icons-material/ArrowCircleUpRounded";
import ArrowCircleDownRoundedIcon from "@mui/icons-material/ArrowCircleDownRounded";
import {Cell, Pie, PieChart, ResponsiveContainer, Tooltip} from "recharts";
import {Button, ExpensesTable} from "../components";
import AddIcon from "@mui/icons-material/Add";
import ArrowCircleLeftIcon from '@mui/icons-material/ArrowCircleLeft';
import json2mq from "json2mq";

export function PeriodDetail() {
    const theme = useTheme();
    const customWidth = {
        '&.MuiSvgIcon-root': {
            width: "38px",
            height: "38px",
        },
    }

    const backButtonStyle = {
        '&.MuiSvgIcon-root': {
            width: "38px",
            height: "38px",
        },
    }

    const xsOnly: boolean = useMediaQuery(theme.breakpoints.only('xs'));
    const mdUp: boolean = useMediaQuery(theme.breakpoints.up('md'));
    const xlCustom = useMediaQuery(
        json2mq({
            maxWidth: 2300,
        }),
    );

    const user = {
        full_name: "Joel",
        username: "test@gmail.com",
        remainder: 14456.21,
        expenses: 8563.05,
        categories: [
            {
                id: "CTGzJeEzCNz6HMTiPKwgPmj",
                name: "Entertainment",
                color: "#ff8733",
                value: 55
            },
            {
                id: "CTGtClGT160UteOl02jIH4F",
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
                id: "CTGrR7fO4ndmI0IthJ7Wg8fs",
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

    const chartHeight: number = 250
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
    };

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
            {/*Page title*/}
            <Grid container>
                <Grid xs={2}>
                    <div className={"flex justify-center items-center h-full"}>
                        {/*@ts-ignore*/}
                        <ArrowCircleLeftIcon sx={customWidth} color="darkGreen"/>
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

            <Grid container spacing={1} justifyContent={"center"}>
                {/*Balance*/}
                <Grid xs={12} sm={6} hidden={mdUp}>
                    <div>
                        <Grid container borderRadius="1rem" p="0.5rem" bgcolor="white.main" boxShadow={"2"}>
                            <Grid xs={3}>
                                <Grid height="100%" container alignContent="center" justifyContent="center">
                                    {/*@ts-ignore*/}
                                    <ArrowCircleUpRoundedIcon sx={customWidth} color="darkGreen"/>
                                </Grid>
                            </Grid>

                            <Grid xs={9}>
                                <Typography variant="h6" fontWeight="bold">Balance</Typography>
                                <Typography lineHeight="unset" variant="h4" color="darkGreen.main">
                                    {new Intl.NumberFormat('en-US', {
                                        style: 'currency',
                                        currency: 'USD'
                                    }).format(user.remainder)}
                                </Typography>
                            </Grid>
                        </Grid>
                    </div>
                </Grid>

                {/*Expenses*/}
                <Grid xs={12} sm={6} hidden={mdUp}>
                    <div>
                        <Grid container mt={xsOnly ? "0.5rem" : ""} borderRadius="1rem" p="0.5rem" bgcolor="white.main"
                              boxShadow={"2"}>
                            <Grid xs={3}>
                                <Grid height="100%" container alignContent="center" justifyContent="center">
                                    {/*@ts-ignore*/}
                                    <ArrowCircleDownRoundedIcon sx={customWidth} color="red"/>
                                </Grid>
                            </Grid>

                            <Grid xs={9}>
                                <Typography variant="h6" fontWeight="bold">Expenses</Typography>
                                <Typography lineHeight="unset" variant="h4" color="red.main">
                                    {new Intl.NumberFormat('en-US', {
                                        style: 'currency',
                                        currency: 'USD'
                                    }).format(user.expenses)}
                                </Typography>
                            </Grid>
                        </Grid>
                    </div>
                </Grid>

                {/*Chart, Current balance and expenses*/}
                <Grid xs={12} maxWidth={"880px"}>
                    <div>
                        <Grid container spacing={1}>
                            {/*Chart section*/}
                            <Grid xs={12} md={6} maxWidth={"430px"}>
                                <div>
                                    <Grid container bgcolor={"white.main"} borderRadius="1rem"
                                          p="1rem" boxShadow="3"
                                          mt="1rem">
                                        <Grid xs={12}>
                                            <Typography variant="h4">
                                                {period.name}
                                            </Typography>
                                            <Typography color="gray.light">
                                                {getPeriodDates()}
                                            </Typography>
                                        </Grid>
                                        {/*Chart*/}
                                        <Grid xs={12} height={chartHeight}>
                                            <ResponsiveContainer width="100%" height="100%">
                                                <PieChart width={350} height={chartHeight}>
                                                    <Pie data={user.categories} label={getCustomLabel} dataKey="value"
                                                         nameKey="name"
                                                         cx="50%"
                                                         cy="50%"
                                                         labelLine={false}
                                                         fill="#8884d8">
                                                        {user.categories.map((category, index) => (
                                                            <Cell key={`cell-${index}`} fill={category.color}/>
                                                        ))}
                                                    </Pie>
                                                    <Tooltip/>
                                                </PieChart>
                                            </ResponsiveContainer>
                                        </Grid>

                                        {/*Chart legend*/}
                                        <Grid xs={12}>
                                            <Grid container width="100%" className="justify-between">
                                                {/*Categories*/}
                                                <Grid xs={6}>
                                                    {user.categories.map((category) => (
                                                        <div key={category.id} className="flex gap-1 items-center">
                                                            <div className="rounded-full w-3 h-3"
                                                                 style={{backgroundColor: category.color}}/>
                                                            <Typography color="gray.light">
                                                                {category.name}
                                                            </Typography>
                                                        </div>
                                                    ))}

                                                </Grid>
                                                {/*Details button*/}
                                                <Grid xs={6}>
                                                    <Grid container className="items-end h-full">
                                                        <Button variant="outlined"
                                                                sx={{
                                                                    textTransform: "capitalize",
                                                                    borderRadius: "1rem",
                                                                    height: "fit-content"
                                                                }}>
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
                                            <div>
                                                <Grid container borderRadius="1rem" p="0.5rem" bgcolor="white.main"
                                                      boxShadow={"2"}>
                                                    <Grid xs={3}>
                                                        <Grid height="100%" container alignContent="center"
                                                              justifyContent="center">
                                                            <ArrowCircleUpRoundedIcon sx={customWidth}
                                                                //@ts-ignore
                                                                                      color="darkGreen"/>
                                                        </Grid>
                                                    </Grid>

                                                    <Grid xs={9}>
                                                        <Typography variant="h6" fontWeight="bold">Balance</Typography>
                                                        <Typography lineHeight="unset" variant="h4"
                                                                    color="darkGreen.main">
                                                            {new Intl.NumberFormat('en-US', {
                                                                style: 'currency',
                                                                currency: 'USD'
                                                            }).format(user.remainder)}
                                                        </Typography>
                                                    </Grid>
                                                </Grid>
                                            </div>
                                        </Grid>

                                        {/*Expenses*/}
                                        <Grid xs={12} hidden={!mdUp}>
                                            <div>
                                                <Grid container mt={xsOnly ? "0.5rem" : ""} borderRadius="1rem"
                                                      p="0.5rem" bgcolor="white.main" boxShadow={"2"}>
                                                    <Grid xs={3}>
                                                        <Grid height="100%" container alignContent="center"
                                                              justifyContent="center">
                                                            {/*@ts-ignore*/}
                                                            <ArrowCircleDownRoundedIcon sx={customWidth} color="red"/>
                                                        </Grid>
                                                    </Grid>

                                                    <Grid xs={9}>
                                                        <Typography variant="h6" fontWeight="bold">Expenses</Typography>
                                                        <Typography lineHeight="unset" variant="h4" color="red.main">
                                                            {new Intl.NumberFormat('en-US', {
                                                                style: 'currency',
                                                                currency: 'USD'
                                                            }).format(user.expenses)}
                                                        </Typography>
                                                    </Grid>
                                                </Grid>
                                            </div>
                                        </Grid>

                                        {/**New expense/income buttons*/}
                                        <Grid xs={12}>
                                            <Button color={"secondary"} variant={"contained"}
                                                    startIcon={<AddIcon/>}>
                                                New expense
                                            </Button>

                                            <Button sx={{marginLeft: "1rem"}} variant={"contained"}
                                                    startIcon={<AddIcon/>}>
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
                    <ExpensesTable/>
                </Grid>
            </Grid>

        </>
    );
}