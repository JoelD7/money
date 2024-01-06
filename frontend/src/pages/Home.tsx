import {Box, Typography, useMediaQuery, useTheme} from "@mui/material";
import Grid from '@mui/material/Unstable_Grid2';
import AddIcon from '@mui/icons-material/Add';
import ArrowCircleUpRoundedIcon from '@mui/icons-material/ArrowCircleUpRounded';
import ArrowCircleDownRoundedIcon from '@mui/icons-material/ArrowCircleDownRounded';
import {Button, Navbar} from "../components";
import {Cell, Pie, PieChart, ResponsiveContainer, Tooltip} from "recharts";
import {Expense} from "../types";
import {DataGrid, GridCell, GridColDef, GridRowsProp} from "@mui/x-data-grid";
import {GridValidRowModel} from "@mui/x-data-grid/models/gridRows";
import {GridCellProps} from "@mui/x-data-grid/components/cell/GridCell";

type RechartsLabelProps = {
    cx: number
    cy: number
    midAngle: number
    innerRadius: number
    outerRadius: number
    percent: number
    index: number
}

export function Home() {
    const theme = useTheme();
    const customWidth = {
        '&.MuiSvgIcon-root': {
            width: "38px",
            height: "38px",
        },
    }

    const gridStyle = {
        '&.MuiDataGrid-root': {
            borderRadius: "1rem",
        },
        '&.MuiDataGrid-root .MuiDataGrid-cellContent': {
            textWrap: "pretty",
            maxHeight: "38px",
        },
        '& .MuiDataGrid-columnHeaderTitle': {
            fontSize: "large",
        }
    }

    const xs: boolean = useMediaQuery(theme.breakpoints.up('xs'));
    const xsOnly: boolean = useMediaQuery(theme.breakpoints.only('xs'));
    const mdUp: boolean = useMediaQuery(theme.breakpoints.up('md'));


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
                color: "#009eb8",
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

    const columns: GridColDef[] = [
        {field: 'amount', headerName: 'Amount', width: 150},
        {field: 'categoryName', headerName: 'Category', width: 150},
        {field: 'notes', headerName: 'Notes', width: 150},
        {field: 'createdDate', headerName: 'Date', width: 200},
    ];

    function getTableRows(expenses: Expense[]): GridRowsProp {
        return expenses.map((expense): GridValidRowModel => {
            return {
                id: expense.expenseID,
                amount: new Intl.NumberFormat('en-US', {
                    style: 'currency', currency: 'USD'
                }).format(expense.amount),
                categoryName: expense.categoryName ? expense.categoryName : "-",
                notes: expense.notes ? expense.notes : "-",
                createdDate: new Intl.DateTimeFormat('en-GB', {
                    weekday: "short",
                    year: "numeric",
                    month: "numeric",
                    day: "numeric",
                    hour: 'numeric',
                    minute: 'numeric',
                }).format(expense.createdDate),
            }
        })
    }

    function customCellComponent(props: GridCellProps) {
        const {field, children} = props;

        return (
            field === "categoryName" ?
                <GridCell {...props}>
                    <Box sx={{
                        backgroundColor: getCellBackgroundColor(String(props.rowId)),
                        padding: "0.25rem 0.5rem",
                        borderRadius: "9999px",
                    }}>
                        <Typography fontSize={"14px"} color={"white.main"}>
                            {props.value}
                        </Typography>
                    </Box>
                </GridCell> :
                <GridCell {...props}>
                    {children}
                </GridCell>
        )
    }

    function getCellBackgroundColor(rowID: string): string {
        let categoryName: string = ""
        let categoryColor: string = ""

        expenses.forEach((expense) => {
            if (expense.expenseID === rowID) {
                categoryName = expense.categoryName ? expense.categoryName : ""
                return
            }
        })

        user.categories.forEach((category) => {
            if (category.name === categoryName) {
                categoryColor = category.color
                return
            }
        })

        return categoryColor
    }

    return (
        <>
            <Navbar>
                <Typography lineHeight="unset" variant="h4">
                    Overview
                </Typography>
            </Navbar>

            <Grid container spacing={1}>
                {/*Balance*/}
                <Grid xs={12} sm={6} hidden={mdUp}>
                    <div>
                        <Grid container borderRadius="1rem" p="0.5rem" bgcolor="gray.main">
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
                        <Grid container mt={xsOnly ? "0.5rem" : ""} borderRadius="1rem" p="0.5rem" bgcolor="gray.main">
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
                <Grid xs={12}>
                    <div>
                        <Grid container spacing={1}>
                            {/*Chart section*/}
                            <Grid xs={12} md={6}>
                                <div>
                                    <Grid container borderRadius="1rem" p="1rem" boxShadow="3" mt="1rem">
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
                            <Grid xs={12} md={6}>
                                <div>
                                    <Grid container mt={"1rem"} spacing={1}>
                                        {/*Balance*/}
                                        <Grid xs={12} hidden={!mdUp}>
                                            <div>
                                                <Grid container borderRadius="1rem" p="0.5rem" bgcolor="gray.main">
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
                                                      p="0.5rem" bgcolor="gray.main">
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
                <Grid xs={12}>
                    <div>
                        <Grid container mt={"2rem"}>
                            <Typography variant={"h4"}>
                                Latest
                            </Typography>

                            <Box boxShadow={"3"} width={"100%"} borderRadius={"1rem"} mt={"0.5rem"}>
                                <DataGrid sx={gridStyle}
                                          rows={getTableRows(expenses)}
                                          columns={columns}
                                          slots={{
                                              cell: customCellComponent,
                                          }}
                                />
                            </Box>
                        </Grid>
                    </div>
                </Grid>
            </Grid>

        </>
    );
}