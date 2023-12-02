type ButtonProps = {
    color: ButtonColor;
    text: string;
};

export enum ButtonColor {
    Red = "RED",
    Orange = "ORANGE",
    Green = "GREEN",
    Gray = "GRAY",
    White = "WHITE",
}

export function Button({color = ButtonColor.White, text = "Button"}: ButtonProps) {
    return (
        <>
            <button className={`${getButtonColorClass(color)} shadow-lg`}>
                {text}
            </button>
        </>
    );
}

function getButtonColorClass(color: ButtonColor): string {
    switch (color) {
        case ButtonColor.Red:
            return "bg-red-100 hover:bg-red-200 active:bg-red-300";
        case ButtonColor.Orange:
            return "bg-orange-100 hover:bg-orange-200 active:bg-orange-300";
        case ButtonColor.Green:
            return "bg-green-100 hover:bg-green-200 active:bg-green-300";
        case ButtonColor.Gray:
            return "bg-gray-100 hover:bg-gray-200 active:bg-gray-300";
        case ButtonColor.White:
            return "border-green-300 border-2 bg-white-100 hover:bg-white-200 active:bg-white-300";
        default:
            return "bg-green-100 hover:bg-green-200 active:bg-green-300";
    }
}