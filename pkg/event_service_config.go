package events

import "log"

type EventServiceConfig struct {
	ErrorHandler ErrorHandler
}

// Returns default config for events service
func GetDefaultConfig() EventServiceConfig {
	return EventServiceConfig{
		ErrorHandler: defaultErrorHandler,
	}
}

// defaultErrorHandler simply logs the error
func defaultErrorHandler(err error) error {
	if err != nil {
		log.Println(err)
	}

	return nil
}
