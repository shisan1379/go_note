package mspool

import (
	"time"
)

const DefaultExpire = 3

type Worker struct {
	pool *Pool
	//task 任务队列
	task chan func()
	//lastTime 执行任务的最后的时间
	lastTime time.Time
}

func (w *Worker) run() {

	go w.running()
}

// 循环运行任务，任务结束后归还 worker
func (w *Worker) running() {
	for f := range w.task {
		f()
		w.pool.PutWorker(w)
	}
}
