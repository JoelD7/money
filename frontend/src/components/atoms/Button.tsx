import {Button as MuiButton, ButtonProps} from "@mui/material";
import {ReactNode} from "react";

type CustomButtonProps = {
    children: ReactNode
} & ButtonProps

export function Button(props: CustomButtonProps) {
    return (
        <>
            <MuiButton sx={{textTransform: "capitalize", borderRadius: "1rem"}} {...props}>
                {props.children}
            </MuiButton>
        </>
    );
}