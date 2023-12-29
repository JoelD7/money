import "tailwindcss/tailwind.css";
import {createTheme, ThemeProvider, Theme, Container} from "@mui/material";
import {Home} from "./pages";
import shadows from "@mui/material/styles/shadows";

declare module '@mui/material/styles' {
    interface Palette {
        custom: Palette['primary'];
    }

    interface PaletteOptions {
        white?: PaletteOptions['primary'];
        red?: PaletteOptions['primary'];
        blue?: PaletteOptions['primary'];
        gray?: PaletteOptions['primary'];
        darkGreen?: PaletteOptions['primary'];
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
            // Use this color as it is the same as the "bg-zinc-100" Tailwind class
            main: '#F4F4F5',
            dark: '#6F6F6F',
            darker: '#4D4D4D',
        },
        darkGreen: {
            main: `#024511`,
        }
    },
})

function App() {

    return (
        <>
            <ThemeProvider theme={theme}>
                <Container maxWidth={false}>
                    <Home/>
                </Container>
            </ThemeProvider>

        </>
    )
}

export default App
