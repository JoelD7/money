import Grid from "@mui/material/Unstable_Grid2";
import ArrowCircleRightOutlinedIcon from "@mui/icons-material/ArrowCircleRightOutlined";
import { Typography, useMediaQuery, useTheme } from "@mui/material";
import { Colors } from "../../assets";
import { CashFlowSkeleton } from "../atoms";

type IncomeCardProps = {
  income?: number;
  loading: boolean;
};

export function IncomeCard({ income, loading }: IncomeCardProps) {
  const theme = useTheme();
  const xsOnly: boolean = useMediaQuery(theme.breakpoints.only("xs"));
  const customWidth = {
    "&.MuiSvgIcon-root": {
      width: "38px",
      height: "38px",
      color: Colors.GREEN,
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
            <Grid xs={3}>
              <Grid
                height="100%"
                container
                alignContent="center"
                justifyContent="center"
              >
                <ArrowCircleRightOutlinedIcon sx={customWidth} />
              </Grid>
            </Grid>

            <Grid xs={9}>
              <Typography variant="h6" fontWeight="bold">
                Income
              </Typography>
              <Typography lineHeight="unset" variant="h4" color="primary">
                {new Intl.NumberFormat("en-US", {
                  style: "currency",
                  currency: "USD",
                }).format(income ? income : 0)}
              </Typography>
            </Grid>
          </>
        )}
      </Grid>
    </div>
  );
}
