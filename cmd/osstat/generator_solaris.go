// +build cgo

package main

var generators []generator

func init() {
	generators = []generator{
		&cpuGenerator{},
	}
}
