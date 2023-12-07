import "tailwindcss/tailwind.css";
import {Tag} from "./components/atoms/Tag.tsx";
import {Textarea} from "./components/atoms/Textarea.tsx";

function App() {
    return (
        <>
            <div style={{width: "1290px"}}>
                <div style={{margin: "auto", width: "fit-content"}}>
                    <Textarea name="textarea"></Textarea>
                </div>
            </div>
        </>
    )
}

export default App
