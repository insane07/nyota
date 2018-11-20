package sysutils

import (
	"os"
	"testing"
	"time"
)

func TestEnv(t *testing.T) {
	// // Getenv
	os.Setenv("var1", "val1")
	if v := Getenv("var1", ""); v != "val1" {
		t.Fatalf("Getenv got:%v exp:val1", v)
	}
	if v := Getenv("var2", ""); v != "" {
		t.Fatalf("Getenv(default empty) got:%v exp:val1", v)
	}
	if v := Getenv("var2", "val2"); v != "val2" {
		t.Fatalf("Getenv(default) got:%v exp:val2", v)
	}

	// GetenvInt
	os.Setenv("varInt1", "61")
	if v := GetenvInt("varInt1", 100); v != 61 {
		t.Fatalf("Getenv got:%v exp:61", v)
	}
	if v := GetenvInt("varInt2", 123); v != 123 {
		t.Fatalf("GetenvInt(default) got:%v exp:123", v)
	}

	// GetenvTime
	os.Setenv("varTime1", "30")
	if v := GetenvTime("varTime1", time.Minute, 20); v != 30*time.Minute {
		t.Fatalf("GetenvTime(default) got:%v exp:30min", v)
	}
	if v := GetenvTime("varTime2", time.Minute, 20); v != 20*time.Minute {
		t.Fatalf("GetenvTime(default) got:%v exp:20min", v)
	}
}
