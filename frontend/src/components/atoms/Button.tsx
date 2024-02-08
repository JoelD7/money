import {Button as MuiButton, ButtonProps, SxProps, Theme} from "@mui/material";
import {ReactNode} from "react";

type CustomButtonProps = {
    children: ReactNode
} & ButtonProps

export function Button(props: CustomButtonProps) {
    const {sx, variant, ...other} = props
    let styles: SxProps<Theme> | undefined = {
        textTransform: "capitalize",
        borderRadius: "1rem",
        ...sx
    }

    if (variant === "outlined") {
        styles = {
            ...styles,
            '&.MuiButton-root': {
                backgroundColor: "#ffffff",
            },
        }
    }

    return (
        <>
            <MuiButton
                sx={styles}
                /*@ts-ignore*/
                variant={variant}
                {...other}>
                {props.children}
            </MuiButton>
        </>
    );
}