import "tailwindcss/tailwind.css";
import {Tag} from "./components/atoms/Tag.tsx";

function App() {
    return (
        <>
            <div style={{width: "1290px"}}>
                <div style={{margin: "auto", width: "fit-content"}}>
                    <Tag color="blue-100" label="Label"/>
                </div>
            </div>
        </>
    )
}

export default App
