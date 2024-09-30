import Grid from "@mui/material/Unstable_Grid2";
import ArrowCircleLeftOutlinedIcon from "@mui/icons-material/ArrowCircleLeftOutlined";
import { Typography, useMediaQuery, useTheme } from "@mui/material";
import { CashFlowSkeleton } from "../atoms";
import {Colors} from "../../assets";

type ExpenseCardProps = {
  expenses?: number;
  loading: boolean;
};

export function ExpenseCard({ expenses, loading }: ExpenseCardProps) {
  const theme = useTheme();
  const xsOnly: boolean = useMediaQuery(theme.breakpoints.only("xs"));
  const customWidth = {
    "&.MuiSvgIcon-root": {
      width: "38px",
      height: "38px",
      color: Colors.RED,
    },
  };

  return (
    <div>
      <Grid
        container
        mt={xsOnly ? "0.5rem" : ""}
        borderRadius="1rem"
        p="0.5rem"
        bgcolor="white.main"
        boxShadow={"2"}
      >
        {loading ? (
          <CashFlowSkeleton />
        ) : (
          <>
            {" "}
            <Grid xs={3}>
              <Grid
                height="100%"
                container
                alignContent="center"
                justifyContent="center"
              >
                <ArrowCircleLeftOutlinedIcon sx={customWidth} />
              </Grid>
            </Grid>
            <Grid xs={9}>
              <Typography variant="h6" fontWeight="bold">
                Expenses
              </Typography>
              <Typography lineHeight="unset" variant="h4" color="red.main">
                {new Intl.NumberFormat("en-US", {
                  style: "currency",
                  currency: "USD",
                }).format(expenses ? expenses : 0)}
              </Typography>
            </Grid>
          </>
        )}
      </Grid>
    </div>
  );
}
