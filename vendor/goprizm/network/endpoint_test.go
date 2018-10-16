package network

import (
	"reflect"
	"testing"

	"goprizm/gods"
	"goprizm/goprof"
	"goprizm/log"
)

func TestFingerprintMerge(t *testing.T) {
	fpMergeMax = 5

	fp, _ := NewFingerprintJSON(`{
		"dhcp": {
			"option55": ["1,2,3", "4,5,6"]
		},
		"host": {
			"user_agent" : ["ua1", "ua2", "ua2", "ua3"]
		}
	}`)

	fpOther, _ := NewFingerprintJSON(`{
		"dhcp": {
			"option55": ["1,2,3", "2,4,6"],
			"option60": ["abcd"]
		},
		"host": {
			"user_agent" : ["ua2", "ua4", "ua4", "ua6", "ua8"]
		},
		"snmp" : {
			"sys_descr": ["device0"]
		}
	}`)

	fpResult, _ := NewFingerprintJSON(`{
		"dhcp": {
			"option55": ["1,2,3", "2,4,6", "4,5,6"],
			"option60": ["abcd"]
		},
		"host": {
			"user_agent" : ["ua1", "ua2", "ua4", "ua6", "ua8"]
		},
		"snmp" : {
			"sys_descr": ["device0"]
		}
	}`)

	fp.Merge(fpOther)
	if !reflect.DeepEqual(fp, fpResult) {
		t.Fatalf("Fingerprint update failed fp:%+v\nexp:%+v", fp, fpResult)
	}
}

func TestFingerprintEquals(t *testing.T) {
	t.Run("EmptySourceFingerprint", func(t *testing.T) {
		fp1 := NewFingerprint()

		fp2, _ := NewFingerprintJSON(`{
			"dhcp": {
				"option55": ["1,2,3", "4,5,6"]
			}
		}`)

		if fp1.Equals(fp2) {
			t.Fatalf("expected fp1!=fp2")
		}
	})

	t.Run("CompareToEmptyFingerprint", func(t *testing.T) {
		fp1, _ := NewFingerprintJSON(`{
			"dhcp": {
				"option55": ["1,2,3", "4,5,6"]
			}
		}`)

		fp2 := NewFingerprint()

		if !fp1.Equals(fp2) {
			t.Fatalf("expected fp1==fp2")
		}
	})

	t.Run("CompareItslef", func(t *testing.T) {
		fp, _ := NewFingerprintJSON(`{
			"dhcp": {
				"option55": ["1,2,3", "4,5,6"]
			}
		}`)

		if !fp.Equals(fp) {
			t.Fatalf("expected fp==fp")
		}
	})

	t.Run("CompareSameFingerprint", func(t *testing.T) {
		fp1, _ := NewFingerprintJSON(`{
			"dhcp": {
				"option55": ["1,2,3", "4,5,6"]
			}
		}`)

		fp2, _ := NewFingerprintJSON(`{
			"dhcp": {
				"option55": ["1,2,3", "4,5,6"]
			}
		}`)

		if !fp1.Equals(fp2) {
			t.Fatalf("expected fp1==fp2")
		}
	})

	t.Run("CompareFPValue", func(t *testing.T) {
		fp1, _ := NewFingerprintJSON(`{
			"dhcp": {
				"option55": ["1,2,3", "4,5,6"]
			}
		}`)

		fp2, _ := NewFingerprintJSON(`{
			"dhcp": {
				"option55": ["1,2,3", "8,9,10"]
			}
		}`)

		if fp1.Equals(fp2) {
			t.Fatalf("expected fp1!=fp2")
		}
	})

	t.Run("CompareFPNamesapce", func(t *testing.T) {
		fp1, _ := NewFingerprintJSON(`{
			"dhcp": {
				"option55": ["1,2,3", "4,5,6"]
			}
		}`)

		fp2, _ := NewFingerprintJSON(`{
			"snmp": {
				"sys_descr": ["ARUBA"]
			}
		}`)

		if fp1.Equals(fp2) {
			t.Fatalf("expected fp1!=fp2")
		}
	})

	t.Run("CompareFPAttribute", func(t *testing.T) {
		fp1, _ := NewFingerprintJSON(`{
			"dhcp": {
				"option55": ["1,2,3", "4,5,6"]
			}
		}`)

		fp2, _ := NewFingerprintJSON(`{
			"dhcp": {
				"option55": ["1,2,3", "4,5,6"],
				"options": ["8,9,10"]
			}
		}`)

		if fp1.Equals(fp2) {
			t.Fatalf("expected fp1!=fp2")
		}
	})

	t.Run("CompareFPValueOrder", func(t *testing.T) {
		fp1, _ := NewFingerprintJSON(`{
			"dhcp": {
				"option55": ["1,2,3", "4,5,6"]
			}
		}`)

		fp2, _ := NewFingerprintJSON(`{
			"dhcp": {
				"option55": ["4,5,6", "1,2,3"]
			}
		}`)

		if !fp1.Equals(fp2) {
			t.Fatalf("expected fp1==fp2")
		}
	})

	t.Run("CompareFPNilValue", func(t *testing.T) {
		fp1, _ := NewFingerprintJSON(`{
			"dhcp": {
				"option55": ["1,2,3", "4,5,6"]
			}
		}`)

		fp2, _ := NewFingerprintJSON(`{
			"dhcp": {
				"option55": []
			}
		}`)

		if !fp1.Equals(fp2) {
			t.Fatalf("expected fp1==fp2")
		}
	})

	t.Run("CompareFPExtraValue", func(t *testing.T) {
		fp1, _ := NewFingerprintJSON(`{
			"dhcp": {
				"option55": ["1,2,3", "4,5,6"]
			}
		}`)

		fp2, _ := NewFingerprintJSON(`{
			"dhcp": {
				"option55": ["1,2,3", "6,7,8", "4,5,6"]
			}
		}`)

		if fp1.Equals(fp2) {
			t.Fatalf("expected fp1!=fp2")
		}
	})
}

func TestCleanFingerprints(t *testing.T) {
	fp := map[string]map[string][]string{
		"dhcp": map[string][]string{
			"option55": {"1,2,3"},
			"options":  nil,
		},
		"host": map[string][]string{
			"ports":      {"22", "34"},
			"user_agent": {},
		},
		"snmp": map[string][]string{
			"sys_descr": {},
		},
	}

	fpClean := map[string]map[string][]string{
		"dhcp": map[string][]string{
			"option55": {"1,2,3"},
		},
		"host": map[string][]string{
			"ports": {"22", "34"},
		},
	}

	Fingerprint(fp).Clean()
	if !reflect.DeepEqual(fp, fpClean) {
		t.Fatalf("Fingerprint clean fp:%+v != fpClean:%+v", fp, fpClean)
	}
}

func loadFps(t *testing.T, fpAttrPool *gods.StringPool) {
	fpJson := `{
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
}`
	recordMemUsage := func() {
		goprof.FreeMem()
		log.Printf("Memory Stats: %+v", goprof.ReadMemStats())
	}

	recordMemUsage()

	var fps []Fingerprint
	for i := 0; i < 200000; i++ {
		fp, err := NewFingerprintJSON(fpJson)
		if err != nil {
			t.Fatalf("load fingerprint: %+v", err)
		}

		if fpAttrPool != nil {
			fp.Intern(fpAttrPool)
		}

		fps = append(fps, fp)
	}
	recordMemUsage()
}

/*
go test -v -run TestFpMemoryUsage
=== RUN   TestFpMemoryUsage
2018/07/30 10:17:03 INFO Memory Stats: &{AllocKB:130 SysKB:2020 HeapAllocKB:130 HeapSysKB:672 HeapIdleKB:240 HeapInuseKB:432 HeapReleasedKB:208 HeapObjects:263 NextGCKB:4096 LastGC:2018-07-30 10:17:03 -0700 PDT}
2018/07/30 10:17:14 INFO Memory Stats: &{AllocKB:140 SysKB:1037611 HeapAllocKB:140 HeapSysKB:973888 HeapIdleKB:973368 HeapInuseKB:520 HeapReleasedKB:973336 HeapObjects:294 NextGCKB:4096 LastGC:2018-07-30 10:17:14 -0700 PDT}
--- PASS: TestFpMemoryUsage (10.83s)
PASS
ok  	goprizm/network	10.845s
*/
func TestFpMemoryUsage(t *testing.T) {
	loadFps(t, nil)
}

/*
go test -v -run TestFpInternMemoryUsage
=== RUN   TestFpInternMemoryUsage
2018/07/30 10:17:38 INFO Memory Stats: &{AllocKB:132 SysKB:2084 HeapAllocKB:132 HeapSysKB:608 HeapIdleKB:152 HeapInuseKB:456 HeapReleasedKB:120 HeapObjects:268 NextGCKB:4096 LastGC:2018-07-30 10:17:38 -0700 PDT}
2018/07/30 10:17:51 INFO Memory Stats: &{AllocKB:140 SysKB:846619 HeapAllocKB:140 HeapSysKB:792640 HeapIdleKB:792104 HeapInuseKB:536 HeapReleasedKB:792072 HeapObjects:294 NextGCKB:4096 LastGC:2018-07-30 10:17:51 -0700 PDT}
--- PASS: TestFpInternMemoryUsage (12.86s)
PASS
ok  	goprizm/network	12.873s
*/
func TestFpInternMemoryUsage(t *testing.T) {
	loadFps(t, gods.NewStringPool(1000000))
}
