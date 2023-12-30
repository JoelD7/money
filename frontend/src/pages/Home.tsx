import {Typography, useMediaQuery, useTheme} from "@mui/material";
import Grid from '@mui/material/Unstable_Grid2';
import ArrowCircleUpRoundedIcon from '@mui/icons-material/ArrowCircleUpRounded';
import ArrowCircleDownRoundedIcon from '@mui/icons-material/ArrowCircleDownRounded';
import {Navbar} from "../components";

export function Home() {
    const theme = useTheme();
    const customWidth = {
        '&.MuiSvgIcon-root': {
            width: "38px",
            height: "38px",
        },
    }
    const md = useMediaQuery(theme.breakpoints.up('md'));

    const user = {
        "full_name": "Joel",
        "username": "test@gmail.com",
        "remainder": 14456.21,
        "expenses": 8563.05,
        "categories": [
            {
                "id": "CTGzJeEzCNz6HMTiPKwgPmj",
                "name": "Entertainment",
                "color": "#ff8733"
            },
            {
                "id": "CTGtClGT160UteOl02jIH4F",
                "name": "Health",
                "color": "#00b85e"
            },
            {
                "id": "CTGrR7fO4ndmI0IthJ7Wg8f",
                "name": "Utilities",
                "color": "#009eb8"
            }
        ],
        "current_period": "2023-5"
    }

    const period = {
        "username": "test@gmail.com",
        "period": "asdf",
        "name": "December",
        "start_date": "2023-11-26T00:00:00Z",
        "end_date": "2023-12-24T00:00:00Z",
    }

    function getPeriodDates(): string {
        return `${new Intl.DateTimeFormat('en-US', {
            month: 'short',
            day: '2-digit',
        }).format(new Date(period.start_date))} - ${new Intl.DateTimeFormat('default', {
            month: 'short',
            day: '2-digit',
        }).format(new Date(period.end_date))}`
    }

    return (
        <>
            <Navbar>
                <Typography lineHeight="unset" variant="h4">
                    Overview
                </Typography>
            </Navbar>

            <Grid container borderRadius="1rem" p="0.5rem" bgcolor="gray.main">
                <Grid xs={3}>
                    <Grid height="100%" container alignContent="center" justifyContent="center">
                        {/*@ts-ignore*/}
                        <ArrowCircleUpRoundedIcon sx={customWidth} color="darkGreen"/>
                    </Grid>
                </Grid>

                <Grid xs={9}>
                    <Typography variant="h6" fontWeight="bold">Balance</Typography>
                    <Typography lineHeight="unset" variant="h4" color="darkGreen.main">
                        {new Intl.NumberFormat('en-US', {style: 'currency', currency: 'USD'}).format(user.remainder)}
                    </Typography>
                </Grid>
            </Grid>

            <Grid container mt="0.5rem" borderRadius="1rem" p="0.5rem" bgcolor="gray.main">
                <Grid xs={3}>
                    <Grid height="100%" container alignContent="center" justifyContent="center">
                        {/*@ts-ignore*/}
                        <ArrowCircleDownRoundedIcon sx={customWidth} color="red"/>
                    </Grid>
                </Grid>

                <Grid xs={9}>
                    <Typography variant="h6" fontWeight="bold">Expenses</Typography>
                    <Typography lineHeight="unset" variant="h4" color="red.main">
                        {new Intl.NumberFormat('en-US', {style: 'currency', currency: 'USD'}).format(user.expenses)}
                    </Typography>
                </Grid>
            </Grid>

            {/*Chart section*/}
            <Grid container borderRadius="1rem" p="1rem" boxShadow="3" mt="1rem">
                <Grid xs={12}>
                    <Typography variant="h4">
                        {period.name}
                    </Typography>
                    <Typography color="gray.light">
                        {getPeriodDates()}
                    </Typography>
                </Grid>
            </Grid>

        </>
    );
}