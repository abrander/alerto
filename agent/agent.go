package agent

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type (
	Agent interface {
		Execute(Request) Result
	}

	Job struct {
		AgentId   string          `json:"agentId" bson:"agentId"`
		Timeout   time.Duration   `json:"timeout"`
		Agent     Agent           `json:"-" bson:"arguments"` // This is some hackery to support both json and bson
		Arguments json.RawMessage `json:"-,omitempty" bson:"-"`
	}

	Request struct {
		Timeout time.Duration
	}

	Result struct {
		Status       Status
		Text         string
		Duration     time.Duration
		Measurements *MeasurementCollection
	}

	Constructor func() Agent

	Status int
)

const (
	Ok     = 0
	Failed = 1
)

var agents = map[string]func() Agent{}

func Register(protocol string, constructor Constructor) {
	_, exists := agents[protocol]
	if exists {
		log.Fatal("agent.Register(): Duplicate protocol: '%s' (%T and %T)\n", protocol, agents[protocol], constructor())
		return
	}

	agents[protocol] = constructor
}

func (job *Job) UnmarshalJSON(data []byte) error {
	m := make(map[string]json.RawMessage)

	err := json.Unmarshal(data, &m)
	if err != nil {
		return err
	}

	agentRaw, found := m["agentId"]
	if !found {
		return fmt.Errorf("agentId not found in document")
	}

	argumentsRaw, found := m["arguments"]
	if !found {
		return fmt.Errorf("arguments not found in document")
	}

	err = json.Unmarshal(agentRaw, &job.AgentId)
	if err != nil {
		return err
	}

	a, found := agents[job.AgentId]
	if !found {
		return fmt.Errorf("unknown agentId '%s'", job.AgentId)
	}

	job.Agent = a()

	err = json.Unmarshal(argumentsRaw, job.Agent)
	if err != nil {
		return err
	}

	return nil
}

func (job *Job) SetBSON(raw bson.Raw) error {
	m := make(map[string]bson.Raw)

	err := bson.Unmarshal(raw.Data, &m)
	if err != nil {
		panic(err.Error())
	}

	agentRaw, found := m["agentId"]
	if !found {
		return fmt.Errorf("agentId not found in document")
	}

	err = agentRaw.Unmarshal(&job.AgentId)
	if err != nil {
		return err
	}

	timeoutRaw, found := m["timeout"]
	if !found {
		job.Timeout = time.Second * 10
	} else {
		err = timeoutRaw.Unmarshal(&job.Timeout)
		if err != nil {
			job.Timeout = time.Second * 10
		}
	}

	a, found := agents[job.AgentId]
	if !found {
		return fmt.Errorf("unknown agent '%s'", job.AgentId)
	}

	job.Agent = a()

	argumentsRaw, found := m["arguments"]
	if !found {
		return fmt.Errorf("arguments not found in document")
	}

	err = argumentsRaw.Unmarshal(job.Agent)
	if err != nil {
		return err
	}

	return nil
}

func (job *Job) Run() Result {
	start := time.Now()

	request := Request{
		Timeout: job.Timeout,
	}

	if request.Timeout == 0 {
		request.Timeout = time.Second
	}

	result := job.Agent.Execute(request)
	result.Duration = time.Now().Sub(start)

	return result
}

func NewResult(status Status, measurements *MeasurementCollection, format string, args ...interface{}) Result {
	return Result{
		Status:       status,
		Text:         fmt.Sprintf(format, args...),
		Measurements: measurements,
	}
}
