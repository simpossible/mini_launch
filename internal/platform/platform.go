package platform

import "github.com/simpossible/mini_launch/internal/service"

// Platform defines the interface for daemon management on different OS platforms.
type Platform interface {
	// Generate creates the daemon configuration file for a service.
	Generate(svc *service.Service) error

	// Start launches the service daemon.
	Start(svc *service.Service) error

	// Stop stops the service daemon.
	Stop(svc *service.Service) error

	// Restart stops then starts the service daemon.
	Restart(svc *service.Service) error

	// Status returns the current status of a service.
	Status(svc *service.Service) (string, error)

	// Remove stops the service and deletes its daemon configuration.
	Remove(svc *service.Service) error

	// IsConfigured checks if the daemon configuration exists for a service.
	IsConfigured(svc *service.Service) bool
}

// New returns the appropriate Platform implementation for the current OS.
func New() Platform {
	return newPlatform()
}
