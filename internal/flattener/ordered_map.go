package flattener

// OrderedMap maintains insertion order for key-value pairs
type OrderedMap struct {
	keys   []string
	values map[string]string
}

// NewOrderedMap creates a new ordered map
func NewOrderedMap() *OrderedMap {
	return &OrderedMap{
		keys:   make([]string, 0),
		values: make(map[string]string),
	}
}

// Set adds or updates a key-value pair
func (om *OrderedMap) Set(key, value string) {
	if _, exists := om.values[key]; !exists {
		om.keys = append(om.keys, key)
	}
	om.values[key] = value
}

// Get retrieves a value by key
func (om *OrderedMap) Get(key string) (string, bool) {
	value, exists := om.values[key]
	return value, exists
}

// ToMap converts the ordered map to a regular map (loses order)
func (om *OrderedMap) ToMap() map[string]string {
	result := make(map[string]string, len(om.values))
	for _, key := range om.keys {
		result[key] = om.values[key]
	}
	return result
}

// Keys returns all keys in insertion order
func (om *OrderedMap) Keys() []string {
	return om.keys
}

// Len returns the number of key-value pairs
func (om *OrderedMap) Len() int {
	return len(om.keys)
}
