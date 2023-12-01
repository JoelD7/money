import {useState} from 'react'
import './App.css'
import {Button} from "./atoms/Button.tsx";
import "tailwindcss/tailwind.css";

function App() {
    const [count, setCount] = useState(0)

    return (
        <>
            <Button text={"Button"}></Button>
        </>
    )
}

export default App
