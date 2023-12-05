import {CSSProperties} from "react";

type Tag = {
    label: string;
    color: string;
    style: CSSProperties | undefined
};

export function Tag({label, color = "blue-100", style}: Tag) {
    const backgroundColor = `bg-${color}`;

    return (
        <>
            <div style={{color: "white"}} className={`${backgroundColor} rounded-full text-sm w-fit p-1`}>
                {label}
            </div>
        </>
    );
}