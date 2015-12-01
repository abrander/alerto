package agent

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type (
	Agent interface {
		Execute(Request) Result
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

	Description struct {
		Parameters []Parameter `json:"parameters"`
	}

	Parameter struct {
		Name        string   `json:"name"`
		Type        string   `json:"type"`
		Description string   `json:"description"`
		EnumValues  []string `json:"enumValues"`
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

func getParams(elem reflect.Type) []Parameter {
	parameters := []Parameter{}

	l := elem.NumField()

	for i := 0; i < l; i++ {
		f := elem.Field(i)

		jsonName := f.Tag.Get("json")

		if f.Anonymous {
			parameters = append(parameters, getParams(f.Type)...)
		} else if jsonName != "" {
			p := Parameter{}

			p.Name = jsonName
			p.Type = f.Type.String()
			p.Description = f.Tag.Get("description")
			enum := f.Tag.Get("enum")
			if enum != "" {
				p.EnumValues = strings.Split(enum, ",")
				p.Type = "enum"
			} else {
				p.EnumValues = []string{}
			}

			parameters = append(parameters, p)
		}
	}

	return parameters
}

func AvailableAgents() map[string]Description {
	r := make(map[string]Description)

	for name, agent := range agents {
		elem := reflect.TypeOf(agent()).Elem()
		parameters := getParams(elem)

		r[name] = Description{Parameters: parameters}
	}

	return r
}
