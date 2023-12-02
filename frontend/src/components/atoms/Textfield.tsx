type TextfieldProps = {
    name: string;
    label: string;
};


export function Textfield({label = "Label", name}: TextfieldProps) {
    return (
        <>
            <div className={"block"}>
                <label htmlFor={`${name}`} className="text-gray-200">{label}</label>
                <input type="text" id={name} name={name} className="border-2  border-gray-100 rounded-md p-2 m-2"
                       placeholder="Text"/>
            </div>
        </>
    );
}