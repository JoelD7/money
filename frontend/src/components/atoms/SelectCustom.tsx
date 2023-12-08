import {InputLabel, MenuItem, Select} from "@mui/material";
import {CSSProperties} from "react";

type InputSelectProps = {
    name: string;
    label?: string;
    values: string[];
};

export function SelectCustom({label = "", name, values}: InputSelectProps) {
    return (
        <>
            {
                label !== "" &&
                <InputLabel id="demo-simple-select-label">{label}</InputLabel>
            }
            <Select name={name} id={name} style={{width: "50px"}} labelId="demo-simple-select-label">
                {values.map((value, index) => {
                    return (
                        <MenuItem inputMode="text" key={index} value={value}>{value}</MenuItem>
                    );
                })}
            </Select>
        </>
    );
}