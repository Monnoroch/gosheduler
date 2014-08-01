// Package provides the task sheduler abstraction, which can shedule parallel tasks and collect their results.
package sheduler

/////////////////////////////////
// TODO: find out best sheduler
/////////////////////////////////

type shedulerNGoroutines struct {
	result chan *TaskResult
	count  int
	input  chan Task
	inputs []chan Task
	free   chan int
}

func (self *shedulerNGoroutines) Run() {
	for i := 0; i < self.count; i++ {
		go func(n int) {
			self.free <- n
			for task := range self.inputs[n] {
				self.result <- &TaskResult{
					Task:   task,
					Output: task.Run(),
				}
				self.free <- n
			}
		}(i)
	}
	for task := range self.input {
		self.inputs[<-self.free] <- task
	}
}

func (self *shedulerNGoroutines) Do(task Task) {
	self.input <- task
}

func (self *shedulerNGoroutines) Stop() {
	if self.input != nil {
		close(self.input)
		self.input = nil
	}

	for i := 0; i < len(self.inputs); i++ {
		if self.inputs[i] != nil {
			close(self.inputs[i])
			self.inputs[i] = nil
		}
	}
}

func (self *shedulerNGoroutines) Close() {
	self.Stop()
	close(self.result)
}

func (self *shedulerNGoroutines) Restart() {
	self.input = make(chan Task)
	for i := 0; i < len(self.inputs); i++ {
		self.inputs[i] = make(chan Task)
	}
	go self.Run()
}

func (self *shedulerNGoroutines) Result() chan *TaskResult {
	return self.result
}

func newShedulerNGoroutines(cnt int) *shedulerNGoroutines {
	r := shedulerNGoroutines{
		count:  cnt,
		result: make(chan *TaskResult),
		input:  make(chan Task),
		inputs: make([]chan Task, cnt),
		free:   make(chan int),
	}

	for i := 0; i < cnt; i++ {
		r.inputs[i] = make(chan Task)
	}

	return &r
}

type shedulerNative struct {
	result   chan *TaskResult
	input    chan Task
	semaphor chan struct{}
}

func (self *shedulerNative) Run() {
	for task := range self.input {
		task := task
		self.semaphor <- struct{}{}
		go func() {
			res := task.Run()
			self.result <- &TaskResult{
				Task:   task,
				Output: res,
			}
			<-self.semaphor
		}()
	}
}

func (self *shedulerNative) Do(task Task) {
	self.input <- task
}

func (self *shedulerNative) Stop() {
	if self.input != nil {
		close(self.input)
		self.input = nil
	}
}

func (self *shedulerNative) Close() {
	self.Stop()
	close(self.result)
}

func (self *shedulerNative) Restart() {
	self.input = make(chan Task)
	go self.Run()
}

func (self *shedulerNative) Result() chan *TaskResult {
	return self.result
}

func newShedulerNative(cnt int) *shedulerNative {
	r := shedulerNative{
		input:    make(chan Task),
		result:   make(chan *TaskResult),
		semaphor: make(chan struct{}, cnt),
	}

	return &r
}
