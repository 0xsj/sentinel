// Package root composes the native process and owns its lifetime.
package root

import (
	"embed"
	wailshost "github.com/0xsj/atelier-wails/internal/host/wails"
)

func Run(assets embed.FS) error {
	log, err := consoleLogger()
	if err != nil {
		return err
	}
	return runLogged(log, func() error { return wailshost.Run(assets) })
}
