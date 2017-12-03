// +build !windows,!linux,!darwin,!solaris

package main

var generators []generator

func init() {
	generators = []generator{
		&loadavgGenerator{},
		&uptimeGenerator{},
		&memoryGenerator{},
		&networkGenerator{},
	}
}
