package plugins

import (
	"io"
	"net"
	"reflect"
	"strings"
	"time"

	"github.com/abrander/alerto/logger"
)

type (
	Plugin interface {
	}

	Agent interface {
		Plugin
		Run(Transport, Request) Result
	}

	Transport interface {
		Plugin
		Exec(cmd string, arguments ...string) (io.Reader, io.Reader, error)
		Dial(network string, address string) (net.Conn, error)
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

	Constructor func() Plugin

	Status int
)

const (
	Ok     = 0
	Failed = 1
)

var plugins = map[string]func() Plugin{}

func Register(protocol string, constructor Constructor) {
	_, exists := plugins[protocol]
	if exists {
		logger.Error("plugins", "plugins.Register(): Duplicate protocol: '%s' (%T and %T)\n", protocol, plugins[protocol], constructor())
		return
	}

	plugins[protocol] = constructor
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

func GetPlugin(pluginId string) (Constructor, bool) {
	p, found := plugins[pluginId]

	return p, found
}

func getPlugins(iType reflect.Type) map[string]Description {
	r := make(map[string]Description)

	for name, plugin := range plugins {
		pType := reflect.TypeOf(plugin())
		elem := reflect.TypeOf(plugin()).Elem()
		if pType.Implements(iType) {
			parameters := getParams(elem)

			r[name] = Description{Parameters: parameters}
		}
	}

	return r
}

func AvailableAgents() map[string]Description {
	return getPlugins(reflect.TypeOf((*Agent)(nil)).Elem())
}

func AvailablePlugins() map[string]Description {
	return getPlugins(reflect.TypeOf((*Plugin)(nil)).Elem())
}

func AvailableTransports() map[string]Description {
	return getPlugins(reflect.TypeOf((*Transport)(nil)).Elem())
}
