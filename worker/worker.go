package worker

import (
	"github.com/vladivolo/go_daemon_skeleton/lib/queue"
	"github.com/vladivolo/go_daemon_skeleton/lib/redis"
	log "github.com/vladivolo/lumber"
	"sync/atomic"
)

type WorkerPool struct {
	workers     []*Worker
	InputQueue  *redis.RedisConn
	OutputQueue *redis.RedisConn
	MsgQueue    *queue.Queue
	Notify_chan chan int
}

type Worker struct {
	input_queue *queue.Queue
	notify_chan chan int
	status      uint32
}

const (
	CMD_NEW_DATA = 1
	CMD_EXIT     = 2

	STATUS_INIT     = 0
	STATUS_WORKING  = 1
	STATUS_SHUTDOWN = 2
	STATUS_FAILED   = 3
)

const CMD_CHAN_SIZE = 1024 * 100

var (
	Pool *WorkerPool
)

type Request struct {
	Data []byte
}

func (p *WorkerPool) Workers() []*Worker {
	return p.workers
}

func (p *WorkerPool) AddWorker(w *Worker) *Worker {
	p.workers = append(p.workers, w)
	return w
}

func (p *WorkerPool) MQ() *queue.Queue {
	return p.MsgQueue
}

func (p *WorkerPool) NotifyChan() chan int {
	return p.Notify_chan
}

func (p *WorkerPool) Notify(cmd int) {
	p.NotifyChan() <- cmd
}

func (p *WorkerPool) PushData(r *Request) {
	p.MQ().Push(r)
	p.Notify(CMD_NEW_DATA)

	return
}

func (p *WorkerPool) RunWorker(w *Worker) {
	go WorkerProcess(w)
}

func NewWorkerPool(workers_count int) *WorkerPool {
	log.Info("NewWorkerPool() %d", workers_count)

	Pool = &WorkerPool{
		MsgQueue:    queue.NewQueue(),
		Notify_chan: make(chan int, CMD_CHAN_SIZE),
	}

	for i := 0; i < workers_count; i++ {
		Pool.RunWorker(Pool.AddWorker(NewWorker(Pool.MQ(), Pool.NotifyChan())))
	}

	return Pool
}

func NewWorker(q *queue.Queue, c chan int) *Worker {
	w := &Worker{
		input_queue: q,
		notify_chan: c,
		status:      STATUS_INIT,
	}

	return w
}

func (p *WorkerPool) AddInputQueue(Params string) (err error) {
	p.InputQueue, err = redis.NewRedisConnection(Params)
	if err != nil {
		return err
	}

	p.InputQueue.OpenQueue()
	p.InputQueue.StartConsuming(InputQueueHandle)

	return nil
}

func InputQueueHandle(payload string) {
	log.Info("InputQueueHandle: %s", payload)

	Pool.MQ().Push(Request{Data: []byte(payload)})
	Pool.Notify(CMD_NEW_DATA)
}

func (w *Worker) ProcessHandler(request interface{}, len int) {
	log.Debug("ProcessHandler: ")
}

func (w *Worker) SetStatus(status uint32) {
	atomic.StoreUint32(&w.status, status)
}

func (w *Worker) Status() uint32 {
	return atomic.LoadUint32(&w.status)
}

func WorkerProcess(w *Worker) {
	w.SetStatus(STATUS_WORKING)

	for {
		select {
		case cmd := <-w.notify_chan:
			switch cmd {
			case CMD_NEW_DATA:
				w.ProcessHandler(Pool.MQ().Poll(), Pool.MQ().Len())
			case CMD_EXIT:
				log.Info("Worker GOT EXIT!!!")
				w.SetStatus(STATUS_SHUTDOWN)
				return
			}
		}
	}
}
