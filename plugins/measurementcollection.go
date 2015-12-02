package plugins

import (
	"fmt"
	"time"

	"github.com/abrander/alerto/logger"
)

type (
	MeasurementCollection map[string]Measurement
)

func NewMeasurementCollection(args ...interface{}) *MeasurementCollection {
	c := MeasurementCollection(make(map[string]Measurement))

	l := len(args)

	if l > 0 && l%2 == 0 {
		l /= 2
		for i := 0; i < l; i++ {
			key := args[i*2].(string)
			value := args[i*2+1]

			switch value.(type) {
			case int:
				c[key] = Measurement(value.(int))
			case int64:
				c[key] = Measurement(value.(int64))
			case uint64:
				c[key] = Measurement(value.(uint64))
			case float32:
				c[key] = Measurement(value.(float32))
			case float64:
				c[key] = Measurement(value.(float64))
			case time.Duration:
				c[key] = Measurement(value.(time.Duration).Nanoseconds())
			default:
				logger.Error("plugins", "Unsupported type")
			}
		}
	} else if l > 0 {
		logger.Error("plugins", "Wrong number of arguments to NewMeasurementCollection()")
	}

	return &c
}

func (c MeasurementCollection) String() string {
	result := ""
	for key, value := range c {
		result += fmt.Sprintf("%s:%f ", key, value)
	}

	return result
}

func (c MeasurementCollection) AddInt(key string, value int) {
	c[key] = Measurement(value)
}

func (c MeasurementCollection) AddFloat64(key string, value float64) {
	c[key] = Measurement(value)
}
