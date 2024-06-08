import { Box, Typography } from "@mui/material";
import { DataGrid, GridCell, GridColDef, GridRowsProp } from "@mui/x-data-grid";
import { Category, Expense } from "../../types";
import { GridValidRowModel } from "@mui/x-data-grid/models/gridRows";
import { GridCellProps } from "@mui/x-data-grid/components/cell/GridCell";
import { useState } from "react";
import { Colors } from "../../assets";
import { Button } from "../atoms";
import { CategorySelect } from "./CategorySelect.tsx";

type ExpensesTableProps = {
  expenses: Expense[];
  categories: Category[];
};

export function ExpensesTable({ expenses, categories }: ExpensesTableProps) {
  const gridStyle = {
    "&.MuiDataGrid-root": {
      borderRadius: "1rem",
      backgroundColor: "#ffffff",
    },
    "&.MuiDataGrid-root .MuiDataGrid-cellContent": {
      textWrap: "pretty",
      maxHeight: "38px",
    },
    "& .MuiDataGrid-columnHeaderTitle": {
      fontSize: "large",
    },
  };

  const colorsByExpense: Map<string, string> = getColorsByExpense();
  const [filteredExpenses, setFilteredExpenses] = useState<Expense[]>(expenses);
  const [selectedCategories, setSelectedCategories] = useState<string[]>([]);

  const columns: GridColDef[] = [
    { field: "amount", headerName: "Amount", width: 150 },
    { field: "categoryName", headerName: "Category", width: 150 },
    { field: "notes", headerName: "Notes", flex: 1, minWidth: 150 },
    { field: "createdDate", headerName: "Date", width: 200 },
  ];

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
    expenses.forEach((expense) => {
      categories.forEach((category) => {
        if (category.name === expense.category_name) {
          colorsByExpense.set(expense.expense_id, category.color);
        }
      });
    });

    return colorsByExpense;
  }

  function onCategorySelectedChange() {
    if (selectedCategories.length === 0) {
      setFilteredExpenses(expenses);
      return;
    }

    const newFilteredExpenses: Expense[] = expenses.filter((expense) => {
      return selectedCategories.includes(expense.category_name || "");
    });

    setFilteredExpenses(newFilteredExpenses);
  }

  function clearFilter(): void {
    setSelectedCategories([]);
  }

  return (
    <div>
      <CategorySelect
        width={"700px"}
        multiple
        selected={selectedCategories}
        onSelectedUpdate={(selected) => setSelectedCategories(selected)}
        categories={categories}
        label={"Filter by categories"}
      />
      <div className="flex mt-2">
        <Button variant="outlined" onClick={() => onCategorySelectedChange()}>
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
        <Box boxShadow={"3"} width={"100%"} borderRadius={"1rem"}>
          <DataGrid
            sx={gridStyle}
            rows={getTableRows(filteredExpenses)}
            columns={columns}
            slots={{
              cell: customCellComponent,
            }}
          />
        </Box>
      </div>
    </div>
  );
}
