package data_model

// PortItem represents a JSON stream port entry
type PortItem struct {
	//ID is the item identifier of a port
	ID string
	// Port represents the data about a port, i.e. the main data model item of the application.
	// probably not the best thing @TODO refactor!, it should give lint warnings!
	Port Port
	// Error
	Error error
}
