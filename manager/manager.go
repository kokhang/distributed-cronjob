package manager

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/kokhang/distributed-cronjob/models"
	"github.com/kokhang/distributed-cronjob/worker"
)

const workerNum = 3

// current worker to schedule job. Do round robin
var curWorker = 0
var workers [workerNum]chan models.Job

// lock object for data access
var mutex = &sync.Mutex{}

func InitWorkers() {
	for i := 0; i < workerNum; i++ {
		workerChannel := worker.AddWorker(i)
		workers[i] = workerChannel
	}
}

func ScheduleNewJob(job models.Job) {
	// Looking objects to synchronize with worker threads. Not needed if using zookeeper
	mutex.Lock()
	defer mutex.Unlock()

	fmt.Printf("Scheduling new job %v\n", job)
	job.ID = strconv.Itoa(models.JobCounter)
	job.Status = models.SCHEDULED
	workers[curWorker] <- job
	curWorker++
	if curWorker == workerNum {
		curWorker = 0
	}
	models.JobCounter++
	models.Jobs = append(models.Jobs, job)
}
