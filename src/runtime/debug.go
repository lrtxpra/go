// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	"runtime/internal/atomic"
	"unsafe"
)

// GOMAXPROCS sets the maximum number of CPUs that can be executing
// simultaneously and returns the previous setting.  If n < 1, it does not
// change the current setting.
// The number of logical CPUs on the local machine can be queried with NumCPU.
// This call will go away when the scheduler improves.
func GOMAXPROCS(n int) int {
	if n > _MaxGomaxprocs {
		n = _MaxGomaxprocs
	}
	lock(&sched.lock)
	ret := int(gomaxprocs)
	unlock(&sched.lock)
	if n <= 0 || n == ret {
		return ret
	}

	stopTheWorld("GOMAXPROCS")

	// newprocs will be processed by startTheWorld
	newprocs = int32(n)

	startTheWorld()
	return ret
}

// NumCPU returns the number of logical CPUs usable by the current process.
func NumCPU() int {
	return int(ncpu)
}

// NumCgoCall returns the number of cgo calls made by the current process.
func NumCgoCall() int64 {
	var n int64
	for mp := (*m)(atomic.Loadp(unsafe.Pointer(&allm))); mp != nil; mp = mp.alllink {
		n += int64(mp.ncgocall)
	}
	return n
}

// NumGoroutine returns the number of goroutines that currently exist.
func NumGoroutine() int {
	return int(gcount())
}
