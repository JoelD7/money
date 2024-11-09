import {CircularProgress, FormControl, InputLabel, MenuItem, Select, Typography,} from "@mui/material";
import {BackgroundRefetchErrorSnackbar, Container, ErrorSnackbar, Navbar, NoRowsDataGrid,} from "../components";
import {useGetIncome, useGetPeriods} from "../queries";
import {Income, IncomeList, Period, PeriodList} from "../types";
import {DataGrid, GridColDef, GridPaginationModel, GridRowsProp,} from "@mui/x-data-grid";
import React, {useRef, useState} from "react";
import {useLocation, useNavigate} from "@tanstack/react-router";
import {GridValidRowModel} from "@mui/x-data-grid/models/gridRows";
import {v4 as uuidv4} from "uuid";
import {InfiniteData} from "@tanstack/react-query";

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

    const periodErrSnackbar = {
        open: true,
        title: "Error fetching periods. Refresh the page to try again",
    }

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

    const getPeriods = useGetPeriods();
    const periods: Period[] = getPeriodArray(getPeriods.data)

    function getPeriodArray(data?: InfiniteData<PeriodList>): Period[] {
        if (!data || !data.pages) {
            return new Array<Period>()
        }

        const arr = new Array<Period>()
        let includedSelected = false

        data.pages.forEach((page) => {
            page.periods.forEach((period) => {
                if (period.name === selectedPeriod) {
                    includedSelected = true
                }

                arr.push(period)
            })
        })

        if (selectedPeriod !== "" && !includedSelected) {
            arr.push({
                name: selectedPeriod,
                period: selectedPeriod,
                created_date: "",
                end_date: "",
                start_date: "",
                updated_date: "",
                username: ""
            })
        }

        return Array.of(...arr)
    }

    const columns: GridColDef[] = [
        {field: "amount", headerName: "Amount", width: 150},
        {field: "name", headerName: "Name", width: 150},
        {field: "period", headerName: "Period", width: 150},
        {field: "notes", headerName: "Notes", flex: 1, minWidth: 150},
        {field: "createdDate", headerName: "Date", width: 200},
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

    function fetchPeriods(e: React.UIEvent<HTMLDivElement>) {
        const target = e.currentTarget;
        const isAtBottom = target.scrollHeight - Math.ceil(target.scrollTop) === target.clientHeight;
        if (isAtBottom && getPeriods.hasNextPage && !getPeriods.isFetching) {
            getPeriods.fetchNextPage()
        }
    }

    return (
        <Container>
            <BackgroundRefetchErrorSnackbar show={showRefetchErrorSnackbar()}/>
            <Navbar/>

            {getPeriods.isError && (
                <ErrorSnackbar
                    openProp={periodErrSnackbar.open}
                    title={periodErrSnackbar.title}
                    message={getPeriods.error.message}
                />
            )}

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
                                onScroll: fetchPeriods,
                                sx: {
                                    maxHeight: 150,
                                }
                            }
                        }}
                        label={"Period"}
                        value={periods.length > 0 ? selectedPeriod : ""}
                        onChange={(e) => setSelectedPeriod(e.target.value)}
                    >
                        {Array.isArray(periods) && periods.map((p) => (
                            <MenuItem key={p.period} id={p.name} value={p.name}>
                                {p.name}
                            </MenuItem>
                        ))}
                        {getPeriods.isFetchingNextPage &&
                            <MenuItem key={'loading'} id={'loading'} value={'loading'}>
                                <CircularProgress sx={{margin: 'auto'}}/>
                            </MenuItem>}
                    </Select>
                </FormControl>
            </div>

            <div style={{height: "631px"}}>
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
