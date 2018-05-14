package logs

import (
	"fmt"
	"log"
)

// Printfln prints a line with a formatted string.
func Printfln(format string, a ...interface{}) {
	log.Println(fmt.Sprintf(format, a))
}
