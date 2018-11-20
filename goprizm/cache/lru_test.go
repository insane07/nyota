// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cache

import (
	"reflect"
	"testing"
	"time"
)

func TestLRU(t *testing.T) {
	c := NewLRU(2, time.Hour)

	expectMiss := func(k string) {
		v, ok := c.Get(k)
		if ok {
			t.Fatalf("expected cache miss on key %q but hit value %v", k, v)
		}
	}

	expectHit := func(k string, ev interface{}) {
		v, ok := c.Get(k)
		if !ok {
			t.Fatalf("expected cache(%q)=%v; but missed", k, ev)
		}
		if !reflect.DeepEqual(v, ev) {
			t.Fatalf("expected cache(%q)=%v; but got %v", k, ev, v)
		}
	}

	expectMiss("1")
	c.Add("1", "one")
	expectHit("1", "one")

	c.Add("2", "two")
	expectHit("1", "one")
	expectHit("2", "two")

	c.Add("3", "three")
	expectHit("3", "three")
	expectHit("2", "two")
	expectMiss("1")
}

func TestRemoveOldest(t *testing.T) {
	c := NewLRU(2, time.Hour)
	c.Add("1", "one")
	c.Add("2", "two")
	if k, v := c.removeOldest(); k != "1" || v != "one" {
		t.Fatalf("oldest = %q, %q; want 1, one", k, v)
	}
	if k, v := c.removeOldest(); k != "2" || v != "two" {
		t.Fatalf("oldest = %q, %q; want 2, two", k, v)
	}
	if k, v := c.removeOldest(); k != nil || v != nil {
		t.Fatalf("oldest = %v, %v; want \"\", nil", k, v)
	}
}

func TestTTL(t *testing.T) {
	c := NewLRU(2, 500*time.Millisecond)
	c.Add("1", "one")
	time.Sleep(time.Second)
	c.Add("2", "two")

	if v, ok := c.Get("1"); ok {
		t.Fatalf("ttl failed v:%v", v)
	}
	if v, ok := c.Get("2"); !ok || v.(string) != "two" {
		t.Fatalf("ttl valid entry expired v:%v", v)
	}
}
