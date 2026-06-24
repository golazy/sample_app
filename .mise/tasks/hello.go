// &> /dev/null; exec go run "$0" "$@"
//MISE description="Example: run a Go script task with mise"
package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello from a Go mise task")
	fmt.Println("This task is an example of using Go for small project scripts.")
}
