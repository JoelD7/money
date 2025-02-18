package models

import "time"

type SavingGoal struct {
	SavingGoalID string     `json:"saving_goal_id,omitempty"`
	Username     string     `json:"username,omitempty"`
	Name         *string    `json:"name,omitempty"`
	Target       *float64   `json:"target,omitempty"`
	Progress     *float64   `json:"progress,omitempty"`
	Deadline     *time.Time `json:"deadline,omitempty"`
}

func (sg *SavingGoal) SetName(name string) {
	sg.Name = &name
}

func (sg *SavingGoal) SetTarget(target float64) {
	sg.Target = &target
}

func (sg *SavingGoal) SetProgress(progress float64) {
	sg.Progress = &progress
}

func (sg *SavingGoal) SetDeadline(deadline time.Time) {
	sg.Deadline = &deadline
}

func (sg *SavingGoal) GetSavingGoalID() string {
	if sg == nil {
		return ""
	}
	return sg.SavingGoalID
}

func (sg *SavingGoal) GetUsername() string {
	if sg == nil {
		return ""
	}
	return sg.Username
}

func (sg *SavingGoal) GetName() string {
	if sg == nil || sg.Name == nil {
		return ""
	}
	return *sg.Name
}

func (sg *SavingGoal) GetTarget() float64 {
	if sg == nil || sg.Target == nil {
		return 0
	}
	return *sg.Target
}

func (sg *SavingGoal) GetProgress() float64 {
	if sg == nil || sg.Progress == nil {
		return 0
	}
	return *sg.Progress
}

func (sg *SavingGoal) GetDeadline() time.Time {
	if sg == nil || sg.Deadline == nil {
		return time.Time{}
	}
	return *sg.Deadline
}
