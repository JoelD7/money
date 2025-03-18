import {
  BackgroundRefetchErrorSnackbar,
  Button,
  Container,
  ErrorSnackbar,
  Navbar,
  NewSaving,
  PageTitle,
  SavingGoalsTable,
  Snackbar,
  Table,
  TableHeader,
} from "../components";
import Grid from "@mui/material/Unstable_Grid2";
import { Box, Typography } from "@mui/material";
import { Saving, SnackAlert } from "../types";
import {
  GridColDef,
  GridPaginationModel,
  GridRowsProp,
  GridSortModel,
} from "@mui/x-data-grid";
import { useRef, useState } from "react";
import { useGetSavings } from "../queries";
import { GridValidRowModel } from "@mui/x-data-grid/models/gridRows";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faBullseye,
  faCalendar,
  faClock,
  faDollarSign,
} from "@fortawesome/free-solid-svg-icons";
import { Colors } from "../assets";
import { utils } from "../utils";

export function Savings() {
  const startKeysByPage = useRef<{ [page: number]: string }>({ 0: "" });
  const currencyFormatter = new Intl.NumberFormat("en-US", {
    style: "currency",
    currency: "USD",
  });

  const [open, setOpen] = useState(false);
  const [key, setKey] = useState(0);
  const [alert, setAlert] = useState<SnackAlert>({
    open: false,
    type: "success",
    title: "",
  });
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
        return utils.tableDateFormatter.format(params);
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

  function showRefetchErrorSnackbar() {
    return getSavingsQuery.isRefetchError;
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

  function openNewSavingDialog() {
    setOpen(true);
  }

  return (
    <Container>
      <Navbar />
      <BackgroundRefetchErrorSnackbar show={showRefetchErrorSnackbar()} />

      {getSavingsQuery.isError && (
        <ErrorSnackbar
          openProp={true}
          title={"Error fetching savings"}
          message={getSavingsQuery.error.message}
        />
      )}

      <Snackbar
        open={alert.open}
        title={alert.title}
        message={alert.message}
        severity={alert.type}
        onClose={() => setAlert({ ...alert, open: false })}
      />

      <Grid container position={"relative"} spacing={1} marginTop={"20px"}>
        {/*Title and summary*/}
        <Grid xs={12}>
          <PageTitle>Savings</PageTitle>
        </Grid>

        {/*Saving goals table*/}
        <Grid xs={12}>
          <SavingGoalsTable />
        </Grid>

        {/*Latest savings*/}
        <Grid xs={12} mt={"2rem"}>
          {/*Title and buttons*/}
          <Typography variant={"h5"}>Latest savings</Typography>
          <div className={"pt-4"}>
            <Button variant={"contained"} onClick={() => openNewSavingDialog()}>
              New saving
            </Button>
            <Box height={"fit-content"} paddingTop={"10px"}>
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
            </Box>
          </div>
        </Grid>
      </Grid>

      <NewSaving
        key={key}
        open={open}
        onClose={() => {
          setOpen(false);
          setKey(key + 1);
        }}
        onAlert={(a) => {
          if (a) {
            setAlert({ ...a });
          }
        }}
      />
    </Container>
  );
}
