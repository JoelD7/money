import { Colors } from "../../assets";
import { CircularProgress as MuiCircularProgress, Typography } from "@mui/material";

type CircularProgressProps = {
  progress: number;
  target: number;
  // Size in rem units
  size: string;
  subtitle?: string;
};

export function CircularProgress({
  progress,
  target,
  size,
  subtitle,
}: CircularProgressProps) {
  const progressPercentage = (progress / target) * 100;
  //This is an arbitrary value that works. Any value sufficiently larger or smaller breaks the animation
  const maxStrokeOffset = 289;
  const hiddenStroke = (() => {
    const result = maxStrokeOffset - (progressPercentage * maxStrokeOffset) / 100;
    //Prevent the line that fills up the circle from overflowing
    return result < 0 ? 0 : result;
  })();

  return (
    <div className={"flex items-center justify-center p-8"}>
      <MuiCircularProgress
        variant={"determinate"}
        value={100}
        size={size}
        sx={{
          color: Colors.GRAY,
          position: "absolute",
        }}
      />

      <MuiCircularProgress
        variant={"determinate"}
        value={progressPercentage}
        size={size}
        sx={{
          position: "absolute",
          [`& .MuiCircularProgress-circle`]: {
            strokeLinecap: "round",
            animation: "progress-grow 2s ease-out forwards",
            "@keyframes progress-grow": {
              from: {
                strokeDashoffset: `${maxStrokeOffset}%`,
              },
              to: {
                strokeDashoffset: `${hiddenStroke}%`,
              },
            },
          },
          color: Colors.GREEN_DARK,
        }}
      />

      <div className={"flex justify-center flex-col"}>
        <Typography
          variant={"h6"}
          color={"darkGreen.main"}
          sx={{
            width: "100%",
          }}
        >
          {`${progressPercentage.toFixed(2)}%`}
        </Typography>

        {subtitle && (
          <h6 className={`text-sm leading-none text-center text-gray-200`}>{subtitle}</h6>
        )}
      </div>
    </div>
  );
}
