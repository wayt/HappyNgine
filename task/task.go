package task

import (
	"encoding/json"
	"errors"
	"github.com/mitchellh/mapstructure"
	"github.com/wayt/happyngine/env"
	"github.com/wayt/happyngine/log"
	"gopkg.in/redis.v3"
	"reflect"
	"time"
)

var redisCli *redis.Client
var scheduledTasksKey = "scheduled_tasks" // Tasks pushed by the cli, waiting to be push in todo
var todoTasksKey = "todo_tasks"           // Tasks pushed by the scheduler, waiting to be executed
var tasks map[string]*Task

func init() {

	poolSize := env.GetInt("HAPPY_REDIS_TASK_POOL_SIZE")
	if poolSize <= 0 {
		poolSize = 10
	}

	poolTimeout := time.Duration(env.GetInt("HAPPY_REDIS_TASK_POOL_TIMEOUT")) * time.Millisecond
	if poolTimeout <= 0 {
		poolTimeout = time.Second * 5
	}

	redisCli = redis.NewClient(&redis.Options{
		Addr:        env.Get("REDIS_TASK_PORT_6379_TCP_ADDR") + ":" + env.Get("REDIS_TASK_PORT_6379_TCP_PORT"),
		Password:    env.Get("HAPPY_REDIS_TASK_PASSWORD"),
		DB:          int64(env.GetInt("HAPPY_REDIS_TASK_DB")),
		PoolSize:    poolSize,
		PoolTimeout: poolTimeout,
	})

	tasks = make(map[string]*Task)

	for i := 0; i < 4; i++ {
		go taskRunner()
	}
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
	Name string
	Args []interface{}
	Time time.Time
}

func (ts *TaskSchedule) MarshalBinary() ([]byte, error) {
	return json.Marshal(ts)
}

func (ts *TaskSchedule) UnmarshalBinary(data []byte) error {

	return json.Unmarshal(data, &ts)
}

func (t *Task) Schedule(tm time.Time, args ...interface{}) error {

	timestamp := tm.Unix()

	sc := &TaskSchedule{
		Name: t.Name,
		Args: args,
		Time: time.Now(),
	}

	err := redisCli.ZAdd(scheduledTasksKey, redis.Z{
		Score:  float64(timestamp),
		Member: sc,
	}).Err()

	return err
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

func taskRunner() {
	for {
		task, err := redisCli.BLPop(0, todoTasksKey).Result()
		if err != nil {
			log.Errorln("taskRunner: redisCli.BLPop:", err)
			time.Sleep(1 * time.Second)
			continue
		}

		ts := &TaskSchedule{}
		if err := ts.UnmarshalBinary([]byte(task[1])); err != nil {
			log.Errorln("taskRunner: UnmarshalBinary:", err)
			continue
		}

		t, ok := tasks[ts.Name]
		if !ok {
			log.Errorln("taskRunner: unknown task:", ts.Name)
			continue
		}

		log.Debugln("TASK: running:", ts.Name)
		t.call(ts.Args...)
	}
}
