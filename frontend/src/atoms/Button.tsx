type ButtonProps = {
    color: string;
    text: string;
};

export function Button({color, text}: ButtonProps) {
    return (
        <>
            <button>
                {text}
            </button>
        </>
    );
}