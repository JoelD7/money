type TextareaProps = {
    name: string;
    label?: string;
};

export function Textarea({label = "Label", name}: TextareaProps) {
    return (
        <>
            <label htmlFor={`${name}`} className="text-gray-200 block">{label}</label>
            <textarea name={name} id={name} cols="30" rows="5" placeholder="Text"
                      className="border-2 border-gray-100 rounded-lg p-2"></textarea>
        </>
    );
}