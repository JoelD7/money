import {Box, FormControl, InputLabel, MenuItem, Select, SelectChangeEvent} from "@mui/material";
import {v4 as uuidv4} from 'uuid';
import {useState} from "react";

export type ChipSelectOption = {
    label: string;
    color: string;
}

type ChipSelectProps = {
    options: ChipSelectOption[];
    label: string;
}

export function ChipSelect({options, label}: ChipSelectProps) {
    const labelId: string = uuidv4();
    const [selected, setSelected] = useState<string[]>([]);
    const colorMap: Map<string, string> = buildColorMap();

    function onSelectedChange(event: SelectChangeEvent<typeof selected>) {
        const {target: {value}} = event;
        setSelected(
            typeof value === 'string' ? value.split(' ') : value,
        );
    }

    function buildColorMap(): Map<string, string> {
        const colorMap = new Map<string, string>();
        options.forEach((option) => {
            colorMap.set(option.label, option.color);
        });

        return colorMap;
    }

    function getOptionColor(value: string): string {
        return colorMap.get(value) || "gray.main";
    }

    return (
        <>
            <FormControl fullWidth>
                <InputLabel id={labelId}>{label}</InputLabel>
                <Select
                    labelId={labelId}
                    id={label}
                    label={label}
                    value={selected}
                    onChange={onSelectedChange}
                    multiple
                    renderValue={(selected) => (
                        // This is how items will appear on the select input
                        <Box sx={{display: 'flex', flexWrap: 'wrap', gap: 0.5}}>
                            {selected.map((value) => (
                                <Box className="p-1 w-fit text-sm rounded-xl" style={{color: "white"}}
                                     sx={{backgroundColor: getOptionColor(value)}}>
                                    {value}
                                </Box>
                            ))}
                        </Box>
                    )}
                >
                    {
                        // This is how items will appear on the menu
                        options.map((option) => {
                            return (
                                <MenuItem id={option.label} value={option.label}>
                                    <Box className="p-1 w-fit text-sm rounded-xl" style={{color: "white"}}
                                         sx={{backgroundColor: option.color}}>
                                        {option.label}
                                    </Box>
                                </MenuItem>
                            );
                        })

                    }
                </Select>
            </FormControl>
        </>
    );
}