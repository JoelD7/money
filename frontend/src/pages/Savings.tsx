import {
  BackgroundRefetchErrorSnackbar,
  Button,
  Container,
  Navbar,
  SavingGoalCard,
  Table,
} from "../components";
import Grid from "@mui/material/Unstable_Grid2";
import { Typography } from "@mui/material";
import SavingsIcon from "@mui/icons-material/Savings";
import { useGetSavingGoals } from "../queries";
import { SavingGoal } from "../types";
import { GridColDef, GridRowsProp } from "@mui/x-data-grid";
import { GridValidRowModel } from "@mui/x-data-grid/models/gridRows";

export function Savings() {
  const customWidth = {
    "&.MuiSvgIcon-root": {
      width: "28px",
      height: "28px",
      fill: "#024511",
    },
  };

  const columns: GridColDef[] = [
    { field: "name", headerName: "Name", width: 150 },
    { field: "target", headerName: "Target", width: 150 },
    { field: "progress", headerName: "Progress", width: 150 },
    { field: "saved", headerName: "Saved", width: 150 },
    {
      field: "deadline",
      headerName: "Deadline",
      width: 150,
      valueFormatter: (params) => {
        return new Intl.DateTimeFormat("en-GB", {
          weekday: "short",
          year: "numeric",
          month: "numeric",
          day: "numeric",
          hour: "numeric",
          minute: "numeric",
        }).format(params);
      },
    },
  ];

  const getSavingGoalsQuery = useGetSavingGoals();
  const savingGoals: SavingGoal[] | undefined = getSavingGoalsQuery.data?.saving_goals;

  function showRefetchErrorSnackbar() {
    return false;
  }

  function getTableRows(savingGoals: SavingGoal[]): GridRowsProp {
    return savingGoals.map((goal): GridValidRowModel => {
      return {
        id: goal.saving_goal_id,
        name: goal.name,
        target: goal.target,
        progress: "0%",
        saved: "0",
        deadline: goal.deadline,
      };
    });
  }

  return (
    <Container>
      <Navbar />
      <BackgroundRefetchErrorSnackbar show={showRefetchErrorSnackbar()} />

      <Grid container position={"relative"} spacing={1} marginTop={"20px"}>
        {/*Title and summary*/}
        <Grid xs={12}>
          <Typography mt={"2rem"} variant={"h3"}>
            Savings
          </Typography>

          {/*Summary*/}
          <div className={"mt-2"}>
            <Grid
              container
              borderRadius="0.5rem"
              p="1rem"
              bgcolor="white.main"
              maxWidth={"450px"}
              boxShadow={"2"}
              justifyContent={"space-between"}
            >
              <Grid xs={6}>
                <div className={"flex items-center"}>
                  <SavingsIcon sx={customWidth} />
                  <Typography variant={"h6"}>Total Savings</Typography>
                </div>
              </Grid>

              <Grid xs={4}>
                <Typography lineHeight="unset" variant="h6" color="primary">
                  {new Intl.NumberFormat("en-US", {
                    style: "currency",
                    currency: "USD",
                  }).format(585018)}
                </Typography>
              </Grid>
            </Grid>
          </div>
        </Grid>

        {/*Saving cards*/}
        <Grid xs={12} pt={"3rem"}>
          {/*Title and buttons*/}
          <Typography variant={"h5"}>Your saving goals</Typography>
          <Table
            loading={getSavingGoalsQuery.isFetching}
            columns={columns}
            rows={getTableRows(expenses ? expenses : [])}
          />
        </Grid>
      </Grid>
    </Container>
  );
}
