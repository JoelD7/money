import './App.css'
import {Button, ButtonColor} from "./atoms/Button.tsx";
import "tailwindcss/tailwind.css";

function App() {
    return (
        <>
            <Button text={"Button"} color={ButtonColor.White}></Button>
        </>
    )
}

export default App
