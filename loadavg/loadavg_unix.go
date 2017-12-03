// +build !windows,cgo

package loadavg

import (
	"errors"
)

// #if defined(sun) || defined(__sun)
// #include <sys/time.h>
// #include <sys/loadavg.h>
// #else
// #include <stdlib.h>
// #endif
import "C"

// Reference: man 3 getloadavg
func get() (*Stats, error) {
	var loadavgs [3]C.double
	ret := C.getloadavg(&loadavgs[0], 3)
	if ret != 3 {
		return nil, errors.New("failed to get loadavg")
	}
	return &Stats{
		Loadavg1:  float64(loadavgs[0]),
		Loadavg5:  float64(loadavgs[1]),
		Loadavg15: float64(loadavgs[2]),
	}, nil
}
