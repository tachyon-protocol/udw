package udwTask

type Task interface {
	Run()

	Stop()
}

type TaskManager interface {
	AddTask(t Task)

	Wait()

	Close()
}

type TaskFunc func()

func (f TaskFunc) Run() {
	f()
}
func (f TaskFunc) Stop() {
}
