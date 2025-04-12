import { Dialog } from "./Dialog.tsx";
import { Divider, Typography } from "@mui/material";
import { SavingGoal, SnackAlert } from "../../types";
import { Button } from "../atoms";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import api from "../../api";
import { savingGoalKeys } from "../../queries/saving_goals.ts";

type DeleteSavingGoalProps = {
  savingGoal: SavingGoal;
  open: boolean;
  onClose: () => void;
  onAlert: (alert: SnackAlert) => void;
};

export function DeleteSavingGoal({
  savingGoal,
  open,
  onClose,
  onAlert,
}: DeleteSavingGoalProps) {
  const queryClient = useQueryClient();

  const deleteSavingGoalMutation = useMutation({
    mutationFn: (id: string) => api.deleteSavingGoal(id),
    onSuccess: () => {
      onAlert({
        open: true,
        type: "success",
        title: "Saving goal deleted successfully",
      });

      queryClient
        .invalidateQueries({ queryKey: savingGoalKeys.all })
        .then(() => {})
        .catch((e) => {
          console.error("Error invalidating saving goals query", e);
        });
    },
    onError: (error) => {
      onAlert({
        open: true,
        type: "error",
        title: "Error deleting saving goal",
      });
      console.error("Error deleting saving goal", error);
    },
  });

  function handleDelete() {
    onClose();
    deleteSavingGoalMutation.mutate(savingGoal.saving_goal_id);
  }

  return (
    <Dialog open={open} onClose={onClose} fullWidth>
      <div className={"max-w-[600px]"}>
        <Typography variant={"h4"}>Delete goal</Typography>
        <Divider />

        <div className={"pt-4"}>
          <Typography variant={"body1"}>
            {`Are you sure you want to delete goal `}"
            <span className={"font-bold"}>{savingGoal.name}</span>"?
          </Typography>

          <Typography variant={"body1"} marginTop={"0.5rem"}>
            The saving entries associated with this goal won’t be deleted, but will remain
            without a saving goal.
          </Typography>
        </div>

        <div className={"flex justify-end pt-6"}>
          <Button
            variant={"contained"}
            color={"gray"}
            onClick={onClose}
            sx={{ marginRight: "5px" }}
          >
            Cancel
          </Button>
          <Button color={"red"} variant={"contained"} onClick={handleDelete}>
            Delete
          </Button>
        </div>
      </div>
    </Dialog>
  );
}
