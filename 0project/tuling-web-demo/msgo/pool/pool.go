package mspool

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

type sig struct {
}

var (
	ErrorInValidCap    = errors.New("pool cap  can not < 0")
	ErrorInValidExpire = errors.New("pool expire  can not < 0")
)

type Pool struct {
	//cap 容量 pool max cap
	cap int32
	//running 正在运行的worker的数量
	running int32
	//空闲worker
	workers []*Worker
	//expire 过期时间 空闲的worker超过这个时间 回收掉
	expire time.Duration
	//release 释放资源  pool就不能使用了
	release chan sig
	//lock 去保护pool里面的相关资源的安全
	lock sync.Mutex
	//once 释放只能调用一次 不能多次调用
	once sync.Once
}

func NewPool(cap int) (*Pool, error) {
	return NewTimePool(cap, DefaultExpire)
}
func NewTimePool(cap int, expire int) (*Pool, error) {
	if cap <= 0 {
		return nil, ErrorInValidCap
	}
	if expire <= 0 {
		return nil, ErrorInValidExpire
	}
	p := &Pool{
		cap:     int32(cap),
		expire:  time.Duration(expire) * time.Second,
		release: make(chan sig, 1),
	}
	return p, nil
}

// 提交任务
func (p *Pool) Submit(task func()) error {

	//获取池中的一个 worker 然后执行
	w := p.GetWorker()
	w.task <- task
	w.pool.incRunning()
	return nil
}
func (p *Pool) GetWorker() *Worker {

	workers := p.workers
	l := len(workers) - 1

	//如果没有空闲的 worker 则新建
	//如果 运行的worker + 空闲的worker ，如果大于 pool 的容量，就阻塞等待，等待 worker 释放
	if l >= 0 {
		p.lock.Lock()
		worker := workers[l] //取出最后一个
		workers[l] = nil
		p.workers = workers[:l] //
		p.lock.Unlock()
		return worker
	}
	if p.running < p.cap {
		//没有可用的 worker
		w := &Worker{
			pool: p,
			task: make(chan func(), 1),
			//lastTime:
		}
		w.run()
		return w
	}
	// 阻塞等待
	//for {
	//worker := p.GetWorker()
	//if worker == nil {
	//	continue
	//} else {
	//	return worker
	//}
	//}
	//4. 如果正在运行的workers 如果大于pool容量，阻塞等待，worker释放
	for {
		p.lock.Lock()
		idleWorkers := p.workers
		n := len(idleWorkers) - 1
		if n < 0 {
			p.lock.Unlock()
			continue
		}
		w := idleWorkers[n]
		idleWorkers[n] = nil
		p.workers = idleWorkers[:n]
		p.lock.Unlock()
		return w
	}

	return nil
}

func (p *Pool) incRunning() {
	atomic.AddInt32(&p.running, 1)
}
func (p *Pool) decRunning() {
	atomic.AddInt32(&p.running, -1)
}
func (p *Pool) PutWorker(w *Worker) {
	p.lock.Lock()
	defer p.lock.Unlock()

	w.lastTime = time.Now()
	p.workers = append(p.workers, w)
	p.decRunning()

}
