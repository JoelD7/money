import {IconButton} from "@mui/material";
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
            <nav className="flex flex-col items-center h-screen border-black border-2 w-28">
                <img className="mb-10 w-14"
                     src="https://money-static-files.s3.amazonaws.com/images/dollar.png"
                     alt="dollar_logo"/>

                <IconButton>
                    <HomeIcon sx={customWidth}/>
                </IconButton>
                <IconButton size="medium">
                    <AccessTimeFilledIcon sx={customWidth}/>
                </IconButton>
                <IconButton size="medium">
                    <NotificationImportantIcon sx={customWidth}/>
                </IconButton>
                <IconButton size="medium">
                    <LabelIcon sx={customWidth}/>
                </IconButton>
                <IconButton size="medium">
                    <AutoStoriesIcon sx={customWidth}/>
                </IconButton>
                <hr style={{color: "gray.main"}}/>
                <IconButton size="medium">
                    <SettingsIcon sx={customWidth}/>
                </IconButton>
                <IconButton sx={{
                    marginTop: "auto",
                    marginBottom: "20px",
                }} size="medium">
                    <LogoutIcon sx={customWidth}/>
                </IconButton>
            </nav>

        </>
    );
}