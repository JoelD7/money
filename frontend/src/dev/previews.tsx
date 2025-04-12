import React from "react";
import {ComponentPreview, Previews} from "@react-buddy/ide-toolbox";
import {PaletteTree} from "./palette";
import {Button} from "../components/atoms/Button.tsx";
import App from "../App.tsx";
import {DeleteSavingGoal} from "../components/molecules/DeleteSavingGoal.tsx";

const ComponentPreviews = () => {
    return (
      <Previews palette={<PaletteTree />}>
        <ComponentPreview path="/Button">
          <Button />
        </ComponentPreview>
        <ComponentPreview path="/App">
          <App />
        </ComponentPreview>
        <ComponentPreview path="/DeleteSavingGoal">
          <DeleteSavingGoal />
        </ComponentPreview>
      </Previews>
    );
};

export default ComponentPreviews;