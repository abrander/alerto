package icmpping

import (
	"net"
	"time"

	"github.com/abrander/alerto/agent"
)

func init() {
	agent.Register("dns", NewDns)
}

func NewDns() agent.Agent {
	return new(Dns)
}

type (
	Dns struct {
		Target     string `json:"target" description:"The name to resolve"`
		RecordType string `json:"recordType" description:"The record type to lookup" enum:"A,AAAA,A*"`
	}
)

func GetIP(hostname string) ([]net.IP, error) {
	list, err := net.LookupIP(hostname)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func GetIPv4(hostname string) ([]net.IP, error) {
	list, err := GetIP(hostname)
	if err != nil {
		return nil, err
	}

	list4 := []net.IP{}

	for _, ip := range list {
		if ip.To4() != nil {
			list4 = append(list4, ip)
		}
	}

	return list4, nil
}

func GetIPv6(hostname string) ([]net.IP, error) {
	list, err := GetIP(hostname)
	if err != nil {
		return nil, err
	}

	list6 := []net.IP{}

	for _, ip := range list {
		if ip.To4() == nil {
			list6 = append(list6, ip)
		}
	}

	return list6, nil
}

func (i *Dns) Execute(request agent.Request) agent.Result {
	entries := []net.IP{}

	start := time.Now()

	var err error

	switch i.RecordType {
	case "":
		fallthrough
	case "A*":
		entries, err = GetIP(i.Target)
	case "A":
		entries, err = GetIPv4(i.Target)
	case "AAAA":
		entries, err = GetIPv6(i.Target)
	default:
		return agent.NewResult(agent.Failed, nil, "method '%s' not supported", i.RecordType)
	}

	if err != nil {
		return agent.NewResult(agent.Failed, agent.NewMeasurementCollection("time", time.Now().Sub(start)), err.Error())
	}

	if len(entries) > 0 {
		return agent.NewResult(agent.Ok, agent.NewMeasurementCollection("time", time.Now().Sub(start)), "%d addresses", len(entries))
	} else {
		return agent.NewResult(agent.Failed, agent.NewMeasurementCollection("time", time.Now().Sub(start)), "no addresses")
	}
}
