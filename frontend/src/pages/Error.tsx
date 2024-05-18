import { Button, Logo } from "../components";
import { Typography } from "@mui/material";
import { Colors } from "../assets";

export function Error() {
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
                  <Logo/>
                </div>
                <Typography color={"darkGreen.main"} variant={"h3"}>
                  Whoops...
                </Typography>
                <Typography variant={"h5"} color={"gray.darker"}>
                  There has been an error...
                </Typography>
                <div className={"pt-4 md:max-w-sm"}>
                  <p style={{color: Colors.GRAY_DARK, fontSize: "20px"}}>
                    Our servers seem to be having some issues. Please try again in a
                    few minutes.
                  </p>
                </div>

                <Button
                    variant={"contained"}
                    sx={{
                      marginTop: "10px",
                      fontSize: "18px",
                    }}
                    onClick={() => window.location.reload()}
                >
                  Reload page
                </Button>
              </div>
            </div>
          </div>
        </div>
      </div>
  );
}
