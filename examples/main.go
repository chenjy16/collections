package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println("Go Collections Library Examples")
	fmt.Println("================================")

	// Run all container examples
	RunContainerExamples()

	// Run iterator examples
	fmt.Println("\n" + strings.Repeat("=", 50))
	RunIteratorExamples()

	// Run multiset examples
	fmt.Println("\n" + strings.Repeat("=", 50))
	RunMultisetExamples()

	// Run multimap examples
	fmt.Println("\n" + strings.Repeat("=", 50))
	RunMultimapExamples()
}
