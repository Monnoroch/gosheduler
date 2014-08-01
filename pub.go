package sheduler;

// The common interface for task output for type safety
type TaskOutput interface {}

// The task interface, represents any runnable object
type Task interface {
	Run() TaskOutput
}

// The struct for representing the output assotiated with particular task
type TaskResult struct {
	Task Task
	Output TaskOutput
}

// The sheduler interface: the main thing, which runs tasks and returns a channel to receive the task results
type Sheduler interface {
	Run()
	Do(task Task)
	Stop()
	Restart()
	Close()
	Result() chan *TaskResult
}


// The sheduler built with native semaphor, which starts a new goroutine for each task.
// maxJobs is a maximum number of parallel jobs to run
func NewShedulerNative(maxJobs int) Sheduler {
	return newShedulerNative(maxJobs)
}

// The tricky sheduler, which starts maxJobs goroutines at once and then they wait for incoming jobs.
// This one might be better, because of eliding of allocation on new goroutines and associated costs (context ptr, stack, etc).
// maxJobs is a maximum number of parallel jobs to run
func NewShedulerNGoroutines(maxJobs int) Sheduler {
	return newShedulerNGoroutines(maxJobs)
}


