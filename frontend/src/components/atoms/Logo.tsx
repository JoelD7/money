import {Typography} from "@mui/material";

type LogoProps = {
    variant?: string
}

export function Logo({variant = "h4"}: LogoProps) {

    return (
        <>
            <div className="flex items-center">
                <img className="w-14"
                     src="https://money-static-files.s3.amazonaws.com/images/dollar.png"
                     alt="dollar_logo"/>
                {/*// @ts-ignore*/}
                <Typography variant={variant} ml="5px">
                    Money
                </Typography>
            </div>
        </>
    );
}