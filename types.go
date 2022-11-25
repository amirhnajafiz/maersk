package maersk

// Every part of our file is a chunk.
// We assemble the file based on chunks index.
type chunk struct {
	// data stores the file bytes.
	data []byte
	// index is the chunk index in the file.
	index int
}
