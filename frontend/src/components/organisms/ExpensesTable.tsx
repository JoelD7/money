import {Box, Typography} from "@mui/material";
import {DataGrid, GridCell, GridColDef, GridRowsProp} from "@mui/x-data-grid";
import {Expense} from "../../types";
import {GridValidRowModel} from "@mui/x-data-grid/models/gridRows";
import {GridCellProps} from "@mui/x-data-grid/components/cell/GridCell";
import {useState} from "react";
import {Colors} from "../../assets";

export function ExpensesTable() {
    const gridStyle = {
        '&.MuiDataGrid-root': {
            borderRadius: "1rem",
            backgroundColor: "#ffffff"
        },
        '&.MuiDataGrid-root .MuiDataGrid-cellContent': {
            textWrap: "pretty",
            maxHeight: "38px",
        },
        '& .MuiDataGrid-columnHeaderTitle': {
            fontSize: "large",
        }
    }

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

    const [colorsByExpense, setColorsByExpense] = useState<Map<string, string>>(getColorsByExpense())

    const columns: GridColDef[] = [
        {field: 'amount', headerName: 'Amount', width: 150},
        {field: 'categoryName', headerName: 'Category', width: 150},
        {field: 'notes', headerName: 'Notes', flex: 1, minWidth: 150},
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
        let color: string | undefined = colorsByExpense.get(rowID)
        if (color) {
            return color
        }

        return Colors.WHITE
    }

    function getColorsByExpense(): Map<string, string> {
        const colorsByExpense: Map<string, string> = new Map<string, string>()
        expenses.forEach((expense) => {
            user.categories.forEach((category) => {
                if (category.name === expense.categoryName) {
                    colorsByExpense.set(expense.expenseID, category.color)
                }
            })
        })

        return colorsByExpense
    }

    return (
        <Box boxShadow={"3"} width={"100%"} borderRadius={"1rem"}
             mt={"0.5rem"}>
            <DataGrid sx={gridStyle}
                      rows={getTableRows(expenses)}
                      columns={columns}
                      slots={{
                          cell: customCellComponent,
                      }}
            />
        </Box>
    )
}