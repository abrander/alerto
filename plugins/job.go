package plugins

import (
	"encoding/json"
	"fmt"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type (
	Job struct {
		AgentId string        `json:"agentId" bson:"agentId"`
		Timeout time.Duration `json:"timeout"`
		Agent   Agent         `json:"arguments" bson:"arguments"`
	}
)

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

	a, found := plugins[job.AgentId]
	if !found {
		return fmt.Errorf("unknown agentId '%s'", job.AgentId)
	}

	agent, ok := a().(Agent)
	if !ok {
		return fmt.Errorf("plugin '%s' does not implement plugins.Agent", job.AgentId)
	}

	job.Agent = agent

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

	a, found := plugins[job.AgentId]
	if !found {
		return fmt.Errorf("unknown agentId '%s'", job.AgentId)
	}

	agent, ok := a().(Agent)
	if !ok {
		return fmt.Errorf("plugin '%s' does not implement plugins.Agent", job.AgentId)
	}

	job.Agent = agent

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
