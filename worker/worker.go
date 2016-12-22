package worker

import (
	"fmt"

	"github.com/kokhang/distributed-cronjob/models"
)

/* This would normally run as its own process. But for the sake of this prototype, i will run is as a thread.
 */

func AddWorker(workerId int) chan models.Job {
	var job = make(chan models.Job)
	go workerExec(workerId, job)
	return job
}

func workerExec(workerID int, job chan models.Job) {
	fmt.Printf("Worker %v starting...\n", workerID)
	for {
		j := <-job
		fmt.Printf("Scheduling job %v on worker %v\n", j, workerID)
		// TODO
		// Parse job.
		// if job is a script or binary, download from object storage to the node
		// modify crontab to schedule job using exec.Command
	}
}
