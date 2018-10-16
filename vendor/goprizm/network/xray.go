package network

import "goprizm/gods"

// Traffic summary for a time window.
// TODO replace AvgXXX with quantiles(50, 90, 99)
type Traffic struct {
	Proto    string `json:"proto"`     // tcp/udp
	NumFlows uint64 `json:"num_flows"` // number of flows
	Rx       uint64 `json:"rx"`        // sum(rx bytes)
	RxPkts   uint64 `json:"rx_pkts"`   // num of rx packets
	Tx       uint64 `json:"tx"`        // sum(tx bytes)
	TxPkts   uint64 `json:"tx_pkts"`   // num of tx packets
	Duration uint64 `json:"duration"`  // flow duration in milliseconds
}

// Xrays are traffic stats + attributes build from flows.
// Attributes
//      host.dst_conns    = [host:port:proto...]
//      host.ports        = [22, ...]
// 		host.l7_protos    = [rtp,...]
//      host.app_ids
//      host.urls
//      host.ssl_ciphers
//      host.ssl_versions
//      host.cert_subject
//      host.cert_subject_alt_name
//      host.cert_issuer
type Xray struct {
	TenantID   string              `json:"tenant_id"`  // tenantID of endpoint
	DeviceID   string              `json:"device_id"`  // unique-id of device
	StartTime  int                 `json:"start_time"` // timestamp of first flow used
	EndTime    int                 `json:"end_time"`   // timestamp of last flow used
	To         map[string]*Traffic `json:"to"`         // traffic to each server host:port
	From       map[string]*Traffic `json:"from"`       // traffic from each client host:port
	Attributes map[string][]string `json:"attributes"` // fingerprint like attrs extracted.
}

// TrafficReport - stats build from all endpoints of same type(eg quantile)
type TrafficReport struct {
	NumFlows []gods.Percentile `json:"num_flows"` // quantile over all endpoints of same type
	Rx       []gods.Percentile `json:"rx"`
	RxPkts   []gods.Percentile `json:"rx_pkts`
	Tx       []gods.Percentile `json:"tx"`
	TxPkts   []gods.Percentile `json:"tx_pkts"`
	Duration []gods.Percentile `json:"duration"`
}

type XrayReport struct {
	To         TrafficReport       `json:"to"`
	From       TrafficReport       `json:"from"`
	Attributes map[string][]string `json:"attributes"`
}
