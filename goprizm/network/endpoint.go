package network

import (
	"encoding/json"
	"goprizm/gods"
	"sort"
)

var (
	DefaultDevice = Device{"Generic", "Generic", "Unclassified Device"}
)

// Device profile obtained after profiling.
type Device struct {
	Category string `json:"device_category"` // eg Computer, Printer etc     `db:"device_category"`
	Family   string `json:"device_family"`   // eg Windows, Epson
	Name     string `json:"device_name"`     // eg Windows 8, HTC Android
}

func (d Device) IsNil() bool {
	return (d.Category == "") || (d.Category == DefaultDevice.Category)
}

func (d *Device) Clean() {
	if d.IsNil() {
		*d = DefaultDevice
	}
}

// Fingerprint of one endpoint.
//
// Sources(namespaces) and attributes:
//    dhcp:
//      options:
//      option55:
//      option60:
//
//    snmp:
//      lldp_sys_descr:
//      cdp_cache_platform:
//      sys_descr:
//      hr_device_descr:
//      device_type:
//      name:
//
//    host:
//      os_type:
//      user_agent:
//      ports:
//      services:
//      device_type:
//
//    nmap:
//      device:
//
//    ssh:
//       device_name:
//
//    wmi:
//      os_name:

/***
	Example
 	=======
{
    "mac":"703eac1e10b4",
    "ip":"15.111.201.109",
    "hostname":"andrea",
    "fingerprint":{
        "dhcp":{
            "option55":[
                "1,3,6,15,44,46,47"
            ],
            "option60":[
                "1,3,6,15,44,46,47"
            ],
            "option":[
                "1,3,6,15,44,46,47"
            ]
        },
        "snmp":{
            "name":[
                "ArubaSwitch-10.2.51.244"
            ],
            "device_type":[
                "Switch"
            ],
            "sys_descr":[
                "ArubaOS (MODEL: ArubaS3500-24P-US), Version 7.4.1.5 (55591)"
            ],
            "lldp_sys_descr":[
                "ArubaOS (MODEL: ArubaS3500-24P-US), Version 7.4.1.5 (55591)"
            ]
        },
        "ssh":{
            "device_name":[
                "Aruba7220-US"
            ]
        },
        "wmi":{
            "os_name":[
                "Windows"
            ]
        },
        "tcp":{
            "device":[
                "ESX"
            ]
        },
        "host":{
            "user_agent":[
                "Opera/9.80 ( Linux armv6l; Opera TV Store/5581; (SonyBDP/BDV13)) Presto/2.12.362 Version/12.11"
            ],
            "ports":[
                "22",
                "111",
                "389",
                "3306"
            ],
            "services":[
                "22:ssh - OpenSSH Version: 5.3",
                "111:rpcbind Version:2-4",
                "389:ldap",
                "3306:mysql - MySQL Version: 5.1.73"
            ]
        }
    }
}
***/

type Fingerprint map[string]map[string][]string

func NewFingerprint() Fingerprint {
	return make(Fingerprint)
}

// NewFingerprintJSON creates Fingerprint from json. json could of one of the following
func NewFingerprintJSON(js string) (fp Fingerprint, err error) {
	err = json.Unmarshal([]byte(js), &fp)
	return
}

func (fp Fingerprint) Clone() Fingerprint {
	clone := make(Fingerprint)
	for ns, attrs := range fp {
		for attr, values := range attrs {
			cattrs, ok := clone[ns]
			if !ok {
				cattrs = make(map[string][]string)
				clone[ns] = cattrs
			}
			v0 := make([]string, len(values))
			copy(v0, values)
			cattrs[attr] = v0
		}
	}
	return clone
}

func (fp Fingerprint) IsEmpty() bool {
	return len(fp) <= 0
}

// Marshal Fingerprint to JSON
func (fp Fingerprint) ToJSON() (string, error) {
	bytes, err := json.Marshal(fp)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func (fp Fingerprint) String() string {
	fpString, err := fp.ToJSON()
	if err != nil {
		return ""
	}

	return fpString
}

func (fp Fingerprint) Add(source, field string, values ...string) {
	fields, ok := fp[source]
	if !ok {
		fields = make(map[string][]string)
		fp[source] = fields
	}

	fields[field] = append(fields[field], values...)
}

// replace the values
func (fp Fingerprint) Replace(source, field string, values []string) {
	fields, ok := fp[source]
	if !ok {
		fields = make(map[string][]string)
		fp[source] = fields
	}

	fields[field] = values
}

// delete the field
func (fp Fingerprint) Delete(source, field string) {
	fields, ok := fp[source]
	if !ok {
		return
	}

	delete(fields, field)
}

// Values return value of attribute(source.field) as slice of strings.
func (fp Fingerprint) Values(source, field string) []string {
	fields, ok := fp[source]
	if !ok {
		return nil
	}

	return fields[field]
}

// Value return value of attribute(source.field) as string.
func (fp Fingerprint) Value(source, field string) string {
	values := fp.Values(source, field)
	if len(values) > 0 {
		return values[0]
	}
	return ""
}

var (
	fpMergeMax = 10
)

// Merge merge fingerprint with other fingerprint
func (fp Fingerprint) Merge(other Fingerprint) bool {
	modified := false
	for ns, attrs := range other {
		fpAttrs, ok := fp[ns]

		// ns does not exist, add all attr-values.
		if !ok {
			fp[ns] = attrs
			modified = true
			continue
		}

		// ns attr-values exist, merge values.
		for attr, values := range attrs {
			var mod bool
			fpAttrs[attr], mod = fpMerge(values, fpAttrs[attr], fpMergeMax)
			modified = modified || mod
		}
	}

	return modified
}

func fpMerge(values, valuesBase []string, max int) (result []string, modified bool) {
	uniqueVals := gods.NewStringSet(values...)

	if len(valuesBase) == 0 {
		modified = len(values) > 0
	} else {
		for _, v := range valuesBase {
			if uniqueVals.Len() >= max {
				break
			}

			isNewVal := uniqueVals.Add(v)
			modified = modified || isNewVal
		}
	}

	values = uniqueVals.ToSlice()

	if len(values) > max {
		values = values[:max]
	}

	return values, modified
}

func (fp Fingerprint) MacVendor() string {
	return fp.Value("host", "mac_vendor")
}

// Return true if Fingerprint has fields from given source.
func (fp Fingerprint) HasSource(source string) bool {
	_, ok := fp[source]
	return ok
}

func (fp Fingerprint) Device() Device {
	s, ok := fp["device"]
	if !ok {
		return Device{}
	}

	dev := Device{}

	var values []string
	if values, ok = s["category"]; !ok || len(values) <= 0 {
		return Device{}
	}
	dev.Category = values[0]

	if values, ok = s["family"]; !ok || len(values) <= 0 {
		return Device{}
	}
	dev.Family = values[0]

	if values, ok = s["name"]; !ok || len(values) <= 0 {
		return Device{}
	}
	dev.Name = values[0]

	return dev
}

// Equals returns bool comparing with other fingerprint.
// Returns false if other fingerprint conatines extra sources or attribute values.
// Source has namespaces like DHCP, SNMP
// Order of attribute values is ignored to accomodate user agent, services and port attribute types where order is not fixed.
func (fp Fingerprint) Equals(other Fingerprint) bool {
	// Compare source
	for source, attrs := range other {
		// Source namespace does not exist Ex: DHCP, SNMP
		if _, ok := fp[source]; !ok {
			return false
		}

		// Compare the attributes for source
		for attr, values := range attrs {
			// ignore empty values
			if len(values) == 0 {
				continue
			}

			// Check existence of attribute in particular namespace
			fpValues, ok := fp[source][attr]
			if !ok || len(fpValues) == 0 {
				return false
			}

			// other fingerprint contains more values for attribute.
			if len(values) > len(fpValues) {
				return false
			}

			// ignore ordering by sorting the copied values
			otherValuesOrdered := append([]string{}, values[:]...)
			fpValuesOrdered := append([]string{}, fpValues[:]...)

			sort.Strings(otherValuesOrdered)
			sort.Strings(fpValuesOrdered)

			fpValuesOrdered = fpValuesOrdered[:len(fpValuesOrdered)] // avoid bound checking

			for idx, otherVal := range otherValuesOrdered {
				// values do not match
				if otherVal != fpValuesOrdered[idx] {
					return false
				}
			}
		}
	}
	return true
}

func (fp Fingerprint) Clean() {
	for ns, attrs := range fp {
		for attr, values := range attrs {
			if len(values) == 0 {
				delete(attrs, attr)
			}

			if len(values) > 1 {
				attrs[attr] = gods.NewStringSet(values...).ToSlice()
			}
		}
		if len(attrs) == 0 {
			delete(fp, ns)
		}
	}
}

func (fp Fingerprint) Intern(attrPool *gods.StringPool) {
	for ns, attrs := range fp {
		for attr, values := range attrs {
			var internValues []string
			for _, value := range values {
				internValues = append(internValues, attrPool.Get(value))
			}

			fp[ns][attrPool.Get(attr)] = internValues
		}

		fp[attrPool.Get(ns)] = attrs
	}
}
