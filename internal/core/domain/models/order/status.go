package order

type Status int

const (
	Created = iota
	Assigned
	Completed
)

func (s Status) String() string {
	switch s {
	case Created:
		return "Created"
	case Assigned:
		return "Assigned"
	case Completed:
		return "Completed"
	default:
		return "Unknown"
	}
}
