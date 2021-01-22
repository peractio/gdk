package cronx

import (
	"context"
	"reflect"
	"sync"
	"sync/atomic"
	"time"

	"github.com/robfig/cron/v3"
)

type JobItf interface {
	Run(ctx context.Context) error
}

type Job struct {
	Name    string       `json:"name"`
	Status  StatusCode   `json:"status"`
	Latency string       `json:"latency"`
	Error   string       `json:"error"`
	EntryID cron.EntryID `json:"entry_id"`

	inner   JobItf
	status  uint32
	running sync.Mutex
}

// UpdateStatus updates the current job status to the latest.
func (j *Job) UpdateStatus() StatusCode {
	switch atomic.LoadUint32(&j.status) {
	case statusRunning:
		j.Status = StatusCodeRunning
	case statusIdle:
		j.Status = StatusCodeIdle
	case statusDown:
		j.Status = StatusCodeDown
	case statusError:
		j.Status = StatusCodeError
	default:
		j.Status = StatusCodeUp
	}
	return j.Status
}

// Run executes the current job operation.
func (j *Job) Run() {
	start := time.Now()
	ctx := context.Background()

	// Lock current process.
	j.running.Lock()
	defer j.running.Unlock()

	// Update job status as running.
	atomic.StoreUint32(&j.status, statusRunning)
	j.UpdateStatus()

	// Run the job.
	if err := commandController.Interceptor(ctx, j, func(ctx context.Context, job *Job) error {
		return job.inner.Run(ctx)
	}); err != nil {
		j.Error = err.Error()
		atomic.StoreUint32(&j.status, statusError)
	} else {
		atomic.StoreUint32(&j.status, statusIdle)
	}

	// Record time needed to execute the whole process.
	j.Latency = time.Since(start).String()

	// Update job status after running.
	j.UpdateStatus()
}

// NewJob creates a new job with default status and name.
func NewJob(job JobItf) *Job {
	name := reflect.TypeOf(job).Name()
	if name == "" {
		name = reflect.TypeOf(job).Elem().Name()
	}
	if name == "" {
		name = reflect.TypeOf(job).String()
	}
	if name == "Func" {
		name = "(nameless)"
	}

	return &Job{
		Name:   name,
		Status: StatusCodeUp,
		inner:  job,
		status: statusUp,
	}
}
