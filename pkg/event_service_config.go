package pkg

import "log"

type EventServiceConfig struct {
	ErrorHandler ErrorHandler
}

func GetDefaultConfig() EventServiceConfig {
	return EventServiceConfig{
		ErrorHandler: defaultErrorHandler,
	}
}

func defaultErrorHandler(err error) error {
	if err != nil {
		log.Println(err)
	}

	return nil
}
