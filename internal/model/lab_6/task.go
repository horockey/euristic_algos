package model

type Task struct {
	Name string
	Cost []int
}

func NewTask(name string, cost []int) *Task {
	return &Task{
		Name: name,
		Cost: cost,
	}
}
