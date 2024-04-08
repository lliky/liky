package main

import (
	"fmt"
	"time"
)

type Worker struct {
	name string
}
type Processor struct {
	jobChan chan string
	workers []*Worker
	done    chan *Worker
}

func (w *Worker) processJob(data string, done chan *Worker) {
	go func() {
		fmt.Printf("working on data : %s, name: %s\n", data, w.name)
		time.Sleep(time.Millisecond)
		done <- w
	}()
}

func NewProcessor() *Processor {
	p := &Processor{
		jobChan: make(chan string),
		workers: make([]*Worker, 5),
		done:    make(chan *Worker),
	}
	for i := 0; i < 5; i++ {
		p.workers[i] = &Worker{
			name: fmt.Sprintf("<worker-%d>", i),
		}
	}
	p.StartProcess()
	return p
}

func (p *Processor) StartProcess() {
	go func() {
		for {
			select {
			case w := <-p.done:
				p.workers = append(p.workers, w)
			default:
				if len(p.workers) > 0 {
					w := p.workers[0]
					p.workers = p.workers[1:]
					w.processJob(<-p.jobChan, p.done)
				}
			}
		}
	}()
}

func (p *Processor) PostJob(data string) {
	p.jobChan <- data
}
