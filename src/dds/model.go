package dds

// Model is the client interface for interacting with a distributed data model
type Model interface {

	// Get gets the object
	Get(id string, i interface{}) error

	// Create creates an object
	Create(i interface{}) (string, error)

	// Update creates an object
	Update(ID string, i interface{}) error

	// Delete deletes an object
	Delete(ID string) error

	// Subscribe subscribes to the event
	Subscribe(evt string, ch chan string) func()

	// Close closes the model connection
	Close()
}
