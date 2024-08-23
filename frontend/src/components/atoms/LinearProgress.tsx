import { LinearProgress as MuiLinearProgress } from "@mui/material";

type LinearProgressProps = {
  loading?: boolean;
};

export function LinearProgress({ loading }: LinearProgressProps) {
  return (

    <div className={`${loading ? "block" : "hidden"}`}>
        <MuiLinearProgress />
    </div>
  );
}
