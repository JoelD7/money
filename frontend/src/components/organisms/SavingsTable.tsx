import {BackgroundRefetchErrorSnackbar, ErrorSnackbar, Table, TableHeader} from "../molecules";
import { useGetSavings } from "../../queries";
import { useRef, useState } from "react";
import { Saving } from "../../types";
import {
  GridColDef,
  GridPaginationModel,
  GridRowsProp,
  GridSortModel,
} from "@mui/x-data-grid";
import { currencyFormatter, tableDateFormatter } from "../../utils";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { Colors } from "../../assets";
import {
  faBullseye,
  faCalendar,
  faClock,
  faDollarSign,
} from "@fortawesome/free-solid-svg-icons";
import { GridValidRowModel } from "@mui/x-data-grid/models/gridRows";

export function SavingsTable() {
  const startKeysByPage = useRef<{ [page: number]: string }>({ 0: "" });

  const [sortOrder, setSortOrder] = useState("");
  const [paginationModel, setPaginationModel] = useState({ page: 0, pageSize: 10 });
  const [sortBy, setSortBy] = useState("");

  const startKey: string = startKeysByPage.current[paginationModel.page];
  const pageSize: number = paginationModel.pageSize;
  const getSavingsQuery = useGetSavings(startKey, pageSize, sortOrder, sortBy);

  const savings: Saving[] | undefined = getSavingsQuery.data?.savings;

  const columns: GridColDef[] = [
    {
      field: "amount",
      headerName: "Amount",
      flex: 1,
      minWidth: 180,
      valueFormatter: (params) => currencyFormatter.format(params),
      renderHeader: () => (
        <TableHeader
          headerName={"Amount"}
          icon={<FontAwesomeIcon color={Colors.BLUE} icon={faDollarSign} />}
        />
      ),
    },
    {
      field: "period",
      headerName: "Period",
      flex: 1,
      minWidth: 180,
      sortable: false,
      renderHeader: () => (
        <TableHeader
          headerName={"Period"}
          icon={<FontAwesomeIcon color={Colors.BLUE} icon={faClock} />}
        />
      ),
    },
    {
      field: "goal",
      headerName: "Goal",
      flex: 1,
      minWidth: 180,
      sortable: false,
      renderHeader: () => (
        <TableHeader
          headerName={"Goal"}
          icon={<FontAwesomeIcon color={Colors.BLUE} icon={faBullseye} />}
        />
      ),
    },
    {
      field: "created_date",
      headerName: "Created date",
      flex: 1,
      minWidth: 180,
      valueFormatter: (params) => {
        return tableDateFormatter.format(params);
      },
      renderHeader: () => (
        <TableHeader
          headerName={"Created date"}
          icon={<FontAwesomeIcon color={Colors.BLUE} icon={faCalendar} />}
        />
      ),
    },
  ];

  function getTableRows(savings: Saving[]): GridRowsProp {
    return savings.map((saving): GridValidRowModel => {
      return {
        id: saving.saving_id,
        amount: saving.amount,
        period: saving.period,
        goal: saving.saving_goal_name,
        created_date: saving.created_date ? new Date(saving.created_date) : null,
      };
    });
  }

  function onSortModelChange(newModel: GridSortModel) {
    newModel.forEach((model) => {
      if (model.sort !== sortOrder && model.sort) {
        setSortOrder(model.sort);
        //In this case the page order changes, so we need to reset this map because the pagination order changes
        startKeysByPage.current = { 0: "" };
      }

      if (model.field !== sortBy) {
        setSortBy(model.field);
      }
    });
  }

  function onPaginationModelChange(newModel: GridPaginationModel) {
    if (newModel.pageSize !== paginationModel.pageSize) {
      startKeysByPage.current = { 0: "" };
    }

    const nextKey = getSavingsQuery.data?.next_key;
    if (nextKey) {
      startKeysByPage.current[newModel.page] = nextKey;
    }

    setPaginationModel(newModel);
  }

  function showRefetchErrorSnackbar() {
    return getSavingsQuery.isRefetchError;
  }

  return (
    <>
      {getSavingsQuery.isError && (
          <ErrorSnackbar
              openProp={true}
              title={"Error fetching savings"}
              message={getSavingsQuery.error.message}
          />
      )}

      <BackgroundRefetchErrorSnackbar show={showRefetchErrorSnackbar()} />

      <Table
        sortingMode={"server"}
        loading={getSavingsQuery.isFetching}
        columns={columns}
        rows={getTableRows(savings ? savings : [])}
        onSortModelChange={onSortModelChange}
        paginationModel={paginationModel}
        initialState={{
          pagination: {
            rowCount: -1,
            paginationModel,
          },
        }}
        pageSizeOptions={[5, 10, 25]}
        paginationMode={"server"}
        onPaginationModelChange={onPaginationModelChange}
        paginationMeta={{
          hasNextPage: getSavingsQuery.data?.next_key !== "",
        }}
        noRowsMessage={"No savings found"}
      />
    </>
  );
}
