import {Button, Divider, Drawer, IconButton, useMediaQuery, useTheme} from "@mui/material";
import HomeIcon from '@mui/icons-material/Home';
import MenuIcon from '@mui/icons-material/Menu';
import AccessTimeFilledIcon from '@mui/icons-material/AccessTimeFilled';
import NotificationImportantIcon from '@mui/icons-material/NotificationImportant';
import LabelIcon from '@mui/icons-material/Label';
import AutoStoriesIcon from '@mui/icons-material/AutoStories';
import SettingsIcon from '@mui/icons-material/Settings';
import LogoutIcon from '@mui/icons-material/Logout';
import {useState, ReactNode} from "react";
import {Logo} from "../atoms";

type NavbarProps = {
    children?: ReactNode
}

export function Navbar({children}: NavbarProps) {
    const customWidth = {
        '&.MuiSvgIcon-root': {
            width: "28px",
            height: "28px",
            fill: "#024511"
        },
    }

    const buttonStyle = {
        margin: "5px 0px",
        color: "gray.dark",
        textTransform: "capitalize",
        '&.MuiButton-root': {
            justifyContent: "flex-start",
            width: "100%"
        },
    }

    const [open, setOpen] = useState<boolean>(false)
    const theme = useTheme();
    const mdUp: boolean = useMediaQuery(theme.breakpoints.up('md'));

    return (
        <>
            <div className={mdUp ? "hidden" : "flex p-4 flex-row justify-items-center"}>
                {
                    children ? children : <Logo/>
                }

                <div className="ml-auto mr-3">
                    <IconButton title="Home" sx={{margin: "5px 0px"}}>
                        <HomeIcon sx={customWidth}/>
                    </IconButton>

                    <IconButton sx={{marginLeft: "15px"}} onClick={() => setOpen(true)}>
                        <MenuIcon sx={customWidth}/>
                    </IconButton>
                </div>
            </div>

            {/*TODO: add back button when route isn't Home*/}
            {/*Title and go back*/}
            {/*<Grid xs={12}>*/}
            {/*    <IconButton>*/}
            {/*        /!*@ts-ignore*!/*/}
            {/*        <ArrowCircleLeftIcon sx={backButtonStyle} color={"darkGreen"}/>*/}
            {/*    </IconButton>*/}
            {/*</Grid>*/}

            <Drawer anchor="right" open={open} onClose={() => setOpen(false)}>
                <nav hidden={mdUp} style={{backgroundColor: "white"}}
                     className="flex flex-col h-screen w-44">
                    <div className="flex items-center p-4 justify-center w-full">
                        <Logo variant="h5"/>
                    </div>

                    <div className="pl-3">
                        <Button sx={buttonStyle} startIcon={<HomeIcon sx={customWidth}/>}>
                            Home
                        </Button>

                        <Button sx={buttonStyle} startIcon={<AccessTimeFilledIcon sx={customWidth}/>}>
                            History
                        </Button>

                        <Button sx={buttonStyle} startIcon={<NotificationImportantIcon sx={customWidth}/>}>
                            Notifications
                        </Button>

                        <Button sx={buttonStyle} startIcon={<LabelIcon sx={customWidth}/>}>
                            Categories
                        </Button>

                        <Button sx={buttonStyle} startIcon={<AutoStoriesIcon sx={customWidth}/>}>
                            Savings
                        </Button>
                    </div>

                    <Divider sx={{width: "60%", margin: "20px auto"}}/>

                    <div className="pl-3 h-full">
                        <Button sx={{...buttonStyle, margin: "0px"}}
                                startIcon={<SettingsIcon sx={customWidth}/>}>
                            Settings
                        </Button>

                        <Button sx={{...buttonStyle, marginTop: "auto", marginBottom: "20px"}}
                                startIcon={<LogoutIcon sx={customWidth}/>}>
                            Logout
                        </Button>
                    </div>
                </nav>
            </Drawer>

            <nav style={{backgroundColor: "white"}}
                 className={mdUp ? "flex flex-col h-screen w-44 fixed" : "hidden"}>
                <div className="flex items-center p-4 justify-center w-full">
                    <Logo variant="h5"/>
                </div>

                <div className="pl-3">
                    <Button sx={buttonStyle} startIcon={<HomeIcon sx={customWidth}/>}>
                        Home
                    </Button>

                    <Button sx={buttonStyle} startIcon={<AccessTimeFilledIcon sx={customWidth}/>}>
                        History
                    </Button>

                    <Button sx={buttonStyle} startIcon={<NotificationImportantIcon sx={customWidth}/>}>
                        Notifications
                    </Button>

                    <Button sx={buttonStyle} startIcon={<LabelIcon sx={customWidth}/>}>
                        Categories
                    </Button>

                    <Button sx={buttonStyle} startIcon={<AutoStoriesIcon sx={customWidth}/>}>
                        Savings
                    </Button>
                </div>

                <Divider sx={{width: "60%", margin: "20px auto"}}/>

                <div className="pl-3 h-full">
                    <Button sx={{...buttonStyle, margin: "0px"}}
                            startIcon={<SettingsIcon sx={customWidth}/>}>
                        Settings
                    </Button>

                    <Button sx={{...buttonStyle, marginTop: "auto", marginBottom: "20px"}}
                            startIcon={<LogoutIcon sx={customWidth}/>}>
                        Logout
                    </Button>
                </div>
            </nav>
        </>
    );
}