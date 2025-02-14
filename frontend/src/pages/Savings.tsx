import {
  BackgroundRefetchErrorSnackbar,
  Button,
  Container,
  Navbar,
  SavingGoalCard,
} from "../components";
import Grid from "@mui/material/Unstable_Grid2";
import { Typography } from "@mui/material";
import SavingsIcon from "@mui/icons-material/Savings";
import FilterListIcon from "@mui/icons-material/FilterList";

export function Savings() {
  const customWidth = {
    "&.MuiSvgIcon-root": {
      width: "28px",
      height: "28px",
      fill: "#024511",
    },
  };

  function showRefetchErrorSnackbar() {
    return false;
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
          <div className={"flex pt-2"}>
            <Button variant={"contained"}>New goal</Button>
            <div className={"ml-2"}>
              <Button variant={"outlined"} startIcon={<FilterListIcon />}>
                Filters
              </Button>
            </div>
          </div>

          {/*Saving cards*/}
          <div className={"flex"}>
            <SavingGoalCard goal={870} progress={221} />
          </div>
        </Grid>
      </Grid>
    </Container>
  );
}
