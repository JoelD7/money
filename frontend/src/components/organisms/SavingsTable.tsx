import {
  BackgroundRefetchErrorSnackbar,
  ErrorSnackbar,
  Table,
  TableHeader,
} from "../molecules";
import { useGetSavings } from "../../queries";
import { useRef, useState } from "react";
import { Saving } from "../../types";
import {
  GridColDef,
  GridPaginationModel,
  GridRenderCellParams,
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
  faUpRightFromSquare,
} from "@fortawesome/free-solid-svg-icons";
import { GridValidRowModel } from "@mui/x-data-grid/models/gridRows";
import { Link } from "@tanstack/react-router";
import { Button } from "../atoms";

type SavingsTableProps = {
  savingGoalID?: string;
};

export function SavingsTable({ savingGoalID }: SavingsTableProps) {
  const startKeysByPage = useRef<{ [page: number]: string }>({ 0: "" });

  const [sortOrder, setSortOrder] = useState("");
  const [paginationModel, setPaginationModel] = useState({ page: 0, pageSize: 10 });
  const [sortBy, setSortBy] = useState("");

  const startKey: string = startKeysByPage.current[paginationModel.page];
  const pageSize: number = paginationModel.pageSize;
  const getSavingsQuery = useGetSavings(
    startKey,
    pageSize,
    sortOrder,
    sortBy,
    savingGoalID,
  );

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
      renderCell: (params) => <GoalCell params={params} />,
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
        period: saving.period_name,
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

  function openErrorSnackbar(): boolean {
    if (getSavingsQuery.isError && getSavingsQuery.error.response) {
      return getSavingsQuery.error.response.status !== 404;
    }

    return getSavingsQuery.isError;
  }

  return (
    <>
      {openErrorSnackbar() && (
        <ErrorSnackbar
          openProp={true}
          title={"Error fetching savings"}
          message={getSavingsQuery.error ? getSavingsQuery.error.message : ""}
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

type goalCellProps = {
  params: GridRenderCellParams;
};

function GoalCell({ params }: goalCellProps) {
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
