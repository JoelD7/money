import {
  Box,
  FormControl,
  InputLabel,
  MenuItem,
  Select,
  SelectChangeEvent,
  Typography,
} from "@mui/material";
import { DataGrid, GridCell, GridColDef, GridRowsProp } from "@mui/x-data-grid";
import { Category, Expense } from "../../types";
import { GridValidRowModel } from "@mui/x-data-grid/models/gridRows";
import { GridCellProps } from "@mui/x-data-grid/components/cell/GridCell";
import { useState } from "react";
import { Colors } from "../../assets";
import { v4 as uuidv4 } from "uuid";
import { Button } from "../atoms";

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

  const colorsByExpense : Map<string, string> = getColorsByExpense()
  const [filteredExpenses, setFilteredExpenses] = useState<Expense[]>(expenses);

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

  function getCategoryOptions(): CategorySelectorOption[] {
    const addedCategories = new Set<string>();
    const options: CategorySelectorOption[] = [];

    expenses.forEach((expense) => {
      if (
        expense.category_name &&
        !addedCategories.has(expense.category_name)
      ) {
        options.push({
          label: expense.category_name,
          color: colorsByExpense.get(expense.expense_id) || "gray.main",
        });

        addedCategories.add(expense.category_name);
      }
    });

    return options;
  }

  function onCategorySelectedChange(selected: string[]) {
    if (selected.length === 0) {
      setFilteredExpenses(expenses);
      return;
    }

    const newFilteredExpenses: Expense[] = expenses.filter((expense) => {
      return selected.includes(expense.category_name || "");
    });

    setFilteredExpenses(newFilteredExpenses);
  }

  return (
    <div>
      <CategorySelector
        onSelectedUpdate={onCategorySelectedChange}
        options={getCategoryOptions()}
        label={"Filter by categories"}
      />
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

type CategorySelectorOption = {
  label: string;
  color: string;
};

type CategorySelectorProps = {
  options: CategorySelectorOption[];
  label: string;
  onSelectedUpdate: (selected: string[]) => void;
};

function CategorySelector({
  options,
  label,
  onSelectedUpdate,
}: CategorySelectorProps) {
  const labelId: string = uuidv4();
  const [selected, setSelected] = useState<string[]>([]);
  const colorMap: Map<string, string> = buildColorMap();

  function onSelectedChange(event: SelectChangeEvent<typeof selected>) {
    const {
      target: { value },
    } = event;
    const newValue = typeof value === "string" ? value.split(" ") : value;
    setSelected(newValue);
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

  function clearFilter(): void {
    setSelected([]);
    onSelectedUpdate([]);
  }

  return (
    <>
      <FormControl fullWidth sx={{ background: "white", maxWidth: "460px" }}>
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
            <Box sx={{ display: "flex", flexWrap: "wrap", gap: 0.5 }}>
              {selected.map((value) => (
                <Box
                  key={value}
                  className="p-1 w-fit text-sm rounded-xl"
                  style={{ color: "white" }}
                  sx={{ backgroundColor: getOptionColor(value) }}
                >
                  {value}
                </Box>
              ))}
            </Box>
          )}
        >
          {
            // This is how items will appear on the menu
            options.map((option) => (
              <MenuItem
                key={option.label}
                id={option.label}
                value={option.label}
              >
                <Box
                  className="p-1 w-fit text-sm rounded-xl"
                  style={{ color: "white" }}
                  sx={{ backgroundColor: option.color }}
                >
                  {option.label}
                </Box>
              </MenuItem>
            ))
          }
        </Select>
      </FormControl>

      <div className="flex mt-2">
        <Button variant="outlined" onClick={() => onSelectedUpdate(selected)}>
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
    </>
  );
}
