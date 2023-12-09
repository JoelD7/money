import {Box, FormControl, InputLabel, MenuItem, Select} from "@mui/material";
import {v4 as uuidv4} from 'uuid';

type option = {
    label: string;
    color: string;
}

type ChipSelectProps = {
    options: option[];
    label: string;
}

export function ChipSelect({options, label}: ChipSelectProps) {
    const labelId: string = uuidv4();

    return (
        <>
            <FormControl fullWidth>
                <InputLabel id={labelId}>Age</InputLabel>
                <Select
                    labelId={labelId}
                    id={label}
                    label={label}
                    multiple
                >
                    {
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