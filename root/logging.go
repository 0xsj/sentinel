package root

import (
	"fmt"
	"os"

	"github.com/0xsj/atelier-wails/pkg/clock"
	"github.com/0xsj/atelier-wails/pkg/id"
	"github.com/0xsj/atelier-wails/pkg/logger"
	"github.com/0xsj/atelier-wails/pkg/provenance"
)

func consoleLogger() (*logger.Logger, error) {
	sink, err := logger.NewConsole(os.Stderr)
	if err != nil {
		return nil, err
	}
	clock := clock.System{}
	ids, err := id.NewV7(clock)
	if err != nil {
		return nil, err
	}
	return bootstrapLogger(clock, ids, sink, "atelier-wails")
}

// bootstrapLogger opens one real process execution. The resulting scope is
// bound to every lifecycle record; callers do not need to manufacture IDs or
// provenance for startup diagnostics.
func bootstrapLogger(c provenance.Clock, ids provenance.Generator, sink logger.Sink, service string) (*logger.Logger, error) {
	log, err := logger.New(c, sink, service, logger.Info)
	if err != nil {
		return nil, err
	}
	operation, err := provenance.NewOperation("app.startup")
	if err != nil {
		return nil, err
	}
	executor, err := provenance.NewActor(provenance.System, service)
	if err != nil {
		return nil, err
	}
	factory, err := provenance.NewFactory(c, ids)
	if err != nil {
		return nil, err
	}
	scope, err := factory.Open(provenance.RootSpec{
		Origin:    provenance.Startup,
		Operation: operation,
		Executor:  executor,
	})
	if err != nil {
		return nil, err
	}
	return log.WithScope(scope)
}

// Diagnostics never replace the native host's business result. Do not recursively
// log sink failures or print their potentially sensitive raw diagnostic causes.
func reportLog(err error) {
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "atelier: diagnostic output unavailable")
	}
}
func runLogged(log *logger.Logger, run func() error) error {
	reportLog(log.Log(logger.Info, "app.starting", nil))
	err := run()
	if err != nil {
		reportLog(log.WithError(err).Log(logger.Error, "app.failed", nil))
	} else {
		reportLog(log.Log(logger.Info, "app.stopped", nil))
	}
	reportLog(log.Flush())
	return err
}
