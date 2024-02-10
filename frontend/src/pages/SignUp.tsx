import {Box, Link, TextField, Typography} from "@mui/material";
import {Button} from "../components";
import {Colors} from "../assets";

export function SignUp() {
    return (
        <Box>
            {/*Title*/}
            <div className={"h-[12rem] flex items-center"}>
                <div>
                    <div className="flex items-center justify-center">
                        <img className="w-1/6"
                             src="https://money-static-files.s3.amazonaws.com/images/dollar.png"
                             alt="dollar_logo"/>
                        <Typography color={"darkGreen.main"} variant={"h2"} ml="5px">
                            Money
                        </Typography>
                    </div>
                    <div className={"flex justify-center"}>
                        <Typography variant={"h6"} color={"darkGreen.main"}>
                            Finance tracker
                        </Typography>
                    </div>
                </div>
            </div>

            {/*Input fields*/}
            <div className={"w-11/12 m-auto"}>
                <Typography textAlign={"center"} variant={"h4"}>
                    Sign up
                </Typography>

                <TextField margin={"normal"} fullWidth={true} label={"Full name"} variant={"outlined"}/>
                <TextField margin={"normal"} fullWidth={true} type={"email"} label={"Email"} variant={"outlined"}/>
                <TextField margin={"normal"} fullWidth={true} type={"password"} label={"Password"}
                           variant={"outlined"}/>
            </div>

            {/*Button*/}
            <div className={"w-11/12 m-auto pt-10"}>
                <Button variant={"contained"} fullWidth={true}>
                    Sign up
                </Button>
                <Typography textAlign={"center"}>
                    Already signed up? <Link color={Colors.BLUE_DARK} target={"_blank"} href={"/login"}>Login</Link>
                </Typography>
            </div>
        </Box>
    )
}