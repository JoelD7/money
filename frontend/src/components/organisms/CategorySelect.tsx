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
  selected: Category[];
  onSelectedUpdate: (selected: Category[]) => void;
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
  const categoriesByName: Map<string, Category> = buildCategoryMap();

  function onSelectedChange(event: SelectChangeEvent<string[]>) {
    const {
      target: { value },
    } = event;
    const newValue = typeof value === "string" ? value.split(" ") : value;

    const category = newValue.map((value) => categoriesByName.get(value) as Category);
    onSelectedUpdate(category);
  }

  function buildCategoryMap(): Map<string, Category> {
    const map = new Map<string, Category>();

    categories.forEach((category) => {
      map.set(category.name, category);
    });

    return map
  }

  function getOptionColor(value: string): string {
    const category = categoriesByName.get(value);
    if (category) {
      return category.color
    }

    return "gray.main";
  }

  return (
    <>
      <FormControl fullWidth sx={{ background: "white", maxWidth: width }}>
        <InputLabel id={labelId}>{label}</InputLabel>
        <Select
          labelId={labelId}
          id={label}
          label={label}
          value={selected.map((category) => category.name)}
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
          <MenuItem hidden disabled value={""} />
        </Select>
      </FormControl>
    </>
  );
}
