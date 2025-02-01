import {Typography} from "@mui/material";
import {useNavigate} from "@tanstack/react-router";

type LogoProps = {
    variant?: string
}

export function Logo({variant = "h4"}: LogoProps) {
    const navigate = useNavigate()

    function onLogoClicked() {
        navigate({
            to: "/",
        }).then(() => {
        }).catch((e) => console.log("[money] Couldn't navigate to home from logo", e))
    }

    return (
        <>
            <div className="flex items-center cursor-pointer" onClick={() => onLogoClicked()}>
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