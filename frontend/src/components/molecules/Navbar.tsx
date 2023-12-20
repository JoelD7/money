import {Divider, IconButton, Tooltip} from "@mui/material";
import HomeIcon from '@mui/icons-material/Home';
import AccessTimeFilledIcon from '@mui/icons-material/AccessTimeFilled';
import NotificationImportantIcon from '@mui/icons-material/NotificationImportant';
import LabelIcon from '@mui/icons-material/Label';
import AutoStoriesIcon from '@mui/icons-material/AutoStories';
import SettingsIcon from '@mui/icons-material/Settings';
import LogoutIcon from '@mui/icons-material/Logout';

export function Navbar() {
    const customWidth = {
        '&.MuiSvgIcon-root': {
            width: "28px",
            height: "28px",
        },
    }
    return (
        <>
            <nav className="flex flex-col items-center h-screen  w-28">
                <img className="mb-10 w-14"
                     src="https://money-static-files.s3.amazonaws.com/images/dollar.png"
                     alt="dollar_logo"/>

                <Tooltip title="Home" placement="right">
                    <IconButton sx={{margin: "5px 0px"}}>
                        <HomeIcon sx={customWidth}/>
                    </IconButton>
                </Tooltip>

                <Tooltip title="History" placement="right">
                    <IconButton sx={{margin: "5px 0px"}}>
                        <AccessTimeFilledIcon sx={customWidth}/>
                    </IconButton>
                </Tooltip>

                <Tooltip title="Notifications" placement="right">
                    <IconButton sx={{margin: "5px 0px"}}>
                        <NotificationImportantIcon sx={customWidth}/>
                    </IconButton>
                </Tooltip>

                <Tooltip title="Categories" placement="right">
                    <IconButton sx={{margin: "5px 0px"}}>
                        <LabelIcon sx={customWidth}/>
                    </IconButton>
                </Tooltip>

                <Tooltip title="Savings" placement="right">
                    <IconButton sx={{margin: "5px 0px"}}>
                        <AutoStoriesIcon sx={customWidth}/>
                    </IconButton>
                </Tooltip>

                <Divider sx={{width: "60%", margin: "20px 0px"}}/>

                <Tooltip title="Settings" placement="right">
                    <IconButton>
                        <SettingsIcon sx={customWidth}/>
                    </IconButton>
                </Tooltip>

                <Tooltip title="Logout" placement="right">
                    <IconButton sx={{
                        marginTop: "auto",
                        marginBottom: "20px",
                    }}>
                        <LogoutIcon sx={customWidth}/>
                    </IconButton>
                </Tooltip>
            </nav>

        </>
    );
}