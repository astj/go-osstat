// +build solaris

package uptime

import (
	"time"

	"github.com/siebenmann/go-kstat"
)

func get() (time.Duration, error) {
	tok, err := kstat.Open()
	if err != nil {
		return time.Duration(0), err
	}
	defer tok.Close()

	n, err := tok.GetNamed("unix", -1, "system_misc", "boot_time")
	if err != nil {
		return time.Duration(0), err
	}

	bootTime := time.Unix(int64(n.UintVal), 0)
	return time.Now().Sub(bootTime), nil
}
