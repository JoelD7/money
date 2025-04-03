import { Logo } from "../components";
import { Typography } from "@mui/material";
import { Colors } from "../assets";
import { ReactNode } from "react";

type NotFoundProps = {
  title: string;
  subtitle: string;
  body?: string;
  /**
   *  @description The children nodes are treated as buttons that provided choices to the user on what to do next, like
   *  reloading the page or going back to the previous page
   *  */
  children?: ReactNode;
};

export function NotFound({ title, subtitle, body, children }: NotFoundProps) {
  return (
    <div className={"flex items-center justify-around h-lvh"}>
      <div className={"p-4 flex flex-col md:flex-row-reverse"}>
        {/*Image*/}
        <div className={"w-full md:w-2/3"}>
          <div className={"flex justify-center"}>
            <img
              className={"w-full max-w-lg md:max-w-5xl"}
              src="https://money-static-files.s3.amazonaws.com/images/page-error.jpg"
              alt="page error"
            />
          </div>
        </div>

        {/*Text*/}
        <div className={"w-full md:w-1/3 flex items-center justify-end"}>
          <div>
            <div className="p-4">
              <div className={"pb-5"}>
                <Logo />
              </div>
              <Typography color={"darkGreen.main"} variant={"h3"}>
                {title}
              </Typography>
              <Typography variant={"h5"} color={"gray.darker"}>
                {subtitle}
              </Typography>
              {body && (
                <div className={"pt-4 md:max-w-sm"}>
                  <p style={{ color: Colors.GRAY_DARK, fontSize: "20px" }}>{body}</p>
                </div>
              )}
              {children}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
