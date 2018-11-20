package netutils

import (
	"fmt"
	"testing"
	"time"
)

const (
	expectedMac = "001122aabbcc"
)

func TestNormalizeMac(t *testing.T) {
	var (
		m   string
		err error
	)

	m, err = NormalizeMAC(" 00-11-22-AA-BB-CC ")
	if err != nil || m != expectedMac {
		t.Fatal("Normalize mac failed for seperator=`-`")
	}

	m, err = NormalizeMAC(" 0011:22AA:BB-CC ")
	if err != nil || m != expectedMac {
		t.Fatal("Normalize mac failed for seperator=`:`")
	}

	m, err = NormalizeMAC(" 0011.22AA.BBCC ")
	if err != nil || m != expectedMac {
		t.Fatal("Normalize mac failed for seperator=`.`")
	}

	m, err = NormalizeMAC(" 001122AABBCC ")
	if err != nil || m != expectedMac {
		t.Fatal("Normalize mac failed for no separator")
	}
}

func TestIsIPAddr(t *testing.T) {
	if !IsIPAddr("10.17.4.11") {
		t.Fatalf("IPv4 failed")
	}

	if !IsIPAddr("2001:4860:0:2001::68") {
		t.Fatalf("IPv4 failed")
	}

	if IsIPAddr("1017.4.11") {
		t.Fatalf("Invalid IP succeeded")
	}
}

func TestRetryOp(t *testing.T) {
	i := 0
	err := RetryOp(3, 100*time.Millisecond, "test1", func() error {
		i += 1
		if i == 3 {
			return nil
		}
		return fmt.Errorf("test1: sim error")
	})
	if err != nil {
		t.Fatalf("RetryOp failed err:%v", err)
	}

	err = RetryOp(3, 100*time.Millisecond, "test2", func() error {
		return nil
	})
	if err != nil {
		t.Fatalf("RetryOp failed err:%v", err)
	}

	err = RetryOp(3, 100*time.Millisecond, "test3", func() error {
		return fmt.Errorf("test3: sim error")
	})
	if err == nil {
		t.Fatalf("RetryOp failed err:%v", err)
	}
}
