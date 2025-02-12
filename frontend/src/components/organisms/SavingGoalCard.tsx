import { CircularProgress, circularProgressClasses, Typography } from "@mui/material";
import { Colors } from "../../assets";

type SavingGoalCardProps = {
  goal: number;
  progress: number;
};

export function SavingGoalCard({ goal, progress }: SavingGoalCardProps) {
  const size = "8rem";
  const maxStrokeOffset = 300
  const percentageFormat = new Intl.NumberFormat("en-US", {
    style: "decimal",
    minimumFractionDigits: 0,
  });

  function getProgressPercentage(): number {
    return (progress / goal) * 100;
  }

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
          value={getProgressPercentage()}
          size={size}
          sx={{
              [`& .${circularProgressClasses.circle}`]: {
                  strokeLinecap: 'round',
                  animation: 'progress-grow 2s ease-out forwards',
                  '@keyframes progress-grow': {
                      from: {
                          strokeDashoffset: `${maxStrokeOffset}%`,
                      },
                      to: {
                          strokeDashoffset: `${maxStrokeOffset - (progress * maxStrokeOffset / 100)}%`,
                      }
                  }
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
          {`${percentageFormat.format(getProgressPercentage())}%`}
        </Typography>
      </div>
    </div>
  );
}
