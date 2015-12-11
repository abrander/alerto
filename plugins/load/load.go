package load

import (
	"bufio"
	"strconv"
	"strings"

	"github.com/abrander/alerto/plugins"
)

func init() {
	plugins.Register("load", NewLoad)
}

func NewLoad() plugins.Plugin {
	return new(Load)
}

type (
	Load struct {
	}
)

func (l Load) GetInfo() plugins.HumanInfo {
	return plugins.HumanInfo{
		Name:        "Load",
		Description: "Read system load",
	}
}

func (l *Load) Run(transport plugins.Transport, request plugins.Request) plugins.Result {
	file, err := transport.ReadFile("/proc/loadavg")
	if err != nil {
		return plugins.NewResult(plugins.Failed, nil, err.Error())
	}

	var load1 float64
	var load5 float64
	var load15 float64
	var activeTasks int64
	var tasks int64

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()

		data := strings.Fields(strings.Trim(text, " "))
		if len(data) != 5 {
			continue
		}

		load1, _ = strconv.ParseFloat(data[0], 64)
		load5, _ = strconv.ParseFloat(data[1], 64)
		load15, _ = strconv.ParseFloat(data[2], 64)

		sep := strings.Index(data[3], "/")
		if sep > 0 {
			activeTasks, _ = strconv.ParseInt(data[3][0:sep], 10, 64)
			tasks, _ = strconv.ParseInt(data[3][sep+1:], 10, 64)
		}
	}

	return plugins.NewResult(
		plugins.Ok,
		plugins.NewMeasurementCollection(
			"load1", load1,
			"load5", load5,
			"load15", load15,
			"activeTasks", activeTasks,
			"tasks", tasks),
		"%.02f %.02f %.02f %d/%d", load1, load5, load15, activeTasks, tasks)
}

// Ensure compliance
var _ plugins.Agent = (*Load)(nil)
