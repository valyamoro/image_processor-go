package worker

import (
	"container/heap"
	"image_processor-go/internal/core"
	"sync"
	"time"
)

type PriorityQueue []*core.Job

func (pq PriorityQueue) Len() int			{ return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].Priority < pq[j].Priority }
func (pq PriorityQueue) Swap(i, j int)		{ pq[i], pq[j] = pq[j], pq[i] }

func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(*core.Job))
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n - 1]
	*pq = old[0 : n-1]
	return item
}

type Pool struct {
	workerCount		int
	workers			[]*Worker
	jobQueue		chan *core.Job
	resultChan		chan string
	priorityQueue	*PriorityQueue
	queueMutex		sync.Mutex
	wg				sync.WaitGroup
	stopChan		chan struct{}
}

func NewPool(workerCount, queueSize int) *Pool {
	pq := &PriorityQueue{}
	heap.Init(pq)

	return &Pool{
		workerCount: 	workerCount,
		jobQueue: 		make(chan *core.Job, queueSize),
		resultChan:		make(chan string, 1000),
		priorityQueue:	pq,
		stopChan:		make(chan struct{}),
	}
}

func (p *Pool) Start() {
	for i := 0; i < p.workerCount; i++ {
		worker := NewWorker(i+1, p.jobQueue, p.resultChan)
		p.workers = append(p.workers, worker)

		p.wg.Add(1)
		go func(w *Worker) {
			defer p.wg.Done()
			w.Start()
		}(worker)
	}

	go p.queueDispatcher()
}

func (p *Pool) Submit(job *core.Job) {
	p.queueMutex.Lock()
	heap.Push(p.priorityQueue, job)
	p.queueMutex.Unlock()
}

func (p *Pool) queueDispatcher() {
	for {
		select {
		case <-p.stopChan:
			return
		default:
			p.queueMutex.Lock()
			if p.priorityQueue.Len() > 0 {
				job := heap.Pop(p.priorityQueue).(*core.Job)
				p.queueMutex.Unlock()

				p.jobQueue <- job
			} else {
				p.queueMutex.Unlock()
				time.Sleep(10 * time.Millisecond)
			}
		}
	}
}

func (p *Pool) Stop() {
	close(p.stopChan)
	close(p.jobQueue)
	p.wg.Wait()
	close(p.resultChan)
}

func (p *Pool) GetResultChan() <-chan string {
	return p.resultChan
}

func (p *Pool) GetWorkerStats() []*core.WorkerStats {
    stats := make([]*core.WorkerStats, len(p.workers))
    for i, w := range p.workers {
        stats[i] = w.GetStats()
    }
    return stats
}

func (p *Pool) GetQueueLength() int {
    p.queueMutex.Lock()
    defer p.queueMutex.Unlock()
    return p.priorityQueue.Len()
}
