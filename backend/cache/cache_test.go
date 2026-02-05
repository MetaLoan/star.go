package cache

import (
	"star/core"
	"testing"
	"time"
)

func TestMemoryCache(t *testing.T) {
	mc := NewMemoryCache(100)
	key := "test_key"
	slot := &core.TimeSlot{UserID: "user1"}

	// Test Set and Get
	mc.Set(key, slot, 1*time.Minute)
	val, found := mc.Get(key)
	if !found {
		t.Fatal("Expected to find key in cache")
	}
	if val.UserID != "user1" {
		t.Errorf("Expected UserID user1, got %s", val.UserID)
	}

	// Test Expiration
	mc.Set("expired", slot, -1*time.Second)
	_, found = mc.Get("expired")
	if found {
		t.Error("Expected key to be expired")
	}

	// Test Delete
	mc.Delete(key)
	_, found = mc.Get(key)
	if found {
		t.Error("Expected key to be deleted")
	}
}

func TestMultiLevelCache(t *testing.T) {
	l1 := NewMemoryCache(10)
	l2 := NewMemoryCache(10) // Simulate L2 with another memory cache
	mc := NewMultiLevelCache(l1, l2)

	key := "multi_key"
	slot := &core.TimeSlot{UserID: "user2"}

	// Set in L2 only
	l2.Set(key, slot, 1*time.Minute)

	// Get should fetch from L2 and promote to L1
	val, found := mc.Get(key)
	if !found {
		t.Fatal("Expected to find key in multi-level cache")
	}
	if val.UserID != "user2" {
		t.Errorf("Expected UserID user2, got %s", val.UserID)
	}

	// Check L1 promotion
	_, foundInL1 := l1.Get(key)
	if !foundInL1 {
		t.Error("Expected key to be promoted to L1")
	}
}
