package multiset

import (
    "sync"
    "testing"
)

// Race-focused stress test mixing add/remove with concurrent reads
func TestConcurrentHashMultisetRaceMixed(t *testing.T) {
    ms := NewConcurrentHashMultiset[int]()

    keys := 16
    perKeyAdds := 2000

    // Phase 1: concurrent adds
    var wgAdds sync.WaitGroup
    wgAdds.Add(keys)
    for k := 0; k < keys; k++ {
        key := k
        go func() {
            defer wgAdds.Done()
            for i := 0; i < perKeyAdds; i++ {
                ms.Add(key)
            }
        }()
    }
    wgAdds.Wait()

    // Phase 2: concurrent removes interleaved with reads
    var wgRem sync.WaitGroup
    wgRem.Add(keys + 4) // keys removers + 4 readers

    // Removers
    for k := 0; k < keys; k++ {
        key := k
        go func() {
            defer wgRem.Done()
            for i := 0; i < perKeyAdds; i++ {
                ms.Remove(key)
            }
        }()
    }

    // Readers performing various queries concurrently
    for r := 0; r < 4; r++ {
        go func() {
            defer wgRem.Done()
            // perform a bunch of read operations while removers run
            for i := 0; i < perKeyAdds; i++ {
                _ = ms.TotalSize()
                _ = ms.DistinctElements()
                for _, k := range ms.ElementSet() {
                    _ = ms.Count(k)
                }
                _ = ms.EntrySet()
            }
        }()
    }

    wgRem.Wait()

    // After balanced add/remove, all counts should be zero
    if ms.TotalSize() != 0 {
        t.Fatalf("expected TotalSize() == 0 after balanced operations, got %d", ms.TotalSize())
    }
    for k := 0; k < keys; k++ {
        if c := ms.Count(k); c != 0 {
            t.Fatalf("expected Count(%d) == 0, got %d", k, c)
        }
    }
}