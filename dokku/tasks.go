package dokku

// Task is an operation performed as part of configuring Dokku.
type Task interface {
}

var registeredTasks = make([]Task, 0)

// RegisterTask adds a task to the list of tasks to run when the program is
// invoked.
func RegisterTask(task Task) {
	registeredTasks = append(registeredTasks, task)
}
