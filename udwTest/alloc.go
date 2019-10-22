package udwTest

import "runtime"

func AllocSizeDiff(fn func()) uint64 {
	var memstats1 runtime.MemStats
	var memstats2 runtime.MemStats
	runtime.ReadMemStats(&memstats1)
	fn()
	runtime.ReadMemStats(&memstats2)
	return memstats2.TotalAlloc - memstats1.TotalAlloc
}
