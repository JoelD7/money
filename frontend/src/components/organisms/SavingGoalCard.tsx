import { CircularProgress, Typography } from "@mui/material";
import { Colors } from "../../assets";

type SavingGoalCardProps = {
  goal: number;
  progress: number;
};

export function SavingGoalCard({ goal, progress }: SavingGoalCardProps) {
  const size = "8rem";
  //This is an arbitrary value that works. Any value sufficiently larger or smaller breaks the animation
  const maxStrokeOffset = 289;
  const progressPercentage = (progress / goal) * 100;

  return (
    <div className={"rounded-md bg-white-100 shadow-md mt-4"}>
      {/* Progress circle */}
      <div className={"flex items-center justify-center p-8"}>
        <CircularProgress
          variant={"determinate"}
          value={100}
          size={size}
          sx={{
            color: Colors.GRAY,
            position: "absolute",
          }}
        />

        <CircularProgress
          variant={"determinate"}
          value={progressPercentage}
          size={size}
          sx={{
            [`& .MuiCircularProgress-circle`]: {
              strokeLinecap: "round",
              animation: "progress-grow 2s ease-out forwards",
              "@keyframes progress-grow": {
                from: {
                  strokeDashoffset: `${maxStrokeOffset}%`,
                },
                to: {
                  strokeDashoffset: `${maxStrokeOffset - (progressPercentage * maxStrokeOffset) / 100}%`,
                },
              },
            },
            color: Colors.GREEN_DARK,
          }}
        />

        <Typography
          variant={"h6"}
          color={"darkGreen.main"}
          sx={{
            position: "absolute",
          }}
        >
          {`${progressPercentage.toFixed(2)}%`}
        </Typography>
      </div>
    </div>
  );
}
