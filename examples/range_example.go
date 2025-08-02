package main

import (
	"fmt"
	"github.com/chenjianyu/collections/container/range"
)

func rangeExample() {
	fmt.Println("=== Range Collections Example ===")

	// 1. Basic Range Operations
	fmt.Println("1. Basic Range Operations:")
	r1 := ranges.ClosedRange(1, 10)     // [1, 10]
	r2 := ranges.OpenRange(5, 15)       // (5, 15)
	r3 := ranges.Singleton(7)           // [7, 7]

	fmt.Printf("r1 = %s\n", r1.String())
	fmt.Printf("r2 = %s\n", r2.String())
	fmt.Printf("r3 = %s\n", r3.String())
	fmt.Printf("r1 contains 5: %t\n", r1.Contains(5))
	fmt.Printf("r2 contains 5: %t\n", r2.Contains(5))
	fmt.Printf("r1 and r2 are connected: %t\n", r1.IsConnected(r2))
	
	intersection := r1.Intersection(r2)
	fmt.Printf("Intersection of r1 and r2: %s\n", intersection.String())
	fmt.Println()

	// 2. TreeRangeSet Operations
	fmt.Println("2. TreeRangeSet Operations:")
	rs := ranges.NewTreeRangeSet[int]()
	rs.Add(ranges.ClosedRange(1, 5))
	rs.Add(ranges.ClosedRange(8, 12))
	rs.Add(ranges.ClosedRange(15, 20))

	fmt.Printf("RangeSet: %s\n", rs.String())
	fmt.Printf("Size: %d\n", rs.Size())
	fmt.Printf("Contains value 3: %t\n", rs.ContainsValue(3))
	fmt.Printf("Contains value 7: %t\n", rs.ContainsValue(7))
	fmt.Printf("Contains range [2,4]: %t\n", rs.ContainsRange(ranges.ClosedRange(2, 4)))

	// Add overlapping range
	rs.Add(ranges.ClosedRange(4, 9))
	fmt.Printf("After adding [4,9]: %s\n", rs.String())
	fmt.Println()

	// 3. ImmutableRangeSet Operations
	fmt.Println("3. ImmutableRangeSet Operations:")
	irs := ranges.NewImmutableRangeSet[int]().(*ranges.ImmutableRangeSet[int])
	irs2 := irs.WithAdd(ranges.ClosedRange(1, 5)).(*ranges.ImmutableRangeSet[int])
	irs3 := irs2.WithAdd(ranges.ClosedRange(8, 12)).(*ranges.ImmutableRangeSet[int])

	fmt.Printf("Original: %s\n", irs.String())
	fmt.Printf("After adding [1,5]: %s\n", irs2.String())
	fmt.Printf("After adding [8,12]: %s\n", irs3.String())
	fmt.Println()

	// 4. TreeRangeMap Operations
	fmt.Println("4. TreeRangeMap Operations:")
	rm := ranges.NewTreeRangeMap[int, string]()
	rm.Put(ranges.ClosedRange(1, 5), "first")
	rm.Put(ranges.ClosedRange(8, 12), "second")
	rm.Put(ranges.ClosedRange(15, 20), "third")

	fmt.Printf("RangeMap: %s\n", rm.String())
	fmt.Printf("Size: %d\n", rm.Size())

	if value, ok := rm.Get(3); ok {
		fmt.Printf("Value for key 3: %s\n", value)
	}
	if value, ok := rm.Get(10); ok {
		fmt.Printf("Value for key 10: %s\n", value)
	}
	if _, ok := rm.Get(7); !ok {
		fmt.Printf("No value for key 7\n")
	}

	// Get entry information
	if rangeKey, value, ok := rm.GetEntry(10); ok {
		fmt.Printf("Entry for key 10: range=%s, value=%s\n", rangeKey.String(), value)
	}
	fmt.Println()

	// 5. ImmutableRangeMap Operations
	fmt.Println("5. ImmutableRangeMap Operations:")
	irm := ranges.NewImmutableRangeMap[int, string]().(*ranges.ImmutableRangeMap[int, string])
	irm2 := irm.WithPut(ranges.ClosedRange(1, 5), "alpha").(*ranges.ImmutableRangeMap[int, string])
	irm3 := irm2.WithPut(ranges.ClosedRange(8, 12), "beta").(*ranges.ImmutableRangeMap[int, string])

	fmt.Printf("Original: %s\n", irm.String())
	fmt.Printf("After adding [1,5]->alpha: %s\n", irm2.String())
	fmt.Printf("After adding [8,12]->beta: %s\n", irm3.String())
	fmt.Println()

	// 6. Range Set Operations
	fmt.Println("6. Range Set Operations:")
	rs1 := ranges.NewTreeRangeSet[int]()
	rs1.Add(ranges.ClosedRange(1, 10))
	rs1.Add(ranges.ClosedRange(15, 20))

	rs2 := ranges.NewTreeRangeSet[int]()
	rs2.Add(ranges.ClosedRange(5, 12))
	rs2.Add(ranges.ClosedRange(18, 25))

	fmt.Printf("Set 1: %s\n", rs1.String())
	fmt.Printf("Set 2: %s\n", rs2.String())

	union := rs1.Union(rs2)
	fmt.Printf("Union: %s\n", union.String())

	intersection2 := rs1.Intersection(rs2)
	fmt.Printf("Intersection: %s\n", intersection2.String())

	difference := rs1.Difference(rs2)
	fmt.Printf("Difference (Set1 - Set2): %s\n", difference.String())
	fmt.Println()

	// 7. Practical Example: Time Intervals
	fmt.Println("7. Practical Example: Meeting Room Booking")
	bookings := ranges.NewTreeRangeMap[int, string]()
	
	// Book meeting rooms (time in hours, 24-hour format)
	bookings.Put(ranges.ClosedRange(9, 11), "Team Standup")
	bookings.Put(ranges.ClosedRange(13, 15), "Client Meeting")
	bookings.Put(ranges.ClosedRange(16, 18), "Code Review")

	fmt.Printf("Meeting Schedule: %s\n", bookings.String())

	// Check availability
	checkTimes := []int{8, 10, 12, 14, 17, 19}
	for _, time := range checkTimes {
		if meeting, ok := bookings.Get(time); ok {
			fmt.Printf("Time %d:00 - Booked: %s\n", time, meeting)
		} else {
			fmt.Printf("Time %d:00 - Available\n", time)
		}
	}

	fmt.Println("\n=== Range Collections Example Complete ===")
}