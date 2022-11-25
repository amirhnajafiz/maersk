package maersk

// Chunk
// every part of our file is a chunk.
// we assemble the file based on chunks index.
type Chunk struct {
	// Data stores the file bytes
	Data []byte
	// Index is the chunk index
	Index int
}
