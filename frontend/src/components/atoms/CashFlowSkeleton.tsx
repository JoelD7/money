import Grid from "@mui/material/Unstable_Grid2";
import { Skeleton } from "@mui/material";

// Loading skeleton to be displayed while the balance, expenses and income data is being fetched
export function CashFlowSkeleton() {
  return (
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
  );
}
