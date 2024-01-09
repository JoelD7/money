import "tailwindcss/tailwind.css";
import {createTheme, ThemeProvider, Theme, Container, useMediaQuery} from "@mui/material";
import {Home, PeriodDetail} from "./pages";
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
            contrastText: "#ffffff"
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
            light: '#a3a3a3',
        },
        darkGreen: {
            main: `#024511`,
        }
    },
})

function App() {
    const mdUp: boolean = useMediaQuery(theme.breakpoints.up('md'));
    const containerStyles = {
        backgroundColor: "#fafafa",
        width: "auto",
    }
    return (
        <>
            <ThemeProvider theme={theme}>
                <Navbar/>
                <Container
                    sx={mdUp ? {marginLeft: "11rem", ...containerStyles} : {...containerStyles}}
                    maxWidth={false}>
                    <PeriodDetail/>
                </Container>
            </ThemeProvider>

        </>
    )
}

export default App
