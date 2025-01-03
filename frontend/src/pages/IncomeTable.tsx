import {FormControl, InputLabel, MenuItem, Select, Typography,} from "@mui/material";
import {BackgroundRefetchErrorSnackbar, Container, ErrorSnackbar, Navbar, NoRowsDataGrid,} from "../components";
import {useGetIncome} from "../queries";
import {Income, IncomeList} from "../types";
import {DataGrid, GridColDef, GridPaginationModel, GridRowsProp, GridSortModel,} from "@mui/x-data-grid";
import {useRef, useState} from "react";
import {useLocation, useNavigate} from "@tanstack/react-router";
import {GridValidRowModel} from "@mui/x-data-grid/models/gridRows";
import {v4 as uuidv4} from "uuid";

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

    const labelId: string = uuidv4();

    const incomeListErrSnackbar = {
        open: true,
        title: "Error fetching income. Refresh the page to try again",
    }

    const location = useLocation();
    const [paginationModel, setPaginationModel] = useState(
        getPaginationFromURL(),
    );
    const [selectedPeriod, setSelectedPeriod] = useState(
        getCurrentPeriodFromURL(),
    );

    const navigate = useNavigate();
    const startKeysByPage = useRef<{ [page: number]: string }>({0: ""});

    const getIncome = useGetIncome();
    const incomeList: IncomeList | undefined = getIncome.data;

    const periods: string[] = getPeriodArray()

    function getPeriodArray(): string[] {
        if (!incomeList) {
            return []
        }

        return incomeList.periods ? incomeList.periods : []
    }

    const columns: GridColDef[] = [
        {field: "amount", headerName: "Amount", width: 150},
        {field: "name", headerName: "Name", width: 150},
        {field: "period", headerName: "Period", width: 150, sortable: false},
        {field: "notes", headerName: "Notes", flex: 1, minWidth: 150, sortable: false},
        {
            field: "created_date", headerName: "Date", width: 200,
            valueFormatter: (params) => {
                return new Intl.DateTimeFormat("en-GB", {
                    weekday: "short",
                    year: "numeric",
                    month: "numeric",
                    day: "numeric",
                    hour: "numeric",
                    minute: "numeric",
                }).format(params)
            }
        },
    ];

    function getCurrentPeriodFromURL(): string {
        const params = new URLSearchParams(location.search);
        return params.get("period") || "";
    }

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
        let search = {...location.search};

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

        const nextKey = getIncome.data?.next_key;
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
                created_date: new Date(inc.created_date),
            };
        });
    }

    function showRefetchErrorSnackbar() {
        return false;
    }

    function onSelectedPeriodChange(newPeriod: string) {
        if (selectedPeriod === newPeriod) {
            return
        }

        navigate({
            to: "/income",
            search: {
                ...location.search,
                period: newPeriod,
            },
        }).catch((e) => {
            console.error("[money] - Navigating to /income failed: ", e)
        })

        setSelectedPeriod(newPeriod)
    }

    function onSortModelChange(model: GridSortModel) {
        const search = {...location.search};

        model.forEach((item) => {
            if (search.sortOrder !== item.sort || search.sortBy !== item.field) {
                //In this case the page order changes, so we need to reset this map
                startKeysByPage.current = {0: ""};
            }

            navigate({
                to: "/income",
                search: {
                    ...search, sortBy: item.field, sortOrder: item.sort
                },
            })

            return
        })
    }

    return (
        <Container>
            <BackgroundRefetchErrorSnackbar show={showRefetchErrorSnackbar()}/>
            <Navbar/>

            {
                getIncome.isError && incomeList === undefined && (
                    <ErrorSnackbar
                        openProp={incomeListErrSnackbar.open}
                        title={incomeListErrSnackbar.title}
                        message={getIncome.error.message}
                    />
                )
            }

            <Typography variant={"h3"} sx={{margin: "50px 0px 20px 0px"}}>
                Income
            </Typography>

            {/* Period selector*/}
            <div className={"pb-2"}>
                <FormControl sx={{width: '150px'}}>
                    <InputLabel id={labelId}>Period</InputLabel>

                    <Select
                        labelId={labelId}
                        id={"Period"}
                        MenuProps={{
                            PaperProps: {
                                sx: {
                                    maxHeight: 150,
                                }
                            }
                        }}
                        label={"Period"}
                        value={periods.length > 0 ? selectedPeriod : ""}
                        onChange={(e) => onSelectedPeriodChange(e.target.value)}
                    >
                        {Array.isArray(periods) && periods.map((p) => (
                            <MenuItem key={p} id={p} value={p}>
                                {p}
                            </MenuItem>
                        ))}
                    </Select>
                </FormControl>
            </div>

            <div style={{height: "631px"}}>
                <DataGrid
                    sx={gridStyle}
                    loading={getIncome.isFetching}
                    columns={columns}
                    sortingMode={"server"}
                    onSortModelChange={(model)=> onSortModelChange(model)}
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
                        hasNextPage: getIncome.data?.next_key !== "",
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
