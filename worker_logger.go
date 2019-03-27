package cuber

type WorkerLogger interface {
	Println(...interface{})
	Printf(string, ...interface{})
}
