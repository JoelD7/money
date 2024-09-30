import Grid from "@mui/material/Unstable_Grid2";
import ArrowCircleRightOutlinedIcon from "@mui/icons-material/ArrowCircleRightOutlined";
import { Skeleton, Typography, useMediaQuery, useTheme } from "@mui/material";
import { Colors } from "../../assets";

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
          <>
            <Grid xs={3}>
              <Grid
                height="100%"
                container
                alignContent="center"
                justifyContent="center"
              >
                <Skeleton variant={"circular"} width={40} height={40} />
              </Grid>
            </Grid>

            <Grid xs={9}>
              <Skeleton
                variant={"rectangular"}
                sx={{ marginBottom: "5px" }}
                width={282}
                height={31}
              />
              <Skeleton variant={"rectangular"} width={282} height={50} />
            </Grid>
          </>
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
