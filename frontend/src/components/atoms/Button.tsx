import {Button as MuiButton, ButtonProps} from "@mui/material";
import {ReactNode} from "react";

type CustomButtonProps = {
    children: ReactNode
} & ButtonProps

export function Button(props: CustomButtonProps) {
    const {sx, ...other} = props
    return (
        <>
            <MuiButton sx={{textTransform: "capitalize", borderRadius: "1rem", ...sx}} {...other}>
                {props.children}
            </MuiButton>
        </>
    );
}