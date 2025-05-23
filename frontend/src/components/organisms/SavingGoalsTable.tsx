import {
  Alert,
  AlertTitle,
  Box,
  capitalize,
  IconButton,
  keyframes,
  LinearProgress,
  linearProgressClasses,
  Snackbar,
  Typography,
} from "@mui/material";
import { Button } from "../atoms";
import { DeleteSavingGoal, Table, TableHeader } from "../molecules";
import { MouseEvent, useRef, useState } from "react";
import { useGetSavingGoals } from "../../queries";
import { PaginationModel, SavingGoal, SnackAlert } from "../../types";
import {
  GridColDef,
  GridColumnHeaderParams,
  GridPaginationModel,
  GridRenderCellParams,
  GridRowId,
  GridRowsProp,
  GridSortModel,
} from "@mui/x-data-grid";
import { Colors } from "../../assets";
import Grid from "@mui/material/Unstable_Grid2";
import { GridValidRowModel } from "@mui/x-data-grid/models/gridRows";
import { NewSavingGoal } from "./NewSavingGoal.tsx";
import { faUpRightFromSquare } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { Link } from "@tanstack/react-router";
import { faTrash } from "@fortawesome/free-solid-svg-icons/faTrash";

export function SavingGoalsTable() {
  const headerIconSize = 15;
  const startKeysByPage = useRef<{ [page: number]: string }>({ 0: "" });

  const dateFormatter = new Intl.DateTimeFormat("en-US", {
    year: "numeric",
    month: "short",
    day: "numeric",
  });
  const currencyFormatter = new Intl.NumberFormat("en-US", {
    style: "currency",
    currency: "USD",
  });
  const percentageFormatter = new Intl.NumberFormat("en-US", {
    maximumFractionDigits: 0,
  });

  const [open, setOpen] = useState(false);
  const [paginationModel, setPaginationModel] = useState<PaginationModel>({
    page: 0,
    pageSize: 10,
  });
  const [sortOrder, setSortOrder] = useState<string>();
  const [sortBy, setSortBy] = useState<string>();
  const [key, setKey] = useState<number>(0);
  const [alert, setAlert] = useState<SnackAlert>({
    open: false,
    type: "success",
    title: "",
  });
  const [hoveredRow, setHoveredRow] = useState<GridRowId>("0");
  const [savingGoalToDelete, setSavingGoalToDelete] = useState<SavingGoal | null>(null);
  const [openDeleteDialog, setOpenDeleteDialog] = useState<boolean>(false);

  const startKey: string | undefined = startKeysByPage.current[paginationModel.page];
  const pageSize: number | undefined = paginationModel.pageSize;

  const getSavingGoalsQuery = useGetSavingGoals(startKey, pageSize, sortOrder, sortBy);
  const savingGoals: SavingGoal[] | undefined = getSavingGoalsQuery.data?.saving_goals;
  const savingGoalsByID: Map<string, SavingGoal> = buildSavingGoalByID(
    savingGoals ? savingGoals : [],
  );

  const columns: GridColDef[] = [
    {
      field: "name",
      headerName: "Name",
      flex: 1,
      minWidth: 180,
      renderHeader: renderNameHeader,
      renderCell: (params) => <NameCell params={params} />,
    },
    {
      field: "target",
      headerName: "Target",
      flex: 0.6,
      minWidth: 180,
      renderHeader: renderTargetHeader,
      valueFormatter: (params) => currencyFormatter.format(params),
    },
    {
      field: "saved",
      headerName: "Saved",
      flex: 0.6,
      minWidth: 150,
      sortable: false,
      renderHeader: renderSavedHeader,
      valueFormatter: (params) => currencyFormatter.format(params),
    },
    {
      field: "progress",
      headerName: "Progress",
      flex: 1,
      minWidth: 250,
      sortable: false,
      renderHeader: renderProgressHeader,
      renderCell: renderProgressCell,
    },
    {
      field: "deadline",
      headerName: "Deadline",
      resizable: false,
      flex: 0.7,
      minWidth: 150,
      renderHeader: renderDeadlineHeader,
      valueFormatter: (params) => dateFormatter.format(params),
    },
    {
      field: "actions",
      flex: 0.1,
      minWidth: 50,
      sortable: false,
      resizable: false,
      disableColumnMenu: true,

      renderHeader: () => null,
      renderCell: (params) => {
        const isVisible = hoveredRow === params.id;
        return (
          <>
            {isVisible && (
              <div key={params.id}>
                <IconButton
                  title={"Delete"}
                  size={"small"}
                  onClick={() => handleDeleteSavingGoal(params.id)}
                >
                  <FontAwesomeIcon icon={faTrash} />
                </IconButton>
              </div>
            )}
          </>
        );
      },
    },
  ];

  function buildSavingGoalByID(savingGoals: SavingGoal[]): Map<string, SavingGoal> {
    const savingGoalByID: Map<string, SavingGoal> = new Map<string, SavingGoal>();

    for (const goal of savingGoals) {
      savingGoalByID.set(goal.saving_goal_id, goal);
    }

    return savingGoalByID;
  }

  function renderNameHeader(params: GridColumnHeaderParams) {
    return (
      <TableHeader
        headerName={params.colDef.headerName || "Name"}
        icon={
          <svg
            xmlns="http://www.w3.org/2000/svg"
            height={`${headerIconSize}`}
            width={`${headerIconSize}`}
            fill={Colors.BLUE}
            viewBox="0 0 448 512"
          >
            <path d="M254 52.8C249.3 40.3 237.3 32 224 32s-25.3 8.3-30 20.8L57.8 416 32 416c-17.7 0-32 14.3-32 32s14.3 32 32 32l96 0c17.7 0 32-14.3 32-32s-14.3-32-32-32l-1.8 0 18-48 159.6 0 18 48-1.8 0c-17.7 0-32 14.3-32 32s14.3 32 32 32l96 0c17.7 0 32-14.3 32-32s-14.3-32-32-32l-25.8 0L254 52.8zM279.8 304l-111.6 0L224 155.1 279.8 304z" />
          </svg>
        }
      />
    );
  }

  function renderTargetHeader(params: GridColumnHeaderParams) {
    return (
      <TableHeader
        headerName={params.colDef.headerName || "Target"}
        icon={
          <svg
            xmlns="http://www.w3.org/2000/svg"
            height={`${headerIconSize}`}
            width={`${headerIconSize}`}
            fill={Colors.BLUE}
            viewBox="0 0 512 512"
          >
            <path d="M448 256A192 192 0 1 0 64 256a192 192 0 1 0 384 0zM0 256a256 256 0 1 1 512 0A256 256 0 1 1 0 256zm256 80a80 80 0 1 0 0-160 80 80 0 1 0 0 160zm0-224a144 144 0 1 1 0 288 144 144 0 1 1 0-288zM224 256a32 32 0 1 1 64 0 32 32 0 1 1 -64 0z" />
          </svg>
        }
      />
    );
  }

  function renderSavedHeader(params: GridColumnHeaderParams) {
    return (
      <TableHeader
        headerName={params.colDef.headerName || "Saved"}
        icon={
          <svg
            height={`${headerIconSize}`}
            width={`${headerIconSize}`}
            fill={Colors.BLUE}
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 576 512"
          >
            <path d="M400 96l0 .7c-5.3-.4-10.6-.7-16-.7L256 96c-16.5 0-32.5 2.1-47.8 6c-.1-2-.2-4-.2-6c0-53 43-96 96-96s96 43 96 96zm-16 32c3.5 0 7 .1 10.4 .3c4.2 .3 8.4 .7 12.6 1.3C424.6 109.1 450.8 96 480 96l11.5 0c10.4 0 18 9.8 15.5 19.9l-13.8 55.2c15.8 14.8 28.7 32.8 37.5 52.9l13.3 0c17.7 0 32 14.3 32 32l0 96c0 17.7-14.3 32-32 32l-32 0c-9.1 12.1-19.9 22.9-32 32l0 64c0 17.7-14.3 32-32 32l-32 0c-17.7 0-32-14.3-32-32l0-32-128 0 0 32c0 17.7-14.3 32-32 32l-32 0c-17.7 0-32-14.3-32-32l0-64c-34.9-26.2-58.7-66.3-63.2-112L68 304c-37.6 0-68-30.4-68-68s30.4-68 68-68l4 0c13.3 0 24 10.7 24 24s-10.7 24-24 24l-4 0c-11 0-20 9-20 20s9 20 20 20l31.2 0c12.1-59.8 57.7-107.5 116.3-122.8c12.9-3.4 26.5-5.2 40.5-5.2l128 0zm64 136a24 24 0 1 0 -48 0 24 24 0 1 0 48 0z" />
          </svg>
        }
      />
    );
  }

  function renderProgressHeader(params: GridColumnHeaderParams) {
    return (
      <TableHeader
        headerName={params.colDef.headerName || "Progress"}
        icon={
          <svg
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 512 512"
            height={`${headerIconSize}`}
            width={`${headerIconSize}`}
            fill={Colors.BLUE}
          >
            <path d="M304 48a48 48 0 1 0 -96 0 48 48 0 1 0 96 0zm0 416a48 48 0 1 0 -96 0 48 48 0 1 0 96 0zM48 304a48 48 0 1 0 0-96 48 48 0 1 0 0 96zm464-48a48 48 0 1 0 -96 0 48 48 0 1 0 96 0zM142.9 437A48 48 0 1 0 75 369.1 48 48 0 1 0 142.9 437zm0-294.2A48 48 0 1 0 75 75a48 48 0 1 0 67.9 67.9zM369.1 437A48 48 0 1 0 437 369.1 48 48 0 1 0 369.1 437z" />
          </svg>
        }
      />
    );
  }

  function renderDeadlineHeader(params: GridColumnHeaderParams) {
    return (
      <TableHeader
        headerName={params.colDef.headerName || "Deadline"}
        icon={
          <svg
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 448 512"
            height={`${headerIconSize}`}
            width={`${headerIconSize}`}
            fill={Colors.BLUE}
          >
            <path d="M128 0c17.7 0 32 14.3 32 32l0 32 128 0 0-32c0-17.7 14.3-32 32-32s32 14.3 32 32l0 32 48 0c26.5 0 48 21.5 48 48l0 48L0 160l0-48C0 85.5 21.5 64 48 64l48 0 0-32c0-17.7 14.3-32 32-32zM0 192l448 0 0 272c0 26.5-21.5 48-48 48L48 512c-26.5 0-48-21.5-48-48L0 192zm80 64c-8.8 0-16 7.2-16 16l0 96c0 8.8 7.2 16 16 16l96 0c8.8 0 16-7.2 16-16l0-96c0-8.8-7.2-16-16-16l-96 0z" />
          </svg>
        }
      />
    );
  }

  function renderProgressCell(params: GridRenderCellParams) {
    let progressPercent: number = 0;
    const goal: SavingGoal | undefined = savingGoalsByID.get(params.id as string);

    if (goal) {
      progressPercent = (goal.progress / goal.target) * 100;
    }

    let toTransformX: number = 100 - progressPercent;
    if (toTransformX < 0){
      toTransformX = 0
    }

    const progressGrow = keyframes`
    from {
      transform: translateX(-100%);
    }
    to {
      transform: translateX(-${toTransformX}%);
    }
  `;

    return (
      <Grid container height={"100%"} width={"100%"} alignItems={"center"}>
        <Grid xs={2}>
          <Typography variant={"body2"}>
            {`${percentageFormatter.format(progressPercent)}%`}
          </Typography>
        </Grid>

        <Grid xs={10}>
          <LinearProgress
            variant={"determinate"}
            value={progressPercent}
            sx={{
              width: "100%",
              height: 6,
              borderRadius: 15,
              [`& .${linearProgressClasses.bar}`]: {
                strokeLinecap: "round",
                animation: `${progressGrow} 1s ease-out forwards`,
                backgroundColor: getProgressBackground(progressPercent),
              },
              [`&.${linearProgressClasses.colorPrimary}`]: {
                backgroundColor: Colors.GRAY,
              },
            }}
          />
        </Grid>
      </Grid>
    );
  }

  function getTableRows(savingGoals: SavingGoal[]): GridRowsProp {
    return savingGoals.map((goal): GridValidRowModel => {
      return {
        id: goal.saving_goal_id,
        name: goal.name,
        target: goal.target,
        progress: goal.progress,
        saved: goal.progress,
        deadline: new Date(goal.deadline),
      };
    });
  }

  function getProgressBackground(progress: number): string {
    if (progress >= 66) {
      return Colors.GREEN;
    } else if (progress >= 33) {
      return Colors.YELLOW;
    } else {
      return Colors.ORANGE;
    }
  }

  function onPaginationModelChange(newModel: GridPaginationModel) {
    if (newModel.pageSize !== paginationModel.pageSize) {
      startKeysByPage.current = { 0: "" };
    }

    const key = startKeysByPage.current[newModel.page];

    const nextKey = getSavingGoalsQuery.data?.next_key;
    if (nextKey && key === undefined) {
      startKeysByPage.current[newModel.page] = nextKey;
    }

    setPaginationModel(newModel);
  }

  function openNewGoalDialog() {
    setOpen(true);
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

  function handleAlert(alert?: SnackAlert | undefined) {
    if (alert) {
      setAlert(alert);
    }
  }

  function handleMouseEnter(event: MouseEvent<HTMLDivElement>) {
    const rowID = event.currentTarget.getAttribute("data-id");
    if (rowID !== null) {
      setHoveredRow(rowID);
    }
  }

  function handleMouseLeave() {
    setHoveredRow("");
  }

  function handleDeleteSavingGoal(rowID: GridRowId) {
    const savingGoal: SavingGoal | undefined = savingGoalsByID.get(rowID as string);
    if (savingGoal) {
      setSavingGoalToDelete(savingGoal);
      setOpenDeleteDialog(true);
    }
  }

  function handleDeleteSavingGoalDialogClose() {
    setOpenDeleteDialog(false);
  }

  return (
    <>
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

      {/*Title and buttons*/}
      <Typography variant={"h5"}>Your saving goals</Typography>
      <div className={"pt-4"}>
        <Button variant={"contained"} onClick={() => openNewGoalDialog()}>
          New goal
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
            slotProps={{
              row: {
                onMouseEnter: handleMouseEnter,
                onMouseLeave: handleMouseLeave,
              },
            }}
            paginationMeta={{
              hasNextPage: getSavingGoalsQuery.data?.next_key !== "",
            }}
            noRowsMessage={"No saving goals found"}
          />
        </Box>
      </div>

      <NewSavingGoal
        key={key}
        open={open}
        onClose={() => {
          setOpen(false);
          setKey(key + 1);
        }}
        onAlert={handleAlert}
      />

      {savingGoalToDelete && (
        <DeleteSavingGoal
          open={openDeleteDialog}
          onClose={handleDeleteSavingGoalDialogClose}
          savingGoal={savingGoalToDelete}
          onAlert={(alert)=>handleAlert(alert)}
        />
      )}
    </>
  );
}

type nameCellProps = {
  params: GridRenderCellParams;
};

function NameCell({ params }: nameCellProps) {
  const [show, setShow] = useState(false);

  return (
    <div
      className={"flex items-center justify-between"}
      onMouseEnter={() => setShow(true)}
      onMouseLeave={() => setShow(false)}
    >
      {params.value}

      <div className={show ? "block" : "hidden"}>
        <Link to="/savings/goals/$savingGoalId" params={{ savingGoalId: params.id }}>
          <Button
            size={"small"}
            variant={"contained"}
            sx={{
              height: "fit-content",
            }}
            endIcon={<FontAwesomeIcon icon={faUpRightFromSquare} />}
          >
            Open
          </Button>
        </Link>
      </div>
    </div>
  );
}
