import './App.css'
import {Button, ButtonColor} from "./components/atoms/Button.tsx";
import "tailwindcss/tailwind.css";
import {Textfield} from "./components/atoms/Textfield.tsx";

function App() {
    return (
        <>
            <Textfield name={"textfield"} label={"Label"}></Textfield>
        </>
    )
}

export default App
