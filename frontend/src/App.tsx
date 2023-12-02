import './App.css'
import {Button, ButtonColor} from "./components/atoms/Button.tsx";
import "tailwindcss/tailwind.css";

function App() {
    return (
        <>
            <Button text={"Button"} color={ButtonColor.White}></Button>
        </>
    )
}

export default App
