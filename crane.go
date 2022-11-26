package maersk

import "fmt"

// crane
// manages to keep track of jobs and chunks.
// it resends the failed jobs into jobs channel again.
type crane struct {
	// jobs channel.
	jobs chan job
	// failed jobs channel.
	failed chan job
	// crane program for managing the failed jobs.
	program func(int) error
}

// start
// crane to work.
func (c *crane) start() error {
	for j := range c.failed {
		c.jobs <- j

		if err := c.program(j.index); err != nil {
			return fmt.Errorf("crane stopped: %v", err)
		}
	}

	return nil
}
