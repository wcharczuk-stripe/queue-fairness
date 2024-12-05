package sim

import "math/rand/v2"

func RandomKeyByWeight[K comparable, V Operatable](r *rand.Rand, keyWeights map[K]V) K {
	total := SumValues(keyWeights)
	nf := r.Float64()
	var accum float64
	for key, w := range keyWeights {
		accum += float64(w) / float64(total)
		if nf <= accum {
			return key
		}
	}
	var zero K
	return zero
}

func SumValues[K comparable, V Operatable](m map[K]V) (output V) {
	for _, v := range m {
		output += v
	}
	return
}
