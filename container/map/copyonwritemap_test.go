package maps

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"
)

func TestNewCopyOnWriteMap(t *testing.T) {
	m := NewCopyOnWriteMap[string, int]()
	if m == nil {
		t.Fatal("NewCopyOnWriteMap returned nil")
	}
	if !m.IsEmpty() {
		t.Error("New map should be empty")
	}
	if m.Size() != 0 {
		t.Error("New map size should be 0")
	}
}

func TestNewCopyOnWriteMapWithCapacity(t *testing.T) {
	m := NewCopyOnWriteMapWithCapacity[string, int](10)
	if m == nil {
		t.Fatal("NewCopyOnWriteMapWithCapacity returned nil")
	}
	if !m.IsEmpty() {
		t.Error("New map should be empty")
	}
}

func TestCopyOnWriteMapFromMap(t *testing.T) {
	source := map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
	}
	m := CopyOnWriteMapFromMap(source)
	if m.Size() != 3 {
		t.Errorf("Expected size 3, got %d", m.Size())
	}
	if val, ok := m.Get("two"); !ok || val != 2 {
		t.Error("Failed to get value from copied map")
	}
}

func TestCopyOnWriteMapPutAndGet(t *testing.T) {
	m := NewCopyOnWriteMap[string, int]()

	// Test Put
	oldVal, existed := m.Put("key1", 100)
	if existed {
		t.Error("Key should not exist initially")
	}
	if oldVal != 0 {
		t.Error("Old value should be zero for new key")
	}

	// Test Get
	val, ok := m.Get("key1")
	if !ok {
		t.Error("Key should exist after Put")
	}
	if val != 100 {
		t.Errorf("Expected value 100, got %d", val)
	}

	// Test Put existing key
	oldVal, existed = m.Put("key1", 200)
	if !existed {
		t.Error("Key should exist")
	}
	if oldVal != 100 {
		t.Errorf("Expected old value 100, got %d", oldVal)
	}

	// Verify new value
	val, ok = m.Get("key1")
	if !ok || val != 200 {
		t.Error("Failed to update existing key")
	}
}

func TestCopyOnWriteMapRemove(t *testing.T) {
	m := NewCopyOnWriteMap[string, int]()
	m.Put("key1", 100)
	m.Put("key2", 200)

	// Remove existing key
	oldVal, existed := m.Remove("key1")
	if !existed {
		t.Error("Key should exist")
	}
	if oldVal != 100 {
		t.Errorf("Expected old value 100, got %d", oldVal)
	}

	// Verify removal
	_, ok := m.Get("key1")
	if ok {
		t.Error("Key should not exist after removal")
	}

	// Remove non-existing key
	_, existed = m.Remove("nonexistent")
	if existed {
		t.Error("Non-existing key should not be found")
	}
}

func TestCopyOnWriteMapContains(t *testing.T) {
	m := NewCopyOnWriteMap[string, int]()
	m.Put("key1", 100)
	m.Put("key2", 200)

	// Test ContainsKey
	if !m.ContainsKey("key1") {
		t.Error("Map should contain key1")
	}
	if m.ContainsKey("nonexistent") {
		t.Error("Map should not contain nonexistent key")
	}

	// Test ContainsValue
	if !m.ContainsValue(100) {
		t.Error("Map should contain value 100")
	}
	if m.ContainsValue(999) {
		t.Error("Map should not contain value 999")
	}
}

func TestCopyOnWriteMapSizeAndEmpty(t *testing.T) {
	m := NewCopyOnWriteMap[string, int]()

	if !m.IsEmpty() {
		t.Error("New map should be empty")
	}
	if m.Size() != 0 {
		t.Error("New map size should be 0")
	}

	m.Put("key1", 100)
	if m.IsEmpty() {
		t.Error("Map should not be empty after adding element")
	}
	if m.Size() != 1 {
		t.Errorf("Expected size 1, got %d", m.Size())
	}

	m.Put("key2", 200)
	if m.Size() != 2 {
		t.Errorf("Expected size 2, got %d", m.Size())
	}
}

func TestCopyOnWriteMapClear(t *testing.T) {
	m := NewCopyOnWriteMap[string, int]()
	m.Put("key1", 100)
	m.Put("key2", 200)

	m.Clear()
	if !m.IsEmpty() {
		t.Error("Map should be empty after Clear")
	}
	if m.Size() != 0 {
		t.Error("Map size should be 0 after Clear")
	}
}

func TestCopyOnWriteMapKeysAndValues(t *testing.T) {
	m := NewCopyOnWriteMap[string, int]()
	m.Put("key1", 100)
	m.Put("key2", 200)
	m.Put("key3", 300)

	keys := m.Keys()
	if len(keys) != 3 {
		t.Errorf("Expected 3 keys, got %d", len(keys))
	}

	values := m.Values()
	if len(values) != 3 {
		t.Errorf("Expected 3 values, got %d", len(values))
	}

	// Check if all expected keys are present
	keyMap := make(map[string]bool)
	for _, key := range keys {
		keyMap[key] = true
	}
	expectedKeys := []string{"key1", "key2", "key3"}
	for _, expectedKey := range expectedKeys {
		if !keyMap[expectedKey] {
			t.Errorf("Expected key %s not found", expectedKey)
		}
	}
}

func TestCopyOnWriteMapForEach(t *testing.T) {
	m := NewCopyOnWriteMap[string, int]()
	m.Put("key1", 100)
	m.Put("key2", 200)
	m.Put("key3", 300)

	sum := 0
	count := 0
	m.ForEach(func(k string, v int) {
		sum += v
		count++
	})

	if count != 3 {
		t.Errorf("Expected to iterate over 3 elements, got %d", count)
	}
	if sum != 600 {
		t.Errorf("Expected sum 600, got %d", sum)
	}
}

func TestCopyOnWriteMapString(t *testing.T) {
	m := NewCopyOnWriteMap[string, int]()

	// Empty map
	str := m.String()
	if str != "{}" {
		t.Errorf("Expected '{}', got '%s'", str)
	}

	// Non-empty map
	m.Put("key1", 100)
	str = m.String()
	if str != "{key1=100}" {
		t.Errorf("Expected '{key1=100}', got '%s'", str)
	}
}

func TestCopyOnWriteMapPutAll(t *testing.T) {
	m1 := NewCopyOnWriteMap[string, int]()
	m1.Put("key1", 100)

	m2 := NewCopyOnWriteMap[string, int]()
	m2.Put("key2", 200)
	m2.Put("key3", 300)

	m1.PutAll(m2)

	if m1.Size() != 3 {
		t.Errorf("Expected size 3, got %d", m1.Size())
	}

	if val, ok := m1.Get("key2"); !ok || val != 200 {
		t.Error("Failed to copy key2 from m2")
	}
}

func TestCopyOnWriteMapPutAllFromMap(t *testing.T) {
	m := NewCopyOnWriteMap[string, int]()
	m.Put("key1", 100)

	source := map[string]int{
		"key2": 200,
		"key3": 300,
	}

	m.PutAllFromMap(source)

	if m.Size() != 3 {
		t.Errorf("Expected size 3, got %d", m.Size())
	}

	if val, ok := m.Get("key2"); !ok || val != 200 {
		t.Error("Failed to copy key2 from source map")
	}
}

func TestCopyOnWriteMapToMap(t *testing.T) {
	m := NewCopyOnWriteMap[string, int]()
	m.Put("key1", 100)
	m.Put("key2", 200)

	goMap := m.ToMap()
	if len(goMap) != 2 {
		t.Errorf("Expected map length 2, got %d", len(goMap))
	}

	if goMap["key1"] != 100 {
		t.Error("Failed to convert to Go map")
	}
}

func TestCopyOnWriteMapSnapshot(t *testing.T) {
	m := NewCopyOnWriteMap[string, int]()
	m.Put("key1", 100)
	m.Put("key2", 200)

	snapshot := m.Snapshot()
	if len(snapshot) != 2 {
		t.Errorf("Expected snapshot length 2, got %d", len(snapshot))
	}

	// Modify original map
	m.Put("key3", 300)

	// Snapshot should remain unchanged
	if len(snapshot) != 2 {
		t.Error("Snapshot should not be affected by subsequent modifications")
	}
}

func TestCopyOnWriteMapPutIfAbsent(t *testing.T) {
	m := NewCopyOnWriteMap[string, int]()

	// Put new key
	oldVal, inserted := m.PutIfAbsent("key1", 100)
	if !inserted {
		t.Error("Should have inserted new key")
	}
	if oldVal != 0 {
		t.Error("Old value should be zero for new key")
	}

	// Try to put existing key
	oldVal, inserted = m.PutIfAbsent("key1", 200)
	if inserted {
		t.Error("Should not have inserted existing key")
	}
	if oldVal != 100 {
		t.Errorf("Expected old value 100, got %d", oldVal)
	}

	// Verify value unchanged
	if val, ok := m.Get("key1"); !ok || val != 100 {
		t.Error("Value should remain unchanged")
	}
}

func TestCopyOnWriteMapReplace(t *testing.T) {
	m := NewCopyOnWriteMap[string, int]()
	m.Put("key1", 100)

	// Replace existing key
	oldVal, replaced := m.Replace("key1", 200)
	if !replaced {
		t.Error("Should have replaced existing key")
	}
	if oldVal != 100 {
		t.Errorf("Expected old value 100, got %d", oldVal)
	}

	// Try to replace non-existing key
	_, replaced = m.Replace("nonexistent", 300)
	if replaced {
		t.Error("Should not have replaced non-existing key")
	}
}

func TestCopyOnWriteMapReplaceIf(t *testing.T) {
	m := NewCopyOnWriteMap[string, int]()
	m.Put("key1", 100)

	// Replace with correct old value
	replaced := m.ReplaceIf("key1", 100, 200)
	if !replaced {
		t.Error("Should have replaced with correct old value")
	}

	// Try to replace with incorrect old value
	replaced = m.ReplaceIf("key1", 100, 300)
	if replaced {
		t.Error("Should not have replaced with incorrect old value")
	}

	// Verify current value
	if val, ok := m.Get("key1"); !ok || val != 200 {
		t.Error("Value should be 200")
	}
}

// Concurrent tests
func TestCopyOnWriteMapConcurrentReads(t *testing.T) {
	m := NewCopyOnWriteMap[string, int]()

	// Pre-populate data
	for i := 0; i < 100; i++ {
		m.Put(fmt.Sprintf("key%d", i), i)
	}

	const numReaders = 10
	const readsPerReader = 1000

	var wg sync.WaitGroup
	wg.Add(numReaders)

	// Start multiple concurrent readers
	for i := 0; i < numReaders; i++ {
		go func(readerID int) {
			defer wg.Done()
			for j := 0; j < readsPerReader; j++ {
				key := fmt.Sprintf("key%d", j%100)
				if val, ok := m.Get(key); !ok || val != j%100 {
					t.Errorf("Reader %d: unexpected value for %s", readerID, key)
					return
				}
			}
		}(i)
	}

	wg.Wait()
}

func TestCopyOnWriteMapConcurrentReadWrite(t *testing.T) {
	m := NewCopyOnWriteMap[string, int]()

	const numReaders = 5
	const numWriters = 2
	const operations = 100

	var wg sync.WaitGroup
	wg.Add(numReaders + numWriters)

	// Start readers
	for i := 0; i < numReaders; i++ {
		go func(readerID int) {
			defer wg.Done()
			for j := 0; j < operations; j++ {
				key := fmt.Sprintf("key%d", j%10)
				m.Get(key) // Don't check result, just ensure no panic
				runtime.Gosched()
			}
		}(i)
	}

	// Start writers
	for i := 0; i < numWriters; i++ {
		go func(writerID int) {
			defer wg.Done()
			for j := 0; j < operations; j++ {
				key := fmt.Sprintf("key%d", j%10)
				value := writerID*1000 + j
				m.Put(key, value)
				runtime.Gosched()
			}
		}(i)
	}

	wg.Wait()
}

func TestCopyOnWriteMapConcurrentWrites(t *testing.T) {
	m := NewCopyOnWriteMap[string, int]()

	const numWriters = 5
	const writesPerWriter = 100

	var wg sync.WaitGroup
	wg.Add(numWriters)

	// Start multiple concurrent writers
	for i := 0; i < numWriters; i++ {
		go func(writerID int) {
			defer wg.Done()
			for j := 0; j < writesPerWriter; j++ {
				key := fmt.Sprintf("writer%d_key%d", writerID, j)
				value := writerID*1000 + j
				m.Put(key, value)
			}
		}(i)
	}

	wg.Wait()

	// Verify all writes succeeded
	expectedSize := numWriters * writesPerWriter
	if m.Size() != expectedSize {
		t.Errorf("Expected size %d, got %d", expectedSize, m.Size())
	}
}

// Performance benchmark tests
func BenchmarkCopyOnWriteMapGet(b *testing.B) {
	m := NewCopyOnWriteMap[string, int]()

	// Pre-populate data
	for i := 0; i < 1000; i++ {
		m.Put(fmt.Sprintf("key%d", i), i)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("key%d", i%1000)
			m.Get(key)
			i++
		}
	})
}

func BenchmarkCopyOnWriteMapPut(b *testing.B) {
	m := NewCopyOnWriteMap[string, int]()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key%d", i)
		m.Put(key, i)
	}
}

func BenchmarkCopyOnWriteMapMixed(b *testing.B) {
	m := NewCopyOnWriteMap[string, int]()

	// Pre-populate data
	for i := 0; i < 100; i++ {
		m.Put(fmt.Sprintf("key%d", i), i)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if i%10 == 0 { // 10% write operations
				key := fmt.Sprintf("key%d", i%1000)
				m.Put(key, i)
			} else { // 90% read operations
				key := fmt.Sprintf("key%d", i%100)
				m.Get(key)
			}
			i++
		}
	})
}

// Memory usage test
func TestCopyOnWriteMapMemoryUsage(t *testing.T) {
	m := NewCopyOnWriteMap[string, int]()

	// Add large amount of data
	for i := 0; i < 10000; i++ {
		m.Put(fmt.Sprintf("key%d", i), i)
	}

	// Get snapshot
	snapshot1 := m.Snapshot()

	// Modify data
	m.Put("newkey", 99999)

	// Get new snapshot
	snapshot2 := m.Snapshot()

	// Verify snapshot independence
	if len(snapshot1) == len(snapshot2) {
		t.Error("Snapshots should have different lengths")
	}

	if _, exists := snapshot1["newkey"]; exists {
		t.Error("Old snapshot should not contain new key")
	}

	if _, exists := snapshot2["newkey"]; !exists {
		t.Error("New snapshot should contain new key")
	}
}

// Stress test
func TestCopyOnWriteMapStress(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping stress test in short mode")
	}

	m := NewCopyOnWriteMap[int, string]()
	const duration = 2 * time.Second
	const numGoroutines = 20

	var wg sync.WaitGroup
	stop := make(chan struct{})

	// Start readers
	for i := 0; i < numGoroutines/2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter := 0
			for {
				select {
				case <-stop:
					return
				default:
					m.Get(counter % 100)
					counter++
				}
			}
		}()
	}

	// Start writers
	for i := 0; i < numGoroutines/2; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			counter := 0
			for {
				select {
				case <-stop:
					return
				default:
					m.Put(counter, fmt.Sprintf("value%d_%d", id, counter))
					counter++
				}
			}
		}(i)
	}

	// Run for specified duration
	time.Sleep(duration)
	close(stop)
	wg.Wait()

	t.Logf("Stress test completed. Final map size: %d", m.Size())
}
