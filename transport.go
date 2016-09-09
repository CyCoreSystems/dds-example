package dds

// A Transport does the heavy lifting of connecting the various components together.
type Transport interface {
	Model(mf *ModelFactory, st Storage) error
	Action(af *ActionFactory, h Handler) error

	Close()
}
