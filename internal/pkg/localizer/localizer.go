package localizer

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"path"
	"slices"

	"github.com/localization/internal/pkg/localizer/message"
)

const (
	LOCALE_EN Locale = "en_EN"
	LOCALE_TR Locale = "tr_TR"

	DIR_LOCALES   = "locales"
	FILE_EXT_JSON = ".json"
)

//go:embed locales
var embeddedLocales embed.FS

type Locale string

type Localizer struct {
	defaultLocale    Locale
	supportedLocales []Locale
	messages         map[Locale]map[message.Message]string
}

func New(defaultLocale Locale, supportedLocales []Locale) (*Localizer, error) {
	if !slices.Contains(supportedLocales, defaultLocale) {
		supportedLocales = append(supportedLocales, defaultLocale)
	}

	sub, err := fs.Sub(embeddedLocales, DIR_LOCALES)
	if err != nil {
		return nil, fmt.Errorf("could not create sub-filesystem for the folder %s with error: %w", DIR_LOCALES, err)
	}

	return newFromFS(sub, defaultLocale, supportedLocales, message.All)
}

func newFromFS(fsys fs.FS, defaultLocale Locale,
	supportedLocales []Locale, allMessages []message.Message) (*Localizer, error) {
	loc := &Localizer{
		defaultLocale:    defaultLocale,
		supportedLocales: supportedLocales,
		messages:         make(map[Locale]map[message.Message]string, len(supportedLocales)),
	}

	var errs []error
	for _, locale := range supportedLocales {
		msgs, err := loadLocale(fsys, locale)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		if len(errs) == 0 {
			loc.messages[locale] = msgs
		}
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	for _, locale := range supportedLocales {
		for _, m := range allMessages {
			if _, ok := loc.messages[locale][m]; !ok {
				errs = append(errs, fmt.Errorf("locale %s is missing message of %s key", locale, m))
			}
		}
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return loc, nil
}

func loadLocale(fsys fs.FS, locale Locale) (map[message.Message]string, error) {
	entries, err := fs.ReadDir(fsys, string(locale))
	if err != nil {
		return nil, fmt.Errorf("could not read folder for locale %s with error: %w", locale, err)
	}

	msgs := make(map[message.Message]string)
	found := false

	for _, entry := range entries {
		if entry.IsDir() || path.Ext(entry.Name()) != FILE_EXT_JSON {
			continue
		}

		found = true
		file := path.Join(string(locale), entry.Name())
		data, err := fs.ReadFile(fsys, file)
		if err != nil {
			return nil, fmt.Errorf("could not read the file %s with error: %w", file, err)
		}

		var part map[string]string
		if err := json.Unmarshal(data, &part); err != nil {
			return nil, fmt.Errorf("found invalid JSON from file %s with error: %w", file, err)
		}

		// TODO: Keep the duplicated errors in array and return it as single
		for k, v := range part {
			key := message.Message(k)

			if _, dup := msgs[key]; dup {
				return nil, fmt.Errorf("found duplicated key from locale %s", locale)
			}

			msgs[key] = v
		}

	}

	if !found {
		return nil, fmt.Errorf("no %s files found in locale %s", FILE_EXT_JSON, locale)
	}

	return msgs, nil
}

func (l *Localizer) Translate(m message.Message, locale Locale) (string, error) {
	// TODO: After CI/CD pipeline confirms the all_gen.go generation, error case of the func will be unnecessary
	if !slices.Contains(l.supportedLocales, locale) {
		locale = l.defaultLocale
		// return "", fmt.Errorf("unsupported locale of %s is used", locale)
	}

	text, ok := l.messages[locale][m]
	if !ok {
		return "", fmt.Errorf("locale %s does not contain message key of %s", locale, m)
	}

	return text, nil
}
