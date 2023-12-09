import {IconButton} from "@mui/material";
import HomeIcon from '@mui/icons-material/Home';
import AccessTimeFilledIcon from '@mui/icons-material/AccessTimeFilled';
import NotificationImportantIcon from '@mui/icons-material/NotificationImportant';
import LabelIcon from '@mui/icons-material/Label';
import AutoStoriesIcon from '@mui/icons-material/AutoStories';
import SettingsIcon from '@mui/icons-material/Settings';
import LogoutIcon from '@mui/icons-material/Logout';

export function Navbar() {
    return (
        <>
            <nav className="flex flex-col items-center border-black border-2 w-28">
                <img className="mb-10 w-14"
                     src="https://money-static-files.s3.amazonaws.com/images/dollar.png"
                     alt="dollar_logo"/>

                <IconButton>

                    <HomeIcon sx={{
                        '& .MuiSvgIcon-fontSizeMedium.MuiSvgIcon-root': {
                            width: "52px",
                        },
                    }}/>
                </IconButton>
                <IconButton size="medium">
                    <AccessTimeFilledIcon fontSize="medium"/>
                </IconButton>
                <IconButton size="medium">
                    <NotificationImportantIcon fontSize="medium"/>
                </IconButton>
                <IconButton size="medium">
                    <LabelIcon fontSize="medium"/>
                </IconButton>
                <IconButton size="medium">
                    <AutoStoriesIcon fontSize="medium"/>
                </IconButton>
                <IconButton size="medium">
                    <SettingsIcon fontSize="medium"/>
                </IconButton>
                <IconButton size="medium">
                    <LogoutIcon fontSize="medium"/>
                </IconButton>
            </nav>

        </>
    );
}