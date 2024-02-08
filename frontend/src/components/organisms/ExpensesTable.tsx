import {Box, FormControl, InputLabel, MenuItem, Select, SelectChangeEvent, Typography} from "@mui/material";
import {DataGrid, GridCell, GridColDef, GridRowsProp} from "@mui/x-data-grid";
import {Expense} from "../../types";
import {GridValidRowModel} from "@mui/x-data-grid/models/gridRows";
import {GridCellProps} from "@mui/x-data-grid/components/cell/GridCell";
import {useState} from "react";
import {Colors} from "../../assets";
import {v4 as uuidv4} from "uuid";

type ExpensesTableProps = {
    expenses: Expense[];
}

export function ExpensesTable({expenses}: ExpensesTableProps) {
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

    const [colorsByExpense, setColorsByExpense] = useState<Map<string, string>>(getColorsByExpense())
    const [filteredExpenses, setFilteredExpenses] = useState<Expense[]>(expenses)

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

    function getCategoryOptions(): CategorySelectorOption[] {
        let addedCategories = new Set<String>()
        let options: CategorySelectorOption[] = []

        expenses.forEach((expense) => {
            if (expense.categoryName && !addedCategories.has(expense.categoryName)) {
                options.push({
                    label: expense.categoryName,
                    color: colorsByExpense.get(expense.expenseID) || "gray.main"
                })

                addedCategories.add(expense.categoryName)
            }
        })

        return options
    }

    function onCategorySelectedChange(selected: string[]) {
        if (selected.length === 0) {
            setFilteredExpenses(expenses)
            return
        }

        let newFilteredExpenses: Expense[] = expenses.filter((expense) => {
            return selected.includes(expense.categoryName || "")
        })

        setFilteredExpenses(newFilteredExpenses)
    }

    return (
        <div>
            <CategorySelector onSelectedUpdate={onCategorySelectedChange} options={getCategoryOptions()}
                              label={"Filter by categories"}/>
            <Box boxShadow={"3"} width={"100%"} borderRadius={"1rem"}
                 mt={"0.5rem"}>

                <DataGrid sx={gridStyle}
                          rows={getTableRows(filteredExpenses)}
                          columns={columns}
                          slots={{
                              cell: customCellComponent,
                          }}
                />
            </Box>
        </div>
    )
}

type CategorySelectorOption = {
    label: string;
    color: string;
}

type CategorySelectorProps = {
    options: CategorySelectorOption[];
    label: string;
    onSelectedUpdate: (selected: string[]) => void;
}

function CategorySelector({options, label, onSelectedUpdate}: CategorySelectorProps) {
    const labelId: string = uuidv4();
    const [selected, setSelected] = useState<string[]>([]);
    const colorMap: Map<string, string> = buildColorMap();

    function onSelectedChange(event: SelectChangeEvent<typeof selected>) {
        const {target: {value}} = event;
        let newValue = typeof value === 'string' ? value.split(' ') : value
        onSelectedUpdate(newValue)
        setSelected(newValue)
    }

    function buildColorMap(): Map<string, string> {
        const colorMap = new Map<string, string>();
        options.forEach((option) => {
            colorMap.set(option.label, option.color);
        });

        return colorMap;
    }

    function getOptionColor(value: string): string {
        return colorMap.get(value) || "gray.main";
    }

    return (
        <>
            <FormControl fullWidth sx={{background: "white"}}>
                <InputLabel id={labelId}>{label}</InputLabel>
                <Select
                    labelId={labelId}
                    id={label}
                    label={label}
                    value={selected}
                    onChange={onSelectedChange}
                    multiple
                    renderValue={(selected) => (
                        // This is how items will appear on the select input
                        <Box sx={{display: 'flex', flexWrap: 'wrap', gap: 0.5}}>
                            {selected.map((value) => (
                                <Box key={value} className="p-1 w-fit text-sm rounded-xl" style={{color: "white"}}
                                     sx={{backgroundColor: getOptionColor(value)}}>
                                    {value}
                                </Box>
                            ))}
                        </Box>
                    )}
                >
                    {
                        // This is how items will appear on the menu
                        options.map((option) => (
                            <MenuItem key={option.label} id={option.label} value={option.label}>
                                <Box className="p-1 w-fit text-sm rounded-xl" style={{color: "white"}}
                                     sx={{backgroundColor: option.color}}>
                                    {option.label}
                                </Box>
                            </MenuItem>
                        ))

                    }
                </Select>
            </FormControl>
        </>
    );
}