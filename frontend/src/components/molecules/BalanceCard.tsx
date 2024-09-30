import Grid from "@mui/material/Unstable_Grid2";
import RemoveCircleOutlineOutlinedIcon from "@mui/icons-material/RemoveCircleOutlineOutlined";
import { Typography } from "@mui/material";
import { CashFlowSkeleton } from "../atoms";

type BalanceCardProps = {
  remainder: number;
  loading: boolean;
};

export function BalanceCard({ remainder, loading }: BalanceCardProps) {
  const customWidth = {
    "&.MuiSvgIcon-root": {
      width: "38px",
      height: "38px",
      color: "black",
    },
  };

  return (
    <div>
      <Grid
        container
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
                <RemoveCircleOutlineOutlinedIcon sx={customWidth} />
              </Grid>
            </Grid>

            <Grid xs={9}>
              <Typography variant="h6" fontWeight="bold">
                Balance
              </Typography>
              <Typography lineHeight="unset" variant="h4" color="black">
                {new Intl.NumberFormat("en-US", {
                  style: "currency",
                  currency: "USD",
                }).format(remainder)}
              </Typography>
            </Grid>
          </>
        )}
      </Grid>
    </div>
  );
}
