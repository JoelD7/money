import {
  BackgroundRefetchErrorSnackbar,
  Button,
  Container,
  Navbar,
  SavingGoalsTable,
  Table,
} from "../components";
import Grid from "@mui/material/Unstable_Grid2";
import { Alert, AlertTitle, Box, capitalize, Snackbar, Typography } from "@mui/material";
import { SnackAlert } from "../types";
import { GridSortModel } from "@mui/x-data-grid";
import { useRef, useState } from "react";

export function Savings() {
  const startKeysByPage = useRef<{ [page: number]: string }>({ 0: "" });

  const [open, setOpen] = useState(false);
  const [alert, setAlert] = useState<SnackAlert>({
    open: false,
    type: "success",
    title: "",
  });
  const [sortOrder, setSortOrder] = useState("");
  const [sortBy, setSortBy] = useState("");

  const getSavingsQuery = useGetSavings();

  function showRefetchErrorSnackbar() {
    return false;
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

  function openNewSavingDialog() {
    setOpen(true);
  }

  return (
    <Container>
      <Navbar />
      <BackgroundRefetchErrorSnackbar show={showRefetchErrorSnackbar()} />
      <Snackbar
        open={alert.open}
        onClose={() => setAlert({ ...alert, open: false })}
        autoHideDuration={6000}
        anchorOrigin={{ vertical: "top", horizontal: "right" }}
      >
        <Alert variant={"filled"} severity={alert.type}>
          <AlertTitle>{capitalize(alert.type)}</AlertTitle>
          {alert.title}
        </Alert>
      </Snackbar>

      <Grid container position={"relative"} spacing={1} marginTop={"20px"}>
        {/*Title and summary*/}
        <Grid xs={12}>
          <Typography mt={"2rem"} variant={"h3"}>
            Savings
          </Typography>
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
                loading={getSavingGoalsQuery.isFetching}
                columns={columns}
                rows={getTableRows(savingGoals ? savingGoals : [])}
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
                  hasNextPage: getSavingGoalsQuery.data?.next_key !== "",
                }}
                noRowsMessage={"No saving goals found"}
              />
            </Box>
          </div>
        </Grid>
      </Grid>
    </Container>
  );
}
