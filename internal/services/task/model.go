package task

type TaskCreateInput struct {
	Title       string
	Description string
	Priority    int
}

type TaskUpdateInput struct {
	Title       string
	Description string
	Priority    int
}

type TaskUpdateStatusInput struct {
	Status string
}
