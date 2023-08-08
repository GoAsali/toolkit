package multilingual

import (
	filesUtils "github.com/goasali/toolkit/utils/files"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	log "github.com/sirupsen/logrus"
	"golang.org/x/text/language"
)

var bundle *i18n.Bundle

type Multilingual struct {
	*i18n.Bundle
	Path string
}

func NewMultilingual(b *i18n.Bundle, path string) *Multilingual {
	bundle = b
	return &Multilingual{Bundle: bundle, Path: path}
}

func Bundle() *i18n.Bundle {
	return bundle
}

// Load Load all messages from languages folder
func (m *Multilingual) Load() error {
	defer func() {
		log.Println("complete")
	}()

	log.Println("Loading multilingual list")

	dlMaker := directoryToDirLanguage(m.Path)
	files := make([]fileLanguage, 0)

	var messagesMap map[language.Tag][]*i18n.Message
	var err error

	for _, dir := range filesUtils.Directories(m.Path) {
		dl := dlMaker(dir)
		files = append(files, dl.files...)
	}
	if messagesMap, err = parseFiles(files); err != nil {
		return err
	}

	for lang, messages := range messagesMap {
		if err := m.AddMessages(lang, messages...); err != nil {
			return err
		}
	}

	return nil
}

// ChangeLanguageDirectory Change language directory and reload messages files
func (m *Multilingual) ChangeLanguageDirectory(dirPath string) error {
	m.Path = dirPath
	return m.Load()
}
