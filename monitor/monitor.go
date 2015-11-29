package monitor

import (
	"errors"
	"math/rand"
	"os"
	"sync"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/abrander/alerto/agent"
	"github.com/abrander/alerto/logger"
)

type (
	Monitor struct {
		Id         bson.ObjectId `json:"id" bson:"_id"`
		Interval   time.Duration `json:"interval"`
		Agent      agent.Job     `json:"agent"`
		LastCheck  time.Time     `json:"lastCheck"`
		NextCheck  time.Time     `json:"nextCheck"`
		LastResult agent.Result  `json:"lastResult"`
	}
)

var (
	sess       *mgo.Session
	db         *mgo.Database
	collection *mgo.Collection

	ErrorInvalidId error = errors.New("Invalid id")
)

func init() {
	sess, err := mgo.Dial("127.0.0.1")
	if err != nil {
		logger.Error("monitor", "Can't connect to mongo, go error %v\n", err)
		os.Exit(1)
	}

	db = sess.DB("alerto")
	collection = db.C("monitors")
}

func GetAllMonitors() []Monitor {
	var monitors []Monitor

	err := collection.Find(bson.M{}).All(&monitors)
	if err != nil {
		logger.Red("monitor", "Error getting monitors from Mongo: %s", err.Error())
	}

	return monitors
}

func GetMonitor(id string) (Monitor, error) {
	var monitor Monitor

	if !bson.IsObjectIdHex(id) {
		return monitor, ErrorInvalidId
	}

	err := collection.FindId(bson.ObjectIdHex(id)).One(&monitor)
	if err != nil {
		logger.Red("monitor", "Error getting monitors from Mongo: %s", err.Error())
		return monitor, err
	}

	return monitor, nil
}

func UpdateMonitor(mon *Monitor) error {
	return collection.UpdateId(mon.Id, mon)
}

func AddMonitor(mon *Monitor) error {
	mon.Id = bson.NewObjectId()

	return collection.Insert(mon)
}

func DeleteMonitor(mon *Monitor) error {
	return collection.RemoveId(mon.Id)
}

func Loop(wg sync.WaitGroup) {
	ticker := time.Tick(time.Millisecond * 100)

	inFlight := make(map[bson.ObjectId]bool)
	inFlightLock := sync.RWMutex{}
	for t := range ticker {
		var monitors []Monitor
		err := collection.Find(bson.M{}).All(&monitors)
		if err != nil {
			logger.Red("monitor", "Error getting monitors from Mongo: %s", err.Error())
			continue
		}

		for _, mon := range monitors {
			age := t.Sub(mon.LastCheck)  // positive: past
			wait := mon.NextCheck.Sub(t) // positive: future

			inFlightLock.RLock()
			_, found := inFlight[mon.Id]
			inFlightLock.RUnlock()

			if found {
				// skipping monitors in flight
			} else if age > mon.Interval*2 && wait < -mon.Interval {
				checkIn := time.Duration(rand.Int63n(int64(mon.Interval)))
				mon.NextCheck = t.Add(checkIn)
				logger.Yellow("monitor", "%s %s: Delaying first check by %s", mon.Id.Hex(), mon.Agent.AgentId, checkIn)

				err = collection.UpdateId(mon.Id, mon)
				if err != nil {
					logger.Red("Error updating: %v", err.Error())
				}
			} else if wait < 0 {
				inFlightLock.Lock()
				inFlight[mon.Id] = true
				inFlightLock.Unlock()

				go func(mon Monitor) {
					r := mon.Agent.Run()
					if r.Status == agent.Ok {
						logger.Green("monitor", "%s %s: %s [%s]: %s", mon.Id.Hex(), mon.Agent.AgentId, r.Text, r.Duration, r.Measurements)
					} else {
						logger.Red("monitor", "%s %s: %s [%s]", mon.Id.Hex(), mon.Agent.AgentId, r.Text, r.Duration)
					}
					mon.LastResult = r
					mon.LastCheck = t
					mon.NextCheck = t.Add(mon.Interval)

					err = collection.UpdateId(mon.Id, mon)
					if err != nil {
						logger.Red("monitor", "Error updating: %s\n", err.Error())
					}
					inFlightLock.Lock()
					delete(inFlight, mon.Id)
					inFlightLock.Unlock()
				}(mon)
			}
		}
	}

	wg.Done()
}
