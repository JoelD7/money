import {Button, Divider, Drawer, IconButton} from "@mui/material";
import HomeIcon from '@mui/icons-material/Home';
import MenuIcon from '@mui/icons-material/Menu';
import AccessTimeFilledIcon from '@mui/icons-material/AccessTimeFilled';
import NotificationImportantIcon from '@mui/icons-material/NotificationImportant';
import LabelIcon from '@mui/icons-material/Label';
import AutoStoriesIcon from '@mui/icons-material/AutoStories';
import SettingsIcon from '@mui/icons-material/Settings';
import LogoutIcon from '@mui/icons-material/Logout';
import {useState} from "react";

export function Navbar() {
    const customWidth = {
        '&.MuiSvgIcon-root': {
            width: "28px",
            height: "28px",
            fill: "#024511"
        },
    }

    const buttonStyle = {
        margin: "5px 0px", color: "gray.dark", textTransform: "capitalize"
    }

    const [open, setOpen] = useState<boolean>(false)

    return (
        <>
            <div className="flex p-2 flex-row justify-items-center">
                <img className="w-14"
                     src="https://money-static-files.s3.amazonaws.com/images/dollar.png"
                     alt="dollar_logo"/>

                <div className="ml-auto">
                    <IconButton title="Home" sx={{margin: "5px 0px"}}>
                        <HomeIcon sx={customWidth}/>
                    </IconButton>

                    <IconButton sx={{margin: "5px 0px"}}>
                        <AccessTimeFilledIcon sx={customWidth}/>
                    </IconButton>

                    <IconButton sx={{marginLeft: "15px"}} onClick={() => setOpen(true)}>
                        <MenuIcon sx={customWidth}/>
                    </IconButton>
                </div>
            </div>

            <Drawer anchor="left" open={open} onClose={() => setOpen(false)}>
                <nav style={{backgroundColor: "white"}}
                     className="flex flex-col items-start h-screen w-36">

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
                        <Button sx={{color: "gray.dark", textTransform: "capitalize"}}
                                startIcon={<SettingsIcon sx={customWidth}/>}>
                            Settings
                        </Button>

                        <Button sx={{
                            marginTop: "auto",
                            marginBottom: "20px",
                            color: "gray.dark",
                            textTransform: "capitalize"
                        }} startIcon={<LogoutIcon sx={customWidth}/>}>
                            Logout
                        </Button>
                    </div>
                </nav>
            </Drawer>

            {/*Desktop version*/}
            {/*<nav style={{backgroundColor: "white"}}*/}
            {/*//      className="flex flex-col items-center h-screen w-28">*/}
            {/*//*/}
            {/*//     <Tooltip title="Home" placement="right">*/}
            {/*//         <IconButton sx={{margin: "5px 0px"}}>*/}
            {/*//             <HomeIcon sx={customWidth}/>*/}
            {/*//         </IconButton>*/}
            {/*//     </Tooltip>*/}
            {/*//*/}
            {/*//     <Tooltip title="History" placement="right">*/}
            {/*//         <IconButton sx={{margin: "5px 0px"}}>*/}
            {/*//             <AccessTimeFilledIcon sx={customWidth}/>*/}
            {/*//         </IconButton>*/}
            {/*//     </Tooltip>*/}
            {/*//*/}
            {/*//     <Tooltip title="Notifications" placement="right">*/}
            {/*//         <IconButton sx={{margin: "5px 0px"}}>*/}
            {/*//             <NotificationImportantIcon sx={customWidth}/>*/}
            {/*//         </IconButton>*/}
            {/*//     </Tooltip>*/}
            {/*//*/}
            {/*//     <Tooltip title="Categories" placement="right">*/}
            {/*//         <IconButton sx={{margin: "5px 0px"}}>*/}
            {/*//             <LabelIcon sx={customWidth}/>*/}
            {/*//         </IconButton>*/}
            {/*//     </Tooltip>*/}
            {/*//*/}
            {/*//     <Tooltip title="Savings" placement="right">*/}
            {/*//         <IconButton sx={{margin: "5px 0px"}}>*/}
            {/*//             <AutoStoriesIcon sx={customWidth}/>*/}
            {/*//         </IconButton>*/}
            {/*//     </Tooltip>*/}
            {/*//*/}
            {/*//     <Divider sx={{width: "60%", margin: "20px 0px"}}/>*/}
            {/*//*/}
            {/*//     <Tooltip title="Settings" placement="right">*/}
            {/*//         <IconButton>*/}
            {/*//             <SettingsIcon sx={customWidth}/>*/}
            {/*//         </IconButton>*/}
            {/*//     </Tooltip>*/}
            {/*//*/}
            {/*//     <Tooltip title="Logout" placement="right">*/}
            {/*//         <IconButton sx={{*/}
            {/*//             marginTop: "auto",*/}
            {/*//             marginBottom: "20px",*/}
            {/*//         }}>*/}
            {/*//             <LogoutIcon sx={customWidth}/>*/}
            {/*//         </IconButton>*/}
            {/*//     </Tooltip>*/}
            {/*// </nav>*/}
        </>
    );
}