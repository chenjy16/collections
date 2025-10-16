package multiset

// baseMultisetIterator provides shared iterator logic for multiset implementations.
// It iterates elements according to counts in EntrySet and supports Remove via callback.
type baseMultisetIterator[E comparable] struct {
    entries    []Entry[E]
    index      int
    current    int
    removeFunc func(E) bool
    refresh    func() []Entry[E]
}

func (it *baseMultisetIterator[E]) HasNext() bool {
    return it.index < len(it.entries) && (it.current < it.entries[it.index].Count || it.index+1 < len(it.entries))
}

func (it *baseMultisetIterator[E]) Next() (E, bool) {
    if !it.HasNext() {
        var zero E
        return zero, false
    }
    if it.current >= it.entries[it.index].Count {
        it.index++
        it.current = 0
    }
    element := it.entries[it.index].Element
    it.current++
    return element, true
}

// Reset is optional and not part of the public Iterator interface; kept for internal parity.
func (it *baseMultisetIterator[E]) Reset() {
    if it.refresh != nil {
        it.entries = it.refresh()
    }
    it.index = 0
    it.current = 0
}

func (it *baseMultisetIterator[E]) Remove() bool {
    if it.index >= len(it.entries) || it.current == 0 {
        return false
    }
    if it.removeFunc == nil {
        return false
    }
    element := it.entries[it.index].Element
    _ = it.removeFunc(element)
    // Refresh entries after removal to keep iterator consistent
    if it.refresh != nil {
        it.entries = it.refresh()
    }
    if it.index >= len(it.entries) {
        it.index = len(it.entries)
        it.current = 0
    } else if it.current > it.entries[it.index].Count {
        it.current = it.entries[it.index].Count
    }
    return true
}