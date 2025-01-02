import {Box, Typography} from "@mui/material";
import {
    DataGrid,
    GridColDef,
    GridPaginationModel,
    GridRenderCellParams,
    GridRowsProp,
    GridSortModel,
    useGridApiRef,
} from "@mui/x-data-grid";
import {Category, Expense} from "../../types";
import {GridValidRowModel} from "@mui/x-data-grid/models/gridRows";
import {useRef, useState} from "react";
import {Colors} from "../../assets";
import {Button, NoRowsDataGrid} from "../atoms";
import {CategorySelect} from "./CategorySelect.tsx";
import {useLocation, useNavigate} from "@tanstack/react-router";
import {useGetExpenses} from "../../queries";
import {ErrorSnackbar} from "../molecules";

type ExpensesTableProps = {
    categories: Category[] | undefined;
    period: string;
};

export function ExpensesTable({categories, period}: ExpensesTableProps) {
    const gridStyle = {
        "&.MuiDataGrid-root": {
            borderRadius: "1rem",
            backgroundColor: "#ffffff",
            minHeight: "220px",
        },
        "&.MuiDataGrid-root .MuiDataGrid-cellContent": {
            textWrap: "pretty",
            maxHeight: "38px",
        },
        "& .MuiDataGrid-columnHeaderTitle": {
            fontSize: "large",
        },
    };

    const getExpensesQuery = useGetExpenses(period);
    const location = useLocation();

    const expenses: Expense[] | undefined = getExpensesQuery.data?.expenses;

    const errSnackbar = {
        open: true,
        title: "Error fetching expenses",
    }

    const colorsByExpense: Map<string, string> = getColorsByExpense();

    const [selectedCategories, setSelectedCategories] = useState<Category[]>(
        getSelectedCategoriesFromURL(),
    );
    const [paginationModel, setPaginationModel] = useState(
        getPaginationFromURL(),
    );
    const startKeysByPage = useRef<{ [page: number]: string }>({0: ""});

    function renderCategoryCell(params: GridRenderCellParams) {
        const categoryColor = getCellBackgroundColor(String(params.id));
        return (
            <div className={"h-full flex items-center justify-center"}>
                <Box
                    sx={{
                        backgroundColor: categoryColor,
                        padding: "0.25rem 0.5rem",
                        borderRadius: "9999px",
                        display: "flex",
                        justifyContent: "center",
                    }}
                >
                    <Typography fontSize={"14px"} color={"white.main"}>
                        {params.value}
                    </Typography>
                </Box>
            </div>
        );
    }

    const columns: GridColDef[] = [
        {field: "amount", headerName: "Amount", width: 150},
        {field: "name", headerName: "Name", width: 150},
        {
            field: "categoryName",
            headerName: "Category",
            width: 150,
            renderCell: renderCategoryCell,
        },
        {field: "notes", headerName: "Notes", flex: 1, minWidth: 150},
        {
            field: "createdDate", headerName: "Date", width: 200,
            valueFormatter: (params) => {
                return new Intl.DateTimeFormat("en-GB", {
                    weekday: "short",
                    year: "numeric",
                    month: "numeric",
                    day: "numeric",
                    hour: "numeric",
                    minute: "numeric",
                }).format(params)
            }
        },
    ];

    const navigate = useNavigate();

    function getSelectedCategoriesFromURL(): Category[] {
        const params = new URLSearchParams(location.search);
        const categoryParams: string[] = params.get("categories")?.split(",") || [];

        if (categoryParams.length === 0 || !categories || categories.length === 0) {
            return [];
        }

        return categories.filter((category) =>
            categoryParams.includes(category.id),
        );
    }

    function getPaginationFromURL(): GridPaginationModel {
        const params = new URLSearchParams(location.search);
        const pageSize = params.get("pageSize") || "10";
        const page = params.get("page") || "0";

        return {
            page: parseInt(page),
            pageSize: parseInt(pageSize),
        };
    }

    function getTableRows(expenses: Expense[]): GridRowsProp {
        return expenses.map((expense): GridValidRowModel => {
            return {
                id: expense.expense_id,
                amount: new Intl.NumberFormat("en-US", {
                    style: "currency",
                    currency: "USD",
                }).format(expense.amount),
                name: expense.name,
                categoryName: expense.category_name ? expense.category_name : "-",
                notes: expense.notes ? expense.notes : "-",
                createdDate: new Date(expense.created_date),
            };
        });
    }

    function getCellBackgroundColor(rowID: string): string {
        const color: string | undefined = colorsByExpense.get(rowID);
        if (color) {
            return color;
        }

        return Colors.WHITE;
    }

    function getColorsByExpense(): Map<string, string> {
        const colorsByExpense: Map<string, string> = new Map<string, string>();
        if (!expenses || !categories) {
            return colorsByExpense;
        }

        expenses.forEach((expense) => {
            categories.forEach((category) => {
                if (category.name === expense.category_name) {
                    colorsByExpense.set(expense.expense_id, category.color);
                }
            });
        });

        return colorsByExpense;
    }

    function applyFilters() {
        if (!expenses && selectedCategories.length === 0) {
            return;
        }

        if (selectedCategories.length === 0) {
            navigate({
                to: "/",
            });

            return;
        }

        navigate({
            to: "/",
            search: {
                ...location.search,
                categories: selectedCategories.map((category) => category.id).join(","),
            },
        });
    }

    function clearFilter(): void {
        setSelectedCategories([]);

        navigate({
            to: "/",
        });
    }

    function onPaginationModelChange(newModel: GridPaginationModel) {
        let search = {...location.search};

        if (newModel.pageSize !== paginationModel.pageSize) {
            search = {
                ...search,
                pageSize: newModel.pageSize,
            };
        }

        const startKey = getStartKey(newModel);
        if (newModel.page !== paginationModel.page) {
            search = {
                ...search,
                startKey,
            };
        }

        navigate({
            to: "/",
            search,
        }).then(() => {
            setPaginationModel(newModel);
        });
    }

    function getStartKey(newModel: GridPaginationModel): string | undefined {
        if (newModel.page === 0) {
            return undefined;
        }

        const mappedKey = startKeysByPage.current[newModel.page];
        if (mappedKey) {
            return mappedKey;
        }

        const nextKey = getExpensesQuery.data?.next_key;
        if (nextKey) {
            startKeysByPage.current[newModel.page] = nextKey;
            return nextKey;
        }

        return "";
    }

    function onSortModelChange(model: GridSortModel) {
        const search = {...location.search};

        model.forEach((item) => {
            if (search.sortOrder !== item.sort || search.sortBy !== item.field) {
                //In this case the page order changes, so we need to reset this map
                startKeysByPage.current = {0: ""};
            }

            navigate({
                to: "/",
                search: {
                    ...search, sortBy: item.field, sortOrder: item.sort
                },
            })

            return
        })
    }

    const apiRef = useGridApiRef();
    return (
        <div>

            {getExpensesQuery.isError && (
                <ErrorSnackbar
                    openProp={errSnackbar.open}
                    title={errSnackbar.title}
                    message={getExpensesQuery.error.message}
                />
            )}

            <CategorySelect
                width={"700px"}
                multiple
                selected={selectedCategories}
                onSelectedUpdate={(selected) => setSelectedCategories(selected)}
                categories={categories ? categories : []}
                label={"Filter by categories"}
            />
            <div className="flex mt-2">
                <Button variant="outlined" onClick={applyFilters}>
                    Apply filter
                </Button>
                <Button
                    sx={{marginLeft: "1rem"}}
                    onClick={clearFilter}
                    color={"darkerGray"}
                    variant={"outlined"}
                >
                    Clear filter
                </Button>
            </div>
            <div className={"pt-4"}>
                <Box
                    boxShadow={"3"}
                    height={"631px"}
                    width={"100%"}
                    borderRadius={"1rem"}
                >
                    <DataGrid
                        sx={gridStyle}
                        apiRef={apiRef}
                        loading={getExpensesQuery.isFetching}
                        columns={columns}
                        sortingMode={"server"}
                        onSortModelChange={(model) => onSortModelChange(model)}
                        initialState={{
                            pagination: {
                                rowCount: -1,
                                paginationModel,
                            },
                        }}
                        rows={getTableRows(expenses ? expenses : [])}
                        pageSizeOptions={[10, 25, 50]}
                        paginationMode="server"
                        paginationModel={paginationModel}
                        onPaginationModelChange={onPaginationModelChange}
                        paginationMeta={{
                            hasNextPage: getExpensesQuery.data?.next_key !== "",
                        }}
                        slots={{
                            noRowsOverlay: NoRowsDataGrid,
                        }}
                        slotProps={{
                            noRowsOverlay: {
                                sx: {
                                    height: "100px",
                                },
                            },
                            loadingOverlay: {
                                variant: "linear-progress",
                                noRowsVariant: "skeleton",
                            },
                        }}
                    />
                </Box>
            </div>
        </div>
    );
}
