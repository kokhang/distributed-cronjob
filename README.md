# Distributed Cronjob
===================

This is a controller for an cronjob system.

## Components

The following components are high level view of the cronjob system.

### API
- Server that takes job requests from the client to be scheduled for execution.
- The API validates inputs requests from the client and fowards requests to a manager to handle them.
- Validates jobs. Jobs should have a execution schedule format like crontab (* * * * *) and the executable. The executable could be a external command, a script or a binary.
- Validates authentication information. Also validate admin(priviledge) users (for example if we have priviledge users who can run jobs using root).

### Manager
- The manager is in charge to handle the scheduling of the job from the user.
- It will schedule the jobs by persisting its metadata to a distributed datastore. Like mysql cluster.
- It will also add tasks for these jobs into a distributed scheduler like zookeeper. These tasks will be consume by workers.
- The manager assigns workers to jobs.
- It watches and monitors worker statuses. If a worker goes down or becomes unreachable, the manager will schedule the jobs in a different worker.

### Workers
- A worker runs on every node in the cluster.
- Workers will consume tasks from a persistent store or scheduler (like zookeeper) and will schedule them in crontab by either modifying crontab config files or using the crontab command to schedule them.
- Workers consume tasks from zookeeper, execute the job and update the job states into the persistent storage.
- If a workers fails during the job execution, the task(znode) on zookeeper is not consumed and another worker can resume the execution.
- Worker will read the job metadata and execute the job. If the executable is a command, it will run the command. If its a script of binary, it will
  get the executable file from a object store (S3), download it to the node and execute it.
- The worker at initialization time will get its job data from zookeeper and the datastore and configure crontab appropriatly. It will also unschedule any tasks that is not supposed to schedule. This is to prevent two workers scheduling the same task.

## Job Model
- Each job will have a unique identifier.
- Each job will have a status. The status is an enumerator from "scheduled", "finished" or "error".
- "scheduled" status means that the job has been received and has been scheduled to execute.
- "finished" status means that the job has been executed and no more execution for this job is scheduled. This could happend when the job is not reocurring or if an user
  has removed this job from the scheduler.
- "error" status means that the job failed during the execution.
- Each job will be of different executable type. It could be a external command which can just run on the node. Or it could be a script or binary, in which the file can
  be retrieved from the object store and be executed on the node.
- Each job will have a scheduling format like crontab ie */1 * * * *

## Technologies used:

- Mysql clustering for storing persistent job data.
- Zookeeper to synchronize jobs and have them consume by worker processes.
- HAproxy/nginx for loadbalancing API servers.
- A block storage to persist executable files like binary or scripts.
- NTP server to synchronize clock in all nodes. (or perhaps we can use the ubuntu ntp)

## REST API

#### GET /jobs
List the state of all jobs.
```
[
  {
    "id": "0",
    "status": "scheduled",
    "executable": "ls -la",
    "schedule": "5 * * * *" 
    "type": "command"
  },
  {
    "id": "1",
    "status": "finished",
    "executable": "http://s3.amazonaws.com/executable/myScript.sh",
    "schedule": "*/10 * * * *"
    "type": "script"
  },
  {
    "id": "0",
    "status": "error",
    "schedule": "*/25 * * * *"
    "executable": "http://s3.amazonaws.com/executable/myBinary",
    "type": "binary"
  },
]
```

### GET /jobs/{id}
Show the state of a job.
```
{
    "id": "0",
    "status": "scheduled",
    "executable": "ls -la",
    "schedule": "5 * * * *" 
    "type": "command"
  },
```

### POST /jobs
Creates a new scheduled job. For scripts or binary executable, the POST request will use upload the file as BASE64.
```
{   
    "executable": "ls -la",
    "schedule": "5 * * * *" 
    "type": "command"
}
```

Response: Scheduled job.
```
{
    "id": "1",
    "status": "scheduled,
    "executable": "ls -la"
    "schedule": "5 * * * *" 
    "type": "command"
}
```

To build the runnable, do the following:

1. Install GoLang
1. Create a work directory. ie mkdir /home/user/work
1. Set $GOPATH and $GOBIN env to point to work directory and bin. I normally set GOBIN to be $GOPATH/bin
1. Create a 'src' directory inside the work directory
1. Clone the project in the 'src' dir as github.com/kokhang/distributed-cronjob.
1. Install dependency into #GOPATH.
    * go get github.com/gorilla/mux
1. Build executable: `go install github.com/kokhang/distributed-cronjob/cronjob.go`
    * Executable is built in the $GOBIN
1. Run executable. This will start a REST-based server and listen to port 30000
  
Use `curl` to invoke the API. For example:
* `curl -X GET http://localhost:30000/jobs`
* `curl -H "Content-Type: application/json" -X POST -d '{"executable":"ls -la", "schedule":"4 * * * *", "type": "command"}' http://localhost:30000/jobs`

If you are not able to build the source code, i have pushed the executable so you can run that instead. It is called `cronjob`.
