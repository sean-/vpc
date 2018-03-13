package agent

import (
	"context"

	"github.com/pkg/errors"
	"github.com/sean-/vpc/db"
)

type Agent struct {
	DbPool      *db.Pool
	shutdownCtx context.Context
	shutdown    func()
}

func New() (agent *Agent, err error) {
	a := &Agent{}

	a.shutdownCtx, a.shutdown = context.WithCancel(context.Background())

	return a, nil
}

func (a *Agent) Start() error {
	return nil
}

func (a *Agent) Stop() error {
	if err := a.Shutdown(); err != nil {
		return errors.Wrap(err, "shutdown failed while stopping agent")
	}

	return nil
}

func (a *Agent) Shutdown() error {
	if a.DbPool != nil {
		a.DbPool.Close()
	}

	return nil
}

// 5. Listen on the socket for UDP packets
// 6. Parse packet
// 7. Look up the results in the database
// 8. Respond to packet
func (a *Agent) Run() error {
	return nil
}
