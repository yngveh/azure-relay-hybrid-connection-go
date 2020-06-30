// +build debug

package arhc

import (
	"log"
)

func debug(format string, v ...interface{}) {
	log.Printf("[arhc] "+format, v...)
}
