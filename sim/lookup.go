package sim

type Lookup[K comparable, V KeyProvider[K]] map[K]V

type KeyProvider[K comparable] interface {
	Key() K
}

func (l Lookup[K, V]) HasKey(k K) (ok bool) {
	_, ok = l[k]
	return
}

func (l Lookup[K, V]) Has(v V) (ok bool) {
	_, ok = l[v.Key()]
	return
}

func (l Lookup[K, V]) Add(v V) {
	l[v.Key()] = v
}

func (l Lookup[K, V]) Del(v V) {
	delete(l, v.Key())
}

func (l Lookup[K, V]) Keys() (out []K) {
	out = make([]K, 0, len(l))
	for k := range l {
		out = append(out, k)
	}
	return
}

// Copy returns a copy of the lookup.
func (l Lookup[K, V]) Copy() Lookup[K, V] {
	output := make(map[K]V, len(l))
	for k, v := range l {
		output[k] = v
	}
	return output
}
