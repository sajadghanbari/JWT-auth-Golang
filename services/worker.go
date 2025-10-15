package services

import "log"

var JobQueue chan MailJob

type Worker struct {
	WorkerPool chan chan MailJob
	JobChannel chan MailJob
	QuitChan   chan bool
	Mailer     *Mailer
}

func NewWorker(workerPool chan chan MailJob, mailer *Mailer) *Worker {
	return &Worker{
		WorkerPool: workerPool,
		JobChannel: make(chan MailJob),
		QuitChan:   make(chan bool),
		Mailer:     mailer,
	}
}

func (w *Worker) Start() {
	go func() {
		for {
			w.WorkerPool <- w.JobChannel

			select {
			case job := <-w.JobChannel:
				if err := w.Mailer.SendMail(&job); err != nil {
					// Handle error if needed
					log.Printf("Error sending mail: %s", err)
				}
			case <-w.QuitChan:
				return
			}
		}
	}()
}

func (w Worker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}

type Dispatcher struct {
	WorkerPool chan chan MailJob
	MaxWorkers int
	Mailer     *Mailer
}

func NewDispatcher(maxWorkers int, mailer *Mailer) *Dispatcher {
	pool := make(chan chan MailJob, maxWorkers)
	return &Dispatcher{
		WorkerPool: pool,
		MaxWorkers: maxWorkers,
		Mailer:     mailer,
	}
}

func (d *Dispatcher) Run() {
	for i := 0; i < d.MaxWorkers; i++ {
		worker := NewWorker(d.WorkerPool, d.Mailer)
		worker.Start()
	}
	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	for job := range JobQueue {
		// A job request has been received
		go func(job MailJob) {
			// Try to obtain a worker job channel that is available.
			// This will block until a worker is idle
			jobChannel := <-d.WorkerPool

			// Dispatch the job to the worker job channel
			jobChannel <- job
		}(job)
	}
}