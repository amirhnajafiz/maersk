package maersk

// Every part of our file is a chunk.
// We assemble the file based on chunks index.
type chunk struct {
	// data stores the file bytes.
	data []byte
	// index is the chunk index in the file.
	index int
}

// To download each chunk, we send a job into jobs
// channel of cargo.
type job struct {
	// job index
	index int
	// job size
	size int
	// is the last job or not
	last bool
}
