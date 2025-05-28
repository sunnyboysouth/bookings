package forms

type errors map[string][]string

// Add adds an error message for the given field. If the field already
// has an error message, it appends the new message to the existing
// message.
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

// Get returns the first error message for the given field. if no error
// message exists for the field, it returns an empty string.
func (e errors) Get(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	}
	return es[0]
}
