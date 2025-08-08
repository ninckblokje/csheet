package csheet

type Entries struct {
	Entries []Entry `json:"entries"`
}

func (e *Entries) String() string {
	var result string

	for _, entry := range e.Entries {
		result += entry.String() + "\n"
	}

	return result
}

type Entry struct {
	Subject string `json:"subject"`
	Section string `json:"section"`
}

func (entry *Entry) String() string {
	return entry.Subject + " " + entry.Section
}
