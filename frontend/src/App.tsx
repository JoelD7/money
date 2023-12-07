import "tailwindcss/tailwind.css";
import {Tag} from "./components/atoms/Tag.tsx";
import {Textarea} from "./components/atoms/Textarea.tsx";
import {SelectCustom} from "./components/atoms/SelectCustom.tsx";
import {createTheme, ThemeProvider} from "@mui/material";

function App() {
    const values: string[] = ["Value 1", "Value 2", "Very long value"];
    const theme = createTheme({
        typography: {
            fontFamily: [
                'Outfit',
                'sans-serif',
            ].join(','),
        }
    })
    return (
        <>
            <ThemeProvider theme={theme}>
                <div style={{width: "1290px"}}>
                    <div style={{margin: "auto", width: "fit-content"}}>
                        <SelectCustom name="select" label="Currency" values={values}></SelectCustom>
                    </div>
                </div>
            </ThemeProvider>

        </>
    )
}

export default App
