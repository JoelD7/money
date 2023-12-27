import "tailwindcss/tailwind.css";
import {createTheme, ThemeProvider, Theme, Container} from "@mui/material";
import {Navbar} from "./components";

declare module '@mui/material/styles' {
    interface Palette {
        custom: Palette['primary'];
    }

    interface PaletteOptions {
        white?: PaletteOptions['primary'];
        red?: PaletteOptions['primary'];
        blue?: PaletteOptions['primary'];
        gray?: PaletteOptions['primary'];
    }

    interface PaletteColor {
        darker?: string;
    }

    interface SimplePaletteColorOptions {
        darker?: string;
    }
}

const theme: Theme = createTheme({
    typography: {
        fontFamily: [
            'Outfit',
            'sans-serif',
        ].join(','),
    },
    palette: {
        primary: {
            main: "#009821",
            darker: "#024511",
        },
        secondary: {
            main: "#FF8042",
        },
        white: {
            main: "#FFFFFF",
            dark: "#e6e6e6",
            darker: "#cccccc"
        },
        red: {
            main: '#D90707',
            dark: '#ad0101',
            darker: '#7a0101',
        },
        blue: {
            main: '#0088FE',
            dark: '#006dcc',
            darker: '#004d99',
        },
        gray: {
            main: '#D9D9D9',
            dark: '#6F6F6F',
            darker: '#4D4D4D',
        }
    },
})

function App() {
    
    return (
        <>
            <ThemeProvider theme={theme}>
                <Container maxWidth={false} sx={{backgroundColor: "#f1f1f1"}}>
                    <Navbar/>
                </Container>
            </ThemeProvider>

        </>
    )
}

export default App
