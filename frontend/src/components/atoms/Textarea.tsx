export function Textarea() {
    return (
        <>
            <label htmlFor={`${name}`} className="text-gray-200 block">{label}</label>
            <input type="text" id={name} name={name} className="border-2 border-gray-100 rounded-lg p-2"
                   placeholder="Text"/>
        </>
    );
}