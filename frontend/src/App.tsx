import "tailwindcss/tailwind.css";
import {createTheme, ThemeProvider} from "@mui/material";

function App() {
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
                    <div className="ml-2" style={{width: "fit-content"}}>

                    </div>
                </div>
            </ThemeProvider>

        </>
    )
}

export default App
