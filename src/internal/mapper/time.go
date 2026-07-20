package mapper

import "time"

func UnixSecTimeToTime(xs []float64) []time.Time {
	xsTimestamps := make([]time.Time, len(xs))
	for i, x := range xs {
		// Assuming xs are a unix timestamp in seconds, convert to time.Time
		xsTimestamps[i] = time.Unix(int64(x), 0)
	}
	return xsTimestamps
}

func TimeToUnixSecTime(times []time.Time) []float64 {
	xs := make([]float64, len(times))
	for i, t := range times {
		// convert time.Time to unix timestamp in seconds
		xs[i] = float64(t.Unix())
	}
	return xs
}
