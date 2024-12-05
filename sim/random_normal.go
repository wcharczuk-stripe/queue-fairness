package sim

import "math/rand/v2"

// RandomNormal returns a random value with a normal (binomial) distribution.
func RandomNormal(r *rand.Rand, desiredMean, desiredStdDev float64) float64 {
	nf := r.NormFloat64()
	return nf*float64(desiredStdDev) + float64(desiredMean)
}
