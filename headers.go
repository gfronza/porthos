package porthos

// Headers represents RPC headers (request and response).
type Headers struct {
	headers map[string]interface{}
}

// NewHeaders creates a new Headers object initializing the map.
func NewHeaders() *Headers {
	return &Headers{
		make(map[string]interface{}),
	}
}

// NewHeadersFromMap creates a new Headers from a given map.
func NewHeadersFromMap(m map[string]interface{}) *Headers {
	return &Headers{
		m,
	}
}

// Set a header.
func (h *Headers) Set(key string, value string) {
	h.headers[key] = value
}

// Get a header.
func (h *Headers) Get(key string) string {
	return h.headers[key].(string)
}

// Delete a header.
func (h *Headers) Delete(key string) {
	delete(h.headers, key)
}

func (h *Headers) asMap() map[string]interface{} {
	return h.headers
}
