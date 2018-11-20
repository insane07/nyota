package network

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
)

const (
	EndpointEventType      = "endpoint.attrs"
	ConnEventType          = "network.conn"
	DnsEventType           = "network.dns"
	PortMapEventType       = "network.portmap"
	FlowEventType          = "network.flow"
	ProfileResultEventType = "profile.result"

	// topic to read endpoint events from kafka sent by collector
	EndpointAttrTopic = "collector.endpoint.attr"

	// topic to read endpoint flows from kafka sent by collector
	EndpointFlowTopic = "collector.endpoint.flow"

	// topic to read endpoint from kafka emitted by engine after aggregation
	EndpointInfoTopic = "prizm.endpoint.info"
)

type EndpointEvent struct {
	MAC         string `json:"mac"`
	IP          string `json:"ip"`
	Hostname    string `json:"hostname"`
	NadIP       string `json:"nad_ip"`
	NadPort     string `json:"nad_port"`
	AP          string `json:"ap"`
	SSID        string `json:"ssid"`
	Fingerprint `json:"fingerprint"`
}

func (epEvent EndpointEvent) String() string {
	bytes, err := json.Marshal(epEvent)
	if err != nil {
		return "{}"
	}

	return string(bytes)
}

type ConnEvent struct {
	SrcIP  string   `json:"src_ip"`
	AppIds []string `json:"app_ids"`
}

type DNSEvent struct {
	SrcIP    string   `json:"src_ip"`
	Hostname string   `json:"hostname"`
	Ips      []string `json:"ips"`
}

type FlowEvent struct {
	SrcIP     string `json:"src_ip"`
	SrcPort   uint16 `json:"src_port"`
	DstIP     string `json:"dst_ip"`
	DstPort   uint16 `json:"dst_port"`
	Proto     uint8  `json:"proto"`      // proto value from IP header
	TxBytes   uint64 `json:"tx_bytes"`   // bytes transfered during time interval
	RxBytes   uint64 `json:"rx_bytes"`   // bytes transfered during time interval
	StartTime uint64 `json:"start_time"` // start time (in usec) sent by PPE
	EndTime   uint64 `json:"end_time"`   // end time (start time + duration)
	TxPackets uint64 `json:"tx_packets"`
	RxPackets uint64 `json:"rx_packets"`
	AppId     string `json:"app_id"`
	Version   int    `json:"ssl_version,omitempty"`
	Ciphers   string `json:"ssl_cipher_list,omitempty"`
	Server    string `json:"ssl_server_name,omitempty"`
}

var (
	evFieldSep = []byte(",")
)

// ParseEvent - format <timestamp>,<type>,<json>
func ParseEvent(data []byte) (evTime int64, evType string, ev interface{}, err error) {
	fields := bytes.SplitN(data, evFieldSep, 3)
	if len(fields) != 3 {
		err = fmt.Errorf("invalid event(expected 3 fields)")
		return
	}

	evTime, err = strconv.ParseInt(string(fields[0]), 10, 64)
	if err != nil {
		err = fmt.Errorf("timestamp not int(%v)", err)
		return
	}

	evType = string(fields[1])
	evBytes := fields[2]

	switch evType {
	case EndpointEventType:
		var epEvent EndpointEvent
		err = json.Unmarshal(evBytes, &epEvent)
		ev = epEvent
	case ConnEventType:
		var connEvent ConnEvent
		err = json.Unmarshal(evBytes, &connEvent)
		ev = connEvent
	case DnsEventType:
		var dnsEvent DNSEvent
		err = json.Unmarshal(evBytes, &dnsEvent)
		ev = dnsEvent
	case FlowEventType:
		var flowEvent FlowEvent
		err = json.Unmarshal(evBytes, &flowEvent)
		ev = flowEvent
	default:
		err = fmt.Errorf("invalid event type(%s)", evType)
		return
	}

	return
}
