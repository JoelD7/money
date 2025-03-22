import { FontAwesomeIcon as ReactFontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { IconProp } from "@fortawesome/fontawesome-svg-core";

type FontAwesomeIconProps = {
  icon: IconProp;
  colorClassName?: string;
};

// Custom implementation of the FontAwesomeIcon component that uses a CSS class(like Tailwind's text-blue-500 for example)
// to set the color of the icon.
export function FontAwesomeIcon({ icon, colorClassName }: FontAwesomeIconProps) {
  return (
    <span className={colorClassName}>
      <ReactFontAwesomeIcon icon={icon} />
    </span>
  );
}
