// +build darwin

package memory

import (
	"io"
	"reflect"
	"strings"
	"testing"
)

type memoryGeneratorMock struct {
}

func (generator memoryGeneratorMock) Start() error {
	return nil
}

func (generator memoryGeneratorMock) Output() (io.Reader, error) {
	return strings.NewReader(
		`Mach Virtual Memory Statistics: (page size of 4096 bytes)
Pages free:                               72827.
Pages active:                           2154445.
Pages inactive:                         1511468.
Pages speculative:                         8107.
Pages throttled:                              0.
Pages wired down:                        446975.
Pages purgeable:                         383371.
"Translation faults":                  97589077.
Pages copy-on-write:                    3305869.
Pages zero filled:                     50848672.
Pages reactivated:                         1999.
Pages purged:                           2496610.
File-backed pages:                       677870.
Anonymous pages:                        2996150.
Pages stored in compressor:                   0.
Pages occupied by compressor:                 0.
Decompressions:                               0.
Compressions:                                 0.
Pageins:                                6333901.
Pageouts:                                   353.
Swapins:                                      0.
Swapouts:                                     0.
`), nil
}

func (generator memoryGeneratorMock) Finish() error {
	return nil
}

func Test_GetMemory(t *testing.T) {
	_, err := Get()
	if err != nil {
		t.Errorf("error should be nil but got: %v", err)
	}
}

func Test_collectMemoryStats(t *testing.T) {
	got, err := collectMemoryStats(memoryGeneratorMock{})
	if err != nil {
		t.Errorf("error should be nil but got: %v", err)
	}
	pageSize := 4096
	expected := Memory{
		Total:    uint64((446975 + 2154445 + 1511468 + 8107 + 72827) * pageSize),
		Used:     uint64((446975 + 2154445 + 1511468 + 8107 - (383371 + 677870)) * pageSize),
		Cached:   uint64((383371 + 677870) * pageSize),
		Free:     uint64(72827 * pageSize),
		Active:   uint64(2154445 * pageSize),
		Inactive: uint64(1511468 * pageSize),
	}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("invalid memory value: %v (expected: %v)", got, expected)
	}
}

type swapGeneratorMock struct {
}

func (generator swapGeneratorMock) Start() error {
	return nil
}

func (generator swapGeneratorMock) Output() (io.Reader, error) {
	return strings.NewReader(
		`total = 4096.00M  used = 3184.75M  free = 911.25M  (encrypted)
`), nil
}

func (generator swapGeneratorMock) Finish() error {
	return nil
}

func Test_collectSwapStats(t *testing.T) {
	got, err := collectSwapStats(swapGeneratorMock{})
	if err != nil {
		t.Errorf("error should be nil but got: %v", err)
	}
	expected := memorySwap{
		total: uint64(4096.00 * 1024 * 1024),
		used:  uint64(3184.75 * 1024 * 1024),
		free:  uint64(911.25 * 1024 * 1024),
	}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("invalid memory swap value: %v (expected: %v)", got, expected)
	}
}