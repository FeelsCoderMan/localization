package main

import (
	"fmt"
	"log"

	"github.com/localization/internal/pkg/localizer"
	"github.com/localization/internal/pkg/localizer/message"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	supportedLocales := []localizer.Locale{
		localizer.LOCALE_TR,
	}
	loc, err := localizer.New(localizer.LOCALE_EN, supportedLocales)

	if err != nil {
		return fmt.Errorf("localizer can not be initialized with error: %w", err)
	}

	log.Println("Printing user.logged.in key for default locale:")

	userLoggedInEn, err := loc.Translate(message.USER_LOGGED_IN, localizer.LOCALE_EN)
	if err != nil {
		return err
	}

	userLoggedInTr, err := loc.Translate(message.USER_LOGGED_IN, localizer.LOCALE_TR)
	if err != nil {
		return err
	}

	log.Printf("en_EN: %s\n", userLoggedInEn)
	log.Printf("tr_TR: %s\n", userLoggedInTr)
	return nil
}
