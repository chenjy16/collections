package maps

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestNewConcurrentHashMap(t *testing.T) {
	chm := NewConcurrentHashMap[string, int]()
	if chm == nil {
		t.Fatal("NewConcurrentHashMap returned nil")
	}
	if !chm.IsEmpty() {
		t.Error("New map should be empty")
	}
	if chm.Size() != 0 {
		t.Error("New map size should be 0")
	}
}

func TestNewConcurrentHashMapWithCapacity(t *testing.T) {
	chm := NewConcurrentHashMapWithCapacity[string, int](100)
	if chm == nil {
		t.Fatal("NewConcurrentHashMapWithCapacity returned nil")
	}
	if !chm.IsEmpty() {
		t.Error("New map should be empty")
	}
}

func TestConcurrentHashMapFromMap(t *testing.T) {
	source := map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
	}
	chm := ConcurrentHashMapFromMap(source)
	if chm.Size() != 3 {
		t.Errorf("Expected size 3, got %d", chm.Size())
	}
	if val, ok := chm.Get("two"); !ok || val != 2 {
		t.Error("Failed to get value from copied map")
	}
}

func TestConcurrentHashMapPutAndGet(t *testing.T) {
	chm := NewConcurrentHashMap[string, int]()

	// Test Put
	chm.Put("key1", 100)

	// Test Get
	val, ok := chm.Get("key1")
	if !ok {
		t.Error("Key should exist after Put")
	}
	if val != 100 {
		t.Errorf("Expected value 100, got %d", val)
	}

	// Test Put existing key
	chm.Put("key1", 200)

	// Verify new value
	val, ok = chm.Get("key1")
	if !ok || val != 200 {
		t.Error("Failed to update existing key")
	}
}

func TestConcurrentHashMapRemove(t *testing.T) {
	chm := NewConcurrentHashMap[string, int]()
	chm.Put("key1", 100)
	chm.Put("key2", 200)

	// Remove existing key
	oldVal, existed := chm.Remove("key1")
	if !existed {
		t.Error("Key should exist")
	}
	if oldVal != 100 {
		t.Errorf("Expected old value 100, got %d", oldVal)
	}

	// Verify removal
	_, ok := chm.Get("key1")
	if ok {
		t.Error("Key should not exist after removal")
	}

	// Remove non-existing key
	_, existed = chm.Remove("nonexistent")
	if existed {
		t.Error("Non-existing key should not be found")
	}
}

func TestConcurrentHashMapContains(t *testing.T) {
	chm := NewConcurrentHashMap[string, int]()
	chm.Put("key1", 100)
	chm.Put("key2", 200)

	// Test ContainsKey
	if !chm.ContainsKey("key1") {
		t.Error("Map should contain key1")
	}
	if chm.ContainsKey("nonexistent") {
		t.Error("Map should not contain nonexistent key")
	}

	// Test ContainsValue
	if !chm.ContainsValue(100) {
		t.Error("Map should contain value 100")
	}
	if chm.ContainsValue(999) {
		t.Error("Map should not contain value 999")
	}
}

func TestConcurrentHashMapSizeAndEmpty(t *testing.T) {
	chm := NewConcurrentHashMap[string, int]()

	if !chm.IsEmpty() {
		t.Error("New map should be empty")
	}
	if chm.Size() != 0 {
		t.Error("New map size should be 0")
	}

	chm.Put("key1", 100)
	if chm.IsEmpty() {
		t.Error("Map should not be empty after adding element")
	}
	if chm.Size() != 1 {
		t.Errorf("Expected size 1, got %d", chm.Size())
	}

	chm.Put("key2", 200)
	if chm.Size() != 2 {
		t.Errorf("Expected size 2, got %d", chm.Size())
	}
}

func TestConcurrentHashMapClear(t *testing.T) {
	chm := NewConcurrentHashMap[string, int]()
	chm.Put("key1", 100)
	chm.Put("key2", 200)

	chm.Clear()
	if !chm.IsEmpty() {
		t.Error("Map should be empty after Clear")
	}
	if chm.Size() != 0 {
		t.Error("Map size should be 0 after Clear")
	}
}

func TestConcurrentHashMapKeysAndValues(t *testing.T) {
	chm := NewConcurrentHashMap[string, int]()
	chm.Put("key1", 100)
	chm.Put("key2", 200)
	chm.Put("key3", 300)

	keys := chm.Keys()
	if len(keys) != 3 {
		t.Errorf("Expected 3 keys, got %d", len(keys))
	}

	values := chm.Values()
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

func TestConcurrentHashMapForEach(t *testing.T) {
	chm := NewConcurrentHashMap[string, int]()
	chm.Put("key1", 100)
	chm.Put("key2", 200)
	chm.Put("key3", 300)

	sum := 0
	count := 0
	chm.ForEach(func(k string, v int) {
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

func TestConcurrentHashMapString(t *testing.T) {
	chm := NewConcurrentHashMap[string, int]()

	// Empty map
	str := chm.String()
	if str != "{}" {
		t.Errorf("Expected '{}', got '%s'", str)
	}

	// Non-empty map
	chm.Put("key1", 100)
	str = chm.String()
	if str != "{key1=100}" {
		t.Errorf("Expected '{key1=100}', got '%s'", str)
	}
}

func TestConcurrentHashMapPutIfAbsent(t *testing.T) {
	chm := NewConcurrentHashMap[string, int]()

	// Put new key
	oldVal, inserted := chm.PutIfAbsent("key1", 100)
	if !inserted {
		t.Error("Should have inserted new key")
	}
	if oldVal != 0 {
		t.Error("Old value should be zero for new key")
	}

	// Try to put existing key
	oldVal, inserted = chm.PutIfAbsent("key1", 200)
	if inserted {
		t.Error("Should not have inserted existing key")
	}
	if oldVal != 100 {
		t.Errorf("Expected old value 100, got %d", oldVal)
	}

	// Verify value unchanged
	if val, ok := chm.Get("key1"); !ok || val != 100 {
		t.Error("Value should remain unchanged")
	}
}

func TestConcurrentHashMapReplace(t *testing.T) {
	chm := NewConcurrentHashMap[string, int]()
	chm.Put("key1", 100)

	// Replace existing key
	oldVal, replaced := chm.Replace("key1", 200)
	if !replaced {
		t.Error("Should have replaced existing key")
	}
	if oldVal != 100 {
		t.Errorf("Expected old value 100, got %d", oldVal)
	}

	// Try to replace non-existing key
	_, replaced = chm.Replace("nonexistent", 300)
	if replaced {
		t.Error("Should not have replaced non-existing key")
	}
}

func TestConcurrentHashMapReplaceIf(t *testing.T) {
	chm := NewConcurrentHashMap[string, int]()
	chm.Put("key1", 100)

	// Replace with correct old value
	replaced := chm.ReplaceIf("key1", 100, 200)
	if !replaced {
		t.Error("Should have replaced with correct old value")
	}

	// Try to replace with incorrect old value
	replaced = chm.ReplaceIf("key1", 100, 300)
	if replaced {
		t.Error("Should not have replaced with incorrect old value")
	}

	// Verify current value
	if val, ok := chm.Get("key1"); !ok || val != 200 {
		t.Error("Value should be 200")
	}
}

func TestConcurrentHashMapComputeIfAbsent(t *testing.T) {
	chm := NewConcurrentHashMap[string, int]()

	// Compute for non-existing key
	value := chm.ComputeIfAbsent("key1", func(k string) int {
		return len(k) * 10
	})
	if value != 40 { // len("key1") * 10 = 4 * 10 = 40
		t.Errorf("Expected computed value 40, got %d", value)
	}

	// Verify the value was stored
	if val, ok := chm.Get("key1"); !ok || val != 40 {
		t.Error("Computed value should be stored")
	}

	// Compute for existing key (should return existing value)
	value = chm.ComputeIfAbsent("key1", func(k string) int {
		return 999 // This should not be called
	})
	if value != 40 {
		t.Errorf("Expected existing value 40, got %d", value)
	}
}

func TestConcurrentHashMapComputeIfPresent(t *testing.T) {
	chm := NewConcurrentHashMap[string, int]()
	chm.Put("key1", 100)

	// Compute for existing key
	newValue, computed := chm.ComputeIfPresent("key1", func(k string, v int) int {
		return v * 2
	})
	if !computed {
		t.Error("Should have computed for existing key")
	}
	if newValue != 200 {
		t.Errorf("Expected computed value 200, got %d", newValue)
	}

	// Compute for non-existing key
	_, computed = chm.ComputeIfPresent("nonexistent", func(k string, v int) int {
		return v * 2
	})
	if computed {
		t.Error("Should not have computed for non-existing key")
	}
}

func TestConcurrentHashMapPutAll(t *testing.T) {
	chm1 := NewConcurrentHashMap[string, int]()
	chm1.Put("key1", 100)

	chm2 := NewConcurrentHashMap[string, int]()
	chm2.Put("key2", 200)
	chm2.Put("key3", 300)

	chm1.PutAll(chm2)

	if chm1.Size() != 3 {
		t.Errorf("Expected size 3, got %d", chm1.Size())
	}

	if val, ok := chm1.Get("key2"); !ok || val != 200 {
		t.Error("Failed to copy key2 from chm2")
	}
}

func TestConcurrentHashMapPutAllFromMap(t *testing.T) {
	chm := NewConcurrentHashMap[string, int]()
	chm.Put("key1", 100)

	source := map[string]int{
		"key2": 200,
		"key3": 300,
	}

	chm.PutAllFromMap(source)

	if chm.Size() != 3 {
		t.Errorf("Expected size 3, got %d", chm.Size())
	}

	if val, ok := chm.Get("key2"); !ok || val != 200 {
		t.Error("Failed to copy key2 from source map")
	}
}

func TestConcurrentHashMapToMap(t *testing.T) {
	chm := NewConcurrentHashMap[string, int]()
	chm.Put("key1", 100)
	chm.Put("key2", 200)

	goMap := chm.ToMap()
	if len(goMap) != 2 {
		t.Errorf("Expected map length 2, got %d", len(goMap))
	}

	if goMap["key1"] != 100 {
		t.Error("Failed to convert to Go map")
	}
}

func TestConcurrentHashMapSnapshot(t *testing.T) {
	chm := NewConcurrentHashMap[string, int]()
	chm.Put("key1", 100)
	chm.Put("key2", 200)

	snapshot := chm.Snapshot()
	if len(snapshot) != 2 {
		t.Errorf("Expected snapshot length 2, got %d", len(snapshot))
	}

	// Modify original map
	chm.Put("key3", 300)

	// Snapshot should remain unchanged
	if len(snapshot) != 2 {
		t.Error("Snapshot should not be affected by subsequent modifications")
	}
}

// Concurrent test
func TestConcurrentHashMapConcurrency(t *testing.T) {
	chm := NewConcurrentHashMap[string, int]()

	// Pre-populate data
	for i := 0; i < 100; i++ {
		chm.Put(fmt.Sprintf("key%d", i), i)
	}

	const numGoroutines = 20
	const operationsPerGoroutine = 100

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Start multiple concurrent readers
	for i := 0; i < numGoroutines/2; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < operationsPerGoroutine; j++ {
				key := fmt.Sprintf("key%d", j%100)
				if val, ok := chm.Get(key); !ok || val != j%100 {
					t.Errorf("Goroutine %d: unexpected value for %s", id, key)
					return
				}
			}
		}(i)
	}

	// Start writers and readers
	for i := numGoroutines / 2; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < operationsPerGoroutine; j++ {
				if j%2 == 0 {
					key := fmt.Sprintf("new_key_%d_%d", id, j)
					chm.Put(key, id*1000+j)
				} else {
					key := fmt.Sprintf("key%d", j%100)
					chm.Get(key) // Don't check result, just ensure no panic
				}
			}
		}(i)
	}

	wg.Wait()
}

func TestConcurrentHashMapResize(t *testing.T) {
	chm := NewConcurrentHashMapWithCapacity[int, string](4) // Small capacity to trigger resize

	// Add enough elements to trigger resize
	for i := 0; i < 100; i++ {
		chm.Put(i, fmt.Sprintf("value%d", i))
	}

	// Verify all elements exist
	for i := 0; i < 100; i++ {
		if val, ok := chm.Get(i); !ok || val != fmt.Sprintf("value%d", i) {
			t.Errorf("Failed to get value for key %d", i)
		}
	}

	if chm.Size() != 100 {
		t.Errorf("Expected size 100, got %d", chm.Size())
	}
}

// Performance benchmark tests
func BenchmarkConcurrentHashMapGet(b *testing.B) {
	chm := NewConcurrentHashMap[string, int]()

	// Pre-populate data
	for i := 0; i < 1000; i++ {
		chm.Put(fmt.Sprintf("key%d", i), i)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("key%d", i%1000)
			chm.Get(key)
			i++
		}
	})
}

func BenchmarkConcurrentHashMapPut(b *testing.B) {
	chm := NewConcurrentHashMap[string, int]()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key%d", i)
		chm.Put(key, i)
	}
}

func BenchmarkConcurrentHashMapMixed(b *testing.B) {
	chm := NewConcurrentHashMap[string, int]()

	// Pre-populate data
	for i := 0; i < 100; i++ {
		chm.Put(fmt.Sprintf("key%d", i), i)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if i%10 == 0 { // 10% write operations
				key := fmt.Sprintf("key%d", i%1000)
				chm.Put(key, i)
			} else { // 90% read operations
				key := fmt.Sprintf("key%d", i%100)
				chm.Get(key)
			}
			i++
		}
	})
}

// Stress test
func TestConcurrentHashMapStress(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping stress test in short mode")
	}

	chm := NewConcurrentHashMap[int, string]()
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
					chm.Get(counter % 100)
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
					chm.Put(counter, fmt.Sprintf("value%d_%d", id, counter))
					counter++
				}
			}
		}(i)
	}

	// Run for specified duration
	time.Sleep(duration)
	close(stop)
	wg.Wait()

	t.Logf("Stress test completed. Final map size: %d", chm.Size())
}
