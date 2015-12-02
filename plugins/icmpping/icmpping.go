package icmpping

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"sync"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"

	"github.com/abrander/alerto/logger"
	"github.com/abrander/alerto/plugins"
)

type (
	IcmpPing struct {
		Target  string `json:"target" description:"The IPv4 hostname or IP address to ping"`
		id      int
		seq     int
		payload []byte
	}

	IcmpReply struct {
		Source string
		Status Status
	}

	Status int
)

func init() {
	var err error
	conn, err = icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		if err.Error() == "listen ip4:icmp 0.0.0.0: socket: operation not permitted" {
			logger.Error("icmpping", "Please run:\nsudo setcap cap_net_raw=ep %s\n", os.Args[0])
		} else {
			logger.Error("icmpping", err.Error())
		}

		os.Exit(1)
	}

	active = make(map[uint16]chan IcmpReply)

	plugins.Register("icmp4", NewIcmpPing)

	go ListenLoop()
}

func NewIcmpPing() plugins.Plugin {
	return new(IcmpPing)
}

const (
	Reply Status = iota
	Unreachable
)

var conn *icmp.PacketConn

// Maybe we should use more than just 'id'..?
var active map[uint16]chan IcmpReply
var activeLock sync.RWMutex

func decodeUnreachable(packetData []byte) (id uint16, seq uint16) {
	packet := gopacket.NewPacket(packetData, layers.LayerTypeIPv4, gopacket.NoCopy)

	ll := packet.Layers()

	if len(ll) < 2 {
		return 0, 0
	}

	if ll[1].LayerType() == layers.LayerTypeICMPv4 {
		icmp := ll[1].(*layers.ICMPv4)

		return icmp.Id, icmp.Seq
	}

	return 0, 0
}

func ListenLoop() {
	// Set maximum packet size to 9000 to support jumbo frames
	readBytes := make([]byte, 9000)

	for {
		_, peer, err := conn.ReadFrom(readBytes)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		packet := gopacket.NewPacket(readBytes, layers.LayerTypeICMPv4, gopacket.NoCopy)

		ll := packet.Layers()

		if len(ll) > 0 && ll[0].LayerType() == layers.LayerTypeICMPv4 {
			icmp := ll[0].(*layers.ICMPv4)
			typ := uint8(icmp.TypeCode >> 8)

			if typ == layers.ICMPv4TypeDestinationUnreachable {
				if len(ll) > 1 {
					id, _ := decodeUnreachable(ll[1].LayerContents())

					activeLock.RLock()
					ch, found := active[id]
					activeLock.RUnlock()
					if found {
						ch <- IcmpReply{Source: peer.String(), Status: Unreachable}
						continue
					}
				}
			} else if typ == layers.ICMPv4TypeEchoReply {
				activeLock.RLock()
				ch, found := active[icmp.Id]
				activeLock.RUnlock()
				if found {
					ch <- IcmpReply{Source: peer.String(), Status: Reply}
					continue
				}
			}
		}
	}
}

func (i *IcmpPing) Run(transport plugins.Transport, request plugins.Request) plugins.Result {
	ra, err := net.ResolveIPAddr("ip4:icmp", i.Target)

	if err != nil {
		return plugins.NewResult(plugins.Failed, nil, err.Error())
	}

	i.id = rand.Intn(0xffff)
	i.seq = rand.Intn(0xffff)
	i.payload = []byte("alerto pinger")

	bytes, err := (&icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			ID: i.id, Seq: i.seq,
			Data: i.payload,
		},
	}).Marshal(nil)

	replyChannel := make(chan IcmpReply)

	activeLock.Lock()
	_, found := active[uint16(i.id)]
	for found {
		// NOTE: This becomes subject to the birthday paradox at some point
		i.id = rand.Intn(0xffff)

		_, found = active[uint16(i.id)]
	}
	active[uint16(i.id)] = replyChannel
	activeLock.Unlock()

	if n, err := conn.WriteTo(bytes, ra); err != nil {
		return plugins.NewResult(plugins.Failed, nil, err.Error())
	} else if n != len(bytes) {
		return plugins.NewResult(plugins.Failed, nil, "sent %d bytes; wanted %d", n, len(bytes))
	}

	start := time.Now()
	c := time.After(request.Timeout)

	for {
		select {
		case <-c:
			return plugins.NewResult(plugins.Failed, plugins.NewMeasurementCollection("time", time.Now().Sub(start)), "timeout [%s]", i.Target)

		case reply := <-replyChannel:
			activeLock.Lock()
			delete(active, uint16(i.id))
			activeLock.Unlock()

			switch reply.Status {
			case Reply:
				return plugins.NewResult(plugins.Ok, plugins.NewMeasurementCollection("time", time.Now().Sub(start)), "reply from %s [%s]", reply.Source, i.Target)
			case Unreachable:
				return plugins.NewResult(plugins.Failed, plugins.NewMeasurementCollection("time", time.Now().Sub(start)), "unreachable from %s [%s]", reply.Source, i.Target)
			}
		}
	}
}

// Ensure compliance
var _ plugins.Agent = (*IcmpPing)(nil)
