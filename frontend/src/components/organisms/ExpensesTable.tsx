import { Box, Typography } from "@mui/material";
import {
  GridColDef,
  GridPaginationModel,
  GridRenderCellParams,
  GridRowsProp,
  GridSortModel,
  useGridApiRef,
} from "@mui/x-data-grid";
import { Category, Expense } from "../../types";
import { GridValidRowModel } from "@mui/x-data-grid/models/gridRows";
import { useRef, useState } from "react";
import { Colors } from "../../assets";
import { Button } from "../atoms";
import { CategorySelect } from "./CategorySelect.tsx";
import { useLocation, useNavigate } from "@tanstack/react-router";
import { useGetExpenses } from "../../queries";
import { ErrorSnackbar, Table } from "../molecules";
import { tableDateFormatter } from "../../utils";

type ExpensesTableProps = {
  categories: Category[] | undefined;
  period?: string;
};

export function ExpensesTable({ categories, period }: ExpensesTableProps) {
  const getExpensesQuery = useGetExpenses(period);
  const location = useLocation();

  const expenses: Expense[] | undefined = getExpensesQuery.data?.expenses;
  const route = "/";

  const errSnackbar = {
    open: true,
    title: "Error fetching expenses",
  };

  const startKeysByPage = useRef<{ [page: number]: string }>({ 0: "" });

  const colorsByExpense: Map<string, string> = getColorsByExpense();

  const [selectedCategories, setSelectedCategories] = useState<Category[]>(
    getSelectedCategoriesFromURL(),
  );
  const [paginationModel, setPaginationModel] = useState(getPaginationFromURL());

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
    { field: "amount", headerName: "Amount", width: 150 },
    { field: "name", headerName: "Name", width: 150 },
    {
      field: "category_name",
      headerName: "Category",
      width: 150,
      renderCell: renderCategoryCell,
      sortable: false,
    },
    { field: "notes", headerName: "Notes", flex: 1, minWidth: 150, sortable: false },
    {
      field: "created_date",
      headerName: "Date",
      width: 200,
      valueFormatter: (params) => tableDateFormatter.format(params),
    },
  ];

  const navigate = useNavigate();

  function getSelectedCategoriesFromURL(): Category[] {
    const params = new URLSearchParams(location.search);
    const categoryParams: string[] = params.get("categories")?.split(",") || [];

    if (categoryParams.length === 0 || !categories || categories.length === 0) {
      return [];
    }

    return categories.filter((category) => categoryParams.includes(category.id));
  }

  function getPaginationFromURL(): GridPaginationModel {
    const params = new URLSearchParams(location.search);
    const pageSize = params.get("pageSize") || "10";

    return {
      page: 0,
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
        category_name: expense.category_name ? expense.category_name : "-",
        notes: expense.notes ? expense.notes : "-",
        created_date: new Date(expense.created_date),
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
        to: route,
      });

      return;
    }

    navigate({
      to: route,
      search: {
        ...location.search,
        categories: selectedCategories.map((category) => category.id).join(","),
      },
    });
  }

  function clearFilter(): void {
    setSelectedCategories([]);

    navigate({
      to: route,
    });
  }

  function onPaginationModelChange(newModel: GridPaginationModel) {
    let search = { ...location.search };

    if (newModel.pageSize !== paginationModel.pageSize) {
      startKeysByPage.current = { 0: "" };
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
      to: route,
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
    const search = { ...location.search };

    model.forEach((item) => {
      if (search.sortOrder !== item.sort || search.sortBy !== item.field) {
        //In this case the page order changes, so we need to reset this map
        startKeysByPage.current = { 0: "" };
      }

      navigate({
        to: route,
        search: {
          ...search,
          sortBy: item.field,
          sortOrder: item.sort,
        },
      });

      return;
    });
  }

  function showErrorSnackbar(): boolean {
    if (getExpensesQuery.isError && getExpensesQuery.error.response) {
      return getExpensesQuery.error.response.status !== 404;
    }

    return getExpensesQuery.isError;
  }

  const apiRef = useGridApiRef();
  return (
    <div>
      {showErrorSnackbar() && (
        <ErrorSnackbar
          openProp={errSnackbar.open}
          title={errSnackbar.title}
          message={getExpensesQuery.error ? getExpensesQuery.error.message : ""}
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
          sx={{ marginLeft: "1rem" }}
          onClick={clearFilter}
          color={"darkerGray"}
          variant={"outlined"}
        >
          Clear filter
        </Button>
      </div>
      <div className={"pt-4"}>
        <Box boxShadow={"3"} height={"631px"} width={"100%"} borderRadius={"1rem"}>
          <Table
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
            noRowsMessage={"No expenses found"}
          />
        </Box>
      </div>
    </div>
  );
}
