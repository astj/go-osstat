// +build solaris,cgo

package cpu

import (
	"github.com/siebenmann/go-kstat"
)

// #include <sys/types.h>
// #include <sys/processor.h>
// #include <unistd.h>
import "C"

// Get cpu statistics
func Get() (*Stats, error) {
	return collectCPUStats()
}

// Stats represents cpu statistics for solaris
// mpstat and vmstat report kstat 'kernel' as 'system'
type Stats struct {
	User, System, Stolen, Idle, Wait, Total uint64
	CPUCount                                int
}

func collectCPUStats() (*Stats, error) {
	tok, err := kstat.Open()
	if err != nil {
		return nil, err
	}
	defer tok.Close()

	// ref: illumos: usr/src/cmd/stat/common/acquire.c
	maxCPUs := int(C.sysconf(C._SC_CPUID_MAX)) + 1
	cpu := Stats{}
	for i := 0; i < maxCPUs; i++ {
		state, err := C.p_online(C.processorid_t(i), C.P_STATUS)
		if err != nil || int(state) == -1 {
			// failed to stat
			continue
		}
		if state != C.P_ONLINE && state != C.P_NOINTR {
			// not online/available
			continue
		}

		sys, err := tok.Lookup("cpu", i, "sys")
		if err != nil {
			return nil, err
		}

		if err := getNamedAndAdd(sys, "cpu_ticks_idle", &cpu.Idle); err != nil {
			return nil, err
		}
		if err := getNamedAndAdd(sys, "cpu_ticks_kernel", &cpu.System); err != nil {
			return nil, err
		}
		if err := getNamedAndAdd(sys, "cpu_ticks_stolen", &cpu.Stolen); err != nil {
			return nil, err
		}
		if err := getNamedAndAdd(sys, "cpu_ticks_user", &cpu.User); err != nil {
			return nil, err
		}
		if err := getNamedAndAdd(sys, "cpu_ticks_wait", &cpu.Wait); err != nil {
			return nil, err
		}
		cpu.CPUCount++
	}

	cpu.Total = cpu.Idle + cpu.System + cpu.Stolen + cpu.User + cpu.Wait
	return &cpu, nil
}

func getNamedAndAdd(stat *kstat.KStat, name string, v *uint64) error {
	n, err := stat.GetNamed(name)
	if err != nil {
		return err
	}

	*v += n.UintVal
	return nil
}
