package workers

import (
	"github.com/gocraft/work"
)

func (c *Context) SendEmail(job *work.Job) error {
	// Extract arguments:
	// addr := job.ArgString("address")
	// subject := job.ArgString("subject")
	// if err := job.ArgError(); err != nil {
	//  return err
	// }

	// Go ahead and send the email...
	// sendEmailTo(addr, subject)

	return nil
}
