import { Typography } from "@mui/material";
import {
  BackgroundRefetchErrorSnackbar,
  Container,
  Navbar,
  NoRowsDataGrid,
} from "../components";
import { useGetIncome } from "../queries";
import { Income, IncomeList } from "../types";
import {
  DataGrid,
  GridColDef,
  GridPaginationModel,
  GridRowsProp,
} from "@mui/x-data-grid";
import { useRef, useState } from "react";
import { useLocation, useNavigate } from "@tanstack/react-router";
import { GridValidRowModel } from "@mui/x-data-grid/models/gridRows";

export function IncomeTable() {
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

  const location = useLocation();
  const [paginationModel, setPaginationModel] = useState(
    getPaginationFromURL(),
  );

  const navigate = useNavigate();
  const startKeysByPage = useRef<{ [page: number]: string }>({ 0: "" });

  const getIncome = useGetIncome();
  const incomeList: IncomeList | undefined = getIncome.data?.data;

  const columns: GridColDef[] = [
    { field: "amount", headerName: "Amount", width: 150 },
    { field: "name", headerName: "Name", width: 150 },
    { field: "period", headerName: "Period", width: 150 },
    { field: "notes", headerName: "Notes", flex: 1, minWidth: 150 },
    { field: "createdDate", headerName: "Date", width: 200 },
  ];

  function getPaginationFromURL(): GridPaginationModel {
    const params = new URLSearchParams(location.search);
    const pageSize = params.get("pageSize") || "10";
    const page = params.get("page") || "0";

    return {
      page: parseInt(page),
      pageSize: parseInt(pageSize),
    };
  }

  function onPaginationModelChange(newModel: GridPaginationModel) {
    let search = { ...location.search };

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
      to: "/income",
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

    const nextKey = getIncome.data?.data.next_key;
    if (nextKey) {
      startKeysByPage.current[newModel.page] = nextKey;
      return nextKey;
    }

    return "";
  }

  function getTableRows(income: Income[]): GridRowsProp {
    return income.map((inc): GridValidRowModel => {
      return {
        id: inc.income_id,
        amount: new Intl.NumberFormat("en-US", {
          style: "currency",
          currency: "USD",
        }).format(inc.amount),
        name: inc.name,
        notes: inc.notes ? inc.notes : "-",
        period: inc.period,
        createdDate: new Intl.DateTimeFormat("en-GB", {
          weekday: "short",
          year: "numeric",
          month: "numeric",
          day: "numeric",
          hour: "numeric",
          minute: "numeric",
        }).format(new Date(inc.created_date)),
      };
    });
  }

  function showRefetchErrorSnackbar() {
    return false;
  }

  return (
    <Container>
      <BackgroundRefetchErrorSnackbar show={showRefetchErrorSnackbar()} />
      <Navbar />

      <Typography variant={"h3"}>Income</Typography>
      <div style={{ height: "fit-content" }}>
        <DataGrid
          sx={gridStyle}
          loading={getIncome.isFetching}
          columns={columns}
          initialState={{
            pagination: {
              rowCount: -1,
              paginationModel,
            },
          }}
          rows={getTableRows(incomeList?.income ? incomeList?.income : [])}
          pageSizeOptions={[5, 10, 25, 50]}
          paginationMode="server"
          paginationModel={paginationModel}
          onPaginationModelChange={onPaginationModelChange}
          paginationMeta={{
            hasNextPage: getIncome.data?.data.next_key !== "",
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
      </div>
    </Container>
  );
}
