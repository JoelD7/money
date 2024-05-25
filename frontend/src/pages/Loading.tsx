import { CircularProgress, Typography } from "@mui/material";

export function Loading() {
  return (
    <div className="flex w-full h-screen justify-center items-center">
      <div>
        <div className={"m-auto w-fit"}>
          <CircularProgress size={"7rem"} />
        </div>
          <Typography variant={"h5"} textAlign={"center"}>💸 Preparing your data 💸</Typography>
        <Typography>Fetching your income, expenses and savings...</Typography>
      </div>
    </div>
  );
}
