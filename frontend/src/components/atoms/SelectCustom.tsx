import {MenuItem, Select} from "@mui/material";

type InputSelectProps = {
    name: string;
    label?: string;
    values: string[];
};

export function SelectCustom({label = "", name, values}: InputSelectProps) {
    return (
        <>
            <Select name={name} id={name} className="border-2 border-gray-100 rounded-lg p-2"
            >
                {values.map((value, index) => {
                    return (
                        <MenuItem inputMode="text" key={index} value={value}>{value}</MenuItem>
                    );
                })}
            </Select>
        </>
    );
}