import {Box, Typography} from "@mui/material";
import { DataGrid, GridCell, GridColDef, GridRowsProp } from "@mui/x-data-grid";
import { Category, Expense } from "../../types";
import { GridValidRowModel } from "@mui/x-data-grid/models/gridRows";
import { GridCellProps } from "@mui/x-data-grid/components/cell/GridCell";
import { useState } from "react";
import { Colors } from "../../assets";
import {Button, NoRowsDataGrid} from "../atoms";
import { CategorySelect } from "./CategorySelect.tsx";
import {useGetExpenses} from "./queries.ts";
import {useLocation, useNavigate} from "@tanstack/react-router";

type ExpensesTableProps = {
  categories: Category[] | undefined;
};

export function ExpensesTable({  categories }: ExpensesTableProps) {
  const gridStyle = {
    "&.MuiDataGrid-root": {
      borderRadius: "1rem",
      backgroundColor: "#ffffff",
      minHeight: "220px"
    },
    "&.MuiDataGrid-root .MuiDataGrid-cellContent": {
      textWrap: "pretty",
      maxHeight: "38px",
    },
    "& .MuiDataGrid-columnHeaderTitle": {
      fontSize: "large",
    },
  };

  const getExpensesQuery = useGetExpenses()
  const location = useLocation();

  const expenses: Expense[] | undefined = getExpensesQuery.data?.data.expenses;

  const colorsByExpense: Map<string, string> = getColorsByExpense();

  const [selectedCategories, setSelectedCategories] = useState<Category[]>(getSelectedCategoriesFromURL());

  const columns: GridColDef[] = [
    { field: "amount", headerName: "Amount", width: 150 },
    { field: "categoryName", headerName: "Category", width: 150 },
    { field: "notes", headerName: "Notes", flex: 1, minWidth: 150 },
    { field: "createdDate", headerName: "Date", width: 200 },
  ];

  const navigate = useNavigate()

  function getSelectedCategoriesFromURL(): Category[]{
    const params = new URLSearchParams(location.search);
    const categoryParams: string[] = params.get("categories")?.split(",") || [];

    if (categoryParams.length === 0 || !categories || categories.length === 0){
      return []
    }

    return categories.filter((category)=> categoryParams.includes(category.id))
  }

  function getTableRows(expenses: Expense[]): GridRowsProp {
    return expenses.map((expense): GridValidRowModel => {
      return {
        id: expense.expense_id,
        amount: new Intl.NumberFormat("en-US", {
          style: "currency",
          currency: "USD",
        }).format(expense.amount),
        categoryName: expense.category_name ? expense.category_name : "-",
        notes: expense.notes ? expense.notes : "-",
        createdDate: new Intl.DateTimeFormat("en-GB", {
          weekday: "short",
          year: "numeric",
          month: "numeric",
          day: "numeric",
          hour: "numeric",
          minute: "numeric",
        }).format(new Date(expense.created_date)),
      };
    });
  }

  function customCellComponent(props: GridCellProps) {
    const { field, children } = props;

    return field === "categoryName" ? (
      <GridCell {...props}>
        <Box
          sx={{
            backgroundColor: getCellBackgroundColor(String(props.rowId)),
            padding: "0.25rem 0.5rem",
            borderRadius: "9999px",
          }}
        >
          <Typography fontSize={"14px"} color={"white.main"}>
            {props.value}
          </Typography>
        </Box>
      </GridCell>
    ) : (
      <GridCell {...props}>{children}</GridCell>
    );
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

    if (selectedCategories.length === 0){
      navigate({
        to: "/expenses",
      })

      return
    }

    navigate({
      to: "/expenses",
      search: {
        categories: selectedCategories.map((category) => category.id).join(","),
      },
    })

  }

  function clearFilter(): void {
    setSelectedCategories([]);

    navigate({
      to: "/expenses",
    })
  }

  return (
      <div>
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
          <Box boxShadow={"3"} width={"100%"} borderRadius={"1rem"}>
            <DataGrid
                sx={gridStyle}
                loading={getExpensesQuery.isFetching}
                slots={{
                  cell: customCellComponent,
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
                rows={getTableRows(expenses ? expenses : [])}
                columns={columns}
            />
          </Box>
        </div>
      </div>
  );
}
