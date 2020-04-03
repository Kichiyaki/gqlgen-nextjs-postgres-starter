package i18n

import (
	"os"
	"path/filepath"

	_i18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var Bundle = _i18n.NewBundle(language.English)

func LoadMessageFiles(root string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if path != root {
			Bundle.MustLoadMessageFile(path)
		}
		return nil
	})
}
