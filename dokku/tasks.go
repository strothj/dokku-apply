package dokku

// Task is an operation performed as part of configuring Dokku.
type Task interface {
}

var registeredTasks = make(map[string]Task)
