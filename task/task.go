package task

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"git.my-sign.com/backend/coreapi/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/wayt/happyngine/env"
	"github.com/wayt/happyngine/log"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"time"
)

var tasks map[string]*Task
var scheduledTasks *utils.LCFifo
var logger io.Writer = os.Stdout
var taskAPI = "http://localhost:8080"

func init() {

	if apiUrl := env.Get("TASK_API_URL"); len(apiUrl) > 0 {
		taskAPI = apiUrl
	}
	taskLogFile := env.Get("TASK_LOG_FILE")
	if taskLogFile != "" {
		f, err := os.OpenFile(taskLogFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			panic(err)
		}
		// defer f.Close()

		logger = f
	}

	tasks = make(map[string]*Task)
	scheduledTasks = utils.NewListCFifo()

	taskRunnerThreads := env.GetInt("HAPPY_TASK_RUNNER_THREADS")
	if taskRunnerThreads == 0 {
		taskRunnerThreads = runtime.NumCPU()
	}

	for i := 0; i < taskRunnerThreads; i++ {
		go taskRunner()
	}

	go taskScheduler()
}

type Task struct {
	Name string
	fv   reflect.Value // Kind() == reflect.Func
}

func New(name string, i interface{}) *Task {

	if _, ok := tasks[name]; ok {
		panic(errors.New("duplicate task name: " + name))
	}

	t := &Task{
		Name: name,
		fv:   reflect.ValueOf(i),
	}

	f := t.fv.Type()
	if f.Kind() != reflect.Func {
		panic(errors.New("not a function"))
	}

	tasks[name] = t

	return t
}

type TaskSchedule struct {
	Id     string        `json:"id,omitempty"`
	Name   string        `json:"name"`
	Params []interface{} `json:"params,omitempty"`
	Time   int64         `json:"time"`
}

func (t *Task) Schedule(tm time.Time, args ...interface{}) {

	utc := tm.UTC()

	sc := &TaskSchedule{
		Name:   t.Name,
		Params: args,
		Time:   utc.Unix(),
	}

	scheduledTasks.Enqueue(sc)
}

func (t *Task) call(args ...interface{}) error {

	ft := t.fv.Type()
	in := []reflect.Value{}
	for i, arg := range args {
		var v reflect.Value
		if arg != nil {

			paramType := ft.In(i)

			tmp := reflect.New(paramType)
			mapstructure.Decode(arg, tmp.Interface())

			v = tmp.Elem()
		} else {
			// Task was passed a nil argument, so we must construct
			// the zero value for the argument here.
			n := len(in) // we're constructing the nth argument
			var at reflect.Type
			if !ft.IsVariadic() || n < ft.NumIn()-1 {
				at = ft.In(n)
			} else {
				at = ft.In(ft.NumIn() - 1).Elem()
			}
			v = reflect.Zero(at)
		}
		in = append(in, v)
	}

	t.fv.Call(in)

	return nil
}

func getNewTask() *TaskSchedule {

	resp, err := http.Get(fmt.Sprintf("%s/tasks", taskAPI))
	if err != nil {
		log.Errorln("TASK: fail to get tasks:", err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorln("TASK: fail to read tasks:", err)
		return nil
	}

	task := new(TaskSchedule)
	if err := json.Unmarshal(data, task); err != nil {
		log.Errorln("TASK: fail to unmarshal task:", err)
		return nil
	}

	return task
}

func putTask(id, status string, err error) {

	result := struct {
		Status string `json:"status"`
		Error  string `json:"error,omitempty"`
	}{
		Status: status,
	}

	if err != nil {
		result.Error = err.Error()
	}

	data, _ := json.Marshal(result)

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/tasks/%s", taskAPI, id), bytes.NewReader(data))
	if err != nil {
		log.Errorln("TASK: fail to create PUT request:", err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorln("TASK: fail to PUT task status:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {

		data, err := ioutil.ReadAll(resp.Body)
		log.Errorln("TASK: fail to PUT task status 2:", resp.Status, string(data), err)
		return
	}
}

func taskRunner() {

	time.Sleep(15 * time.Second)

	for {

		ts := getNewTask()
		if ts == nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		t, ok := tasks[ts.Name]
		if !ok {
			log.Errorln("TASK: unknown task:", ts.Name)
			putTask(ts.Id, "error", errors.New("Unknown task name"))
			continue
		}

		log.Debugln("TASK: running:", ts.Name)
		t.call(ts.Params...)

		putTask(ts.Id, "done", nil)

		// took := time.Since(startTime)
		//
		// if _, err := fmt.Fprintf(logger, "%s [%s] %d %d\n", taskName, time.Now().Format("2/Jan/2006:15:04:05 -0700"), status, took.Nanoseconds()/1000000); err != nil {
		// 	log.Errorln("taskRunner:", err)
		// }
	}
}

func taskScheduler() {

	for true {

		i, ok := scheduledTasks.Dequeue()
		if !ok {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		task := i.(*TaskSchedule)

		data, err := json.Marshal(task)
		if err != nil {
			log.Errorln("TASK: failed to marshal:", task, err)
			continue
		}

		resp, err := http.Post(fmt.Sprintf("%s/tasks", taskAPI), "application/json", bytes.NewReader(data))
		if err != nil {
			log.Errorln("TASK: failed to post task:", task, err)

			// Requeue task
			scheduledTasks.Enqueue(task)
			time.Sleep(500 * time.Millisecond)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			data, err = ioutil.ReadAll(resp.Body)
			log.Errorln("TASK: failed to create task:", resp.Status, string(data), err)

			// Requeue task
			scheduledTasks.Enqueue(task)
			time.Sleep(500 * time.Millisecond)
			continue
		}
	}
}
