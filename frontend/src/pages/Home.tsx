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
        "remainder": 14456.21,
        "expenses": 8563.05,
    }

    return (
        <>
            <Navbar>
                <Typography lineHeight="unset" variant="h4">
                    Overview
                </Typography>
            </Navbar>

            <Grid container borderRadius="10px" p="0.5rem" bgcolor="gray.main">
                <Grid xs={3}>
                    <Grid height="100%" container alignContent="center" justifyContent="center">
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

            <Grid container mt="1rem" borderRadius="10px" p="0.5rem" bgcolor="gray.main">
                <Grid xs={3}>
                    <Grid height="100%" container alignContent="center" justifyContent="center">
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
            <div></div>

        </>
    );
}