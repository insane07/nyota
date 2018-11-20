package profiler

import "goprizm/network"

// const is not used since address of topic is required for confluent kafka library and golang
// does not allow to take address of const.
var (
	// ProfileReqTopic - kafka topic to enqueue profiling requests.
	ProfileReqTopic = "prizm.profiler.req"

	// ProfileReqTopic - kafka topic to enqueue profiling results.
	ProfileResultTopic = "prizm.profiler.result"
)

// ProfileReq - request generated and send by correlation service to initiate profiling.
type ProfileReq struct {
	TenantID    string              `json:"tenant_id"`          // tenant for which this req is generated
	RequestID   string              `json:"request_id"`         // request id for correlating across services.
	DeviceID    string              `json:"device_id"`          // unique id assigned to the endpoint
	MAC         string              `json:"mac"`                // mac address of endpoint
	Hostname    string              `json:"hostname,omitempty"` // endpoint hostname
	IP          string              `json:"ip,omitempty"`       // endpoint ip
	Device      network.Device      `json:"device,omitempty"`   // device to which this endpoint is currently classified.
	Fingerprint network.Fingerprint `json:"fingerprint"`        // device fingerprint info.
}

// ClassifyResult - classification result from one classifier.
type ClassifyResult struct {
	Device  network.Device `json:"device"`  // device to which endpoint is classified.
	Score   float64        `json:"score"`   // classification score.
	Explain interface{}    `json:"explain"` // per classifier explain
}

// Conflict - result of conflict detection.
type Conflict struct {
	On    bool           `json:"on"`    // true if conflict is detected.
	Other network.Device `json:"other"` // device to which this endpoint is initially classified to
}

type Explain struct {
	Classifiers map[string][]ClassifyResult `json:"classifiers,omitempty"` // result from different classifiers.
}

func (exp *Explain) Clone() Explain {
	if exp.Classifiers == nil {
		return Explain{}
	}

	clone := Explain{
		Classifiers: make(map[string][]ClassifyResult),
	}

	for classifier, results := range exp.Classifiers {
		r0 := make([]ClassifyResult, len(results))
		copy(r0, results)
		clone.Classifiers[classifier] = r0
	}

	return clone
}

// ProfileResult - result of profiling and clustering.
type ProfileResult struct {
	TenantID  string         `json:"tenant_id"`          // tenant for which profiling is done
	RequestID string         `json:"request_id"`         // request id for correlating across services.
	DeviceID  string         `json:"device_id"`          // unique id assigned to the endpoint
	MacVendor string         `json:"mac_vendor"`         // oui vendor.
	Device    network.Device `json:"device"`             // current profile to be used for conflict detection.
	Cluster   string         `json:"cluster,omitempty"`  // cluster to which device is assigned if no matching device is found.
	Conflict  Conflict       `json:"conflict,omitempty"` // result of conflict detection.
	Explain   Explain        `json:"explain,omitempty"`  // detailed explanation of how classification happened.
}
