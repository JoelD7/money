import Grid from "@mui/material/Unstable_Grid2";
import ArrowCircleUpRoundedIcon from "@mui/icons-material/ArrowCircleUpRounded";
import {Typography} from "@mui/material";

type BalanceCardProps = {
    remainder: number
}

export function BalanceCard({remainder}: BalanceCardProps) {
    const customWidth = {
        '&.MuiSvgIcon-root': {
            width: "38px",
            height: "38px",
        },
    }

    return (
        <div>
            <Grid container borderRadius="1rem" p="0.5rem" bgcolor="white.main" boxShadow={"2"}>
                <Grid xs={3}>
                    <Grid height="100%" container alignContent="center" justifyContent="center">
                        {/*@ts-ignore*/}
                        <ArrowCircleUpRoundedIcon sx={customWidth} color="darkGreen"/>
                    </Grid>
                </Grid>

                <Grid xs={9}>
                    <Typography variant="h6" fontWeight="bold">Balance</Typography>
                    <Typography lineHeight="unset" variant="h4" color="darkGreen.main">
                        {new Intl.NumberFormat('en-US', {
                            style: 'currency',
                            currency: 'USD'
                        }).format(remainder)}
                    </Typography>
                </Grid>
            </Grid>
        </div>
    )
}