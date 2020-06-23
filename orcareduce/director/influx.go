package director

import (
	"fmt"
	"log"
	"time"

	"github.com/influxdata/influxdb/client/v2"
	"github.com/mkuchenbecker/orcareduce/orcareduce"
)

const (
	// MyDB specifies name of database
	MyDB = "go_influx"
)

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Price       float32 `json:"price"`
	Image       string  `json:"image"`
	Description string  `json:"description"`
	Views       string
}

type Run interface {
	StartTime() time.Time
	EndTime() time.Time
	Runtime() time.Duration
	Error() error
	Message() string
	KeyValue() map[string]interface{}
	// Add(key string, value interface{}) Run
}

type EndFunc func(error, string)

type Runtime interface {
	ID() orcareduce.ID
	Runs() []Run
	StartRun() (Run, EndFunc)
	Save() error
}

type Clock interface {
	Now() time.Time
}

type wallClock struct{}

func (w wallClock) Now() time.Time {
	return time.Now()
}

func NewClock() Clock {
	return wallClock{}
}

type simpleRuntime struct {
	identifier orcareduce.ID
	runs       []Run
	clock      Clock
}

func (this *simpleRuntime) ID() orcareduce.ID {
	return this.identifier
}

func (this *simpleRuntime) Runs() []Run {
	return this.runs
}

func (this *simpleRuntime) StartRun() (Run, EndFunc) {
	run := &simpleRun{KV: make(map[string]interface{})}
	run.ID = this.identifier.NewChild()
	this.runs = append(this.runs, run)
	run.Start = this.clock.Now()
	return run, func(err error, msg string) {
		run.End = this.clock.Now()
		run.Err = err
		run.Msg = msg
	}
}

func NewRuntime(parent orcareduce.ID, clock Clock) Runtime {
	return &simpleRuntime{
		identifier: parent.NewScopedChild("runtime"),
		clock:      clock,
		runs:       make([]Run, 0),
	}
}

type simpleRun struct {
	ID    orcareduce.ID
	Start time.Time
	End   time.Time
	Err   error
	Msg   string
	KV    map[string]interface{}
}

func (this *simpleRun) StartTime() time.Time {
	return this.Start
}
func (this *simpleRun) EndTime() time.Time {
	return this.End
}
func (this *simpleRun) Runtime() time.Duration {
	return this.End.Sub(this.Start)
}
func (this *simpleRun) Error() error {
	return this.Err
}
func (this *simpleRun) Message() string {
	return this.Msg
}
func (this *simpleRun) KeyValue() map[string]interface{} {
	this.KV["ID"] = this.ID.String()
	this.KV["startTime"] = this.StartTime()
	this.KV["endTime"] = this.EndTime()
	this.KV["runtime"] = this.Runtime()
	this.KV["error"] = this.Runtime()
	this.KV["Msg"] = this.Message()
	return this.KV
}
func (this *simpleRun) Add(key string, value interface{}) Run {
	this.KV[key] = value
	return this
}

func (s *simpleRuntime) Save() error {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: "http://localhost:8086",
	})
	if err != nil {
		return err
	}
	defer c.Close()
	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  MyDB,
		Precision: "ms",
	})
	if err != nil {
		return err
	}
	for i, run := range s.Runs() {
		tags := map[string]string{
			"runtime": s.ID().String(),
			"run":     fmt.Sprintf("%d", i),
		}
		pt, err := client.NewPoint("run", tags, run.KeyValue(), time.Now())
		if err != nil {
			log.Fatal(err)
		}
		bp.AddPoint(pt)
	}
	if err := c.Write(bp); err != nil {
		return err
	}
	return nil
}

// queryDB convenience function to query the database
func queryDB(cmd string) (res []client.Result, err error) {
	q := client.Query{
		Command:  cmd,
		Database: MyDB,
	}
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: "http://localhost:8086",
	})
	if response, err := c.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		return res, err
	}
	return res, nil
}
