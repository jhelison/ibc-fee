package abci

// State is a mock state with only a height and some data
type State struct {
	Height int64
	Data   map[string][]byte
}
