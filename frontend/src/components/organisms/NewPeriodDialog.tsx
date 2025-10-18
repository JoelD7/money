import {Dialog} from "../molecules";
import Grid from "@mui/material/Unstable_Grid2";
import {Box, Divider, TextField, Typography} from "@mui/material";
import {useState} from "react";
import dayjs, {Dayjs} from "dayjs";
import {DatePicker} from "@mui/x-date-pickers";

type NewPeriodDialogProps = {
    open: boolean;
    onClose: ()=> void;
}

export function NewPeriodDialog({open, onClose}: NewPeriodDialogProps){
    const [date, setDate] = useState<Dayjs | null>(dayjs());
    const[name, setName] = useState("")

    return (
        <Dialog open={open} onClose={onClose}>
            <Box component={"form"}>
                <Grid
                    container
                    spacing={2}
                    bgcolor={"white.main"}
                    borderRadius="1rem"
                    width={"500px"}
                    p="1.5rem"
                >
                    <Grid xs={12}>
                        <Typography variant={"h4"}>Create new period</Typography>
                        <Divider />
                    </Grid>

                    <Grid xs={12}>
                        <TextField
                            margin={"none"}
                            name={"name"}
                            value={name}
                            fullWidth={true}
                            type={"text"}
                            label={"Name"}
                            variant={"outlined"}
                            onChange={(e) => setName(e.target.value)}
                        />
                    </Grid>

                    <Grid xs={12}>
                        <DatePicker
                            label="Date"
                            sx={{ width: "100%" }}
                            value={date}
                            onChange={(newDate) => setDate(newDate)}
                        />
                    </Grid>
                </Grid>
            </Box>
        </Dialog>
    )
}