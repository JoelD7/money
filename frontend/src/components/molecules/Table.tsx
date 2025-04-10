import { DataGrid, DataGridProps } from "@mui/x-data-grid";
import { NoRowsDataGrid } from "../atoms";

type CustomDataGridProps = {
  noRowsMessage?: string;
} & DataGridProps;

export function Table({ noRowsMessage, slots, slotProps,  ...props }: CustomDataGridProps) {
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

  function NoGrid() {
    return (
      <NoRowsDataGrid message={noRowsMessage ? noRowsMessage : "No data available"} />
    );
  }

  return (
    <DataGrid
      sx={gridStyle}
      {...props}
      slots={{
        noRowsOverlay: NoGrid,
        ...slots,
      }}
      slotProps={{
        ...slotProps,
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
  );
}
