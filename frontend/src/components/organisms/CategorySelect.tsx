import { v4 as uuidv4 } from "uuid";
import {
  Box,
  FormControl,
  InputLabel,
  MenuItem,
  Select,
  SelectChangeEvent,
} from "@mui/material";
import { Category } from "../../types";

type CategorySelectorProps = {
  categories: Category[];
  label: string;
  selected: string[];
  onSelectedUpdate: (selected: string[]) => void;
  width: string;
  multiple?: boolean;
};

export function CategorySelect({
  categories,
  label,
  multiple = false,
  width,
  selected,
  onSelectedUpdate,
}: CategorySelectorProps) {
  const labelId: string = uuidv4();
  const colorByCategory: Map<string, string> = buildColorByCategory();

  function onSelectedChange(event: SelectChangeEvent<string[]>) {
    const {
      target: { value },
    } = event;
    const newValue = typeof value === "string" ? value.split(" ") : value;
    onSelectedUpdate(newValue);
  }

  function buildColorByCategory(): Map<string, string> {
    const colorByCategory = new Map<string, string>();
    categories.forEach((option) => {
      colorByCategory.set(option.name, option.color);
    });

    return colorByCategory;
  }

  function getOptionColor(value: string): string {
    return colorByCategory.get(value) || "gray.main";
  }

  return (
    <>
      <FormControl fullWidth sx={{ background: "white", maxWidth: width }}>
        <InputLabel id={labelId}>{label}</InputLabel>
        <Select
          labelId={labelId}
          id={label}
          label={label}
          value={selected}
          onChange={onSelectedChange}
          multiple={multiple}
          renderValue={(selected) => (
            // This is how items will appear on the select input
            <Box sx={{ display: "flex", flexWrap: "wrap", gap: 0.5 }}>
              {selected.map((value) => (
                <Box
                  key={value}
                  className="p-1 w-fit text-sm rounded-xl"
                  style={{ color: "white" }}
                  sx={{ backgroundColor: getOptionColor(value) }}
                >
                  {value}
                </Box>
              ))}
            </Box>
          )}
        >
          {
            // This is how items will appear on the menu
            categories.map((option) => (
              <MenuItem key={option.name} id={option.name} value={option.name}>
                <Box
                  className="p-1 w-fit text-sm rounded-xl"
                  style={{ color: "white" }}
                  sx={{ backgroundColor: option.color }}
                >
                  {option.name}
                </Box>
              </MenuItem>
            ))
          }
          {/*Material UI complains when passing a value(like an empty string) that is not in the list of options. This
          is why we need this line*/}
          <MenuItem hidden value={""} />
        </Select>
      </FormControl>
    </>
  );
}
