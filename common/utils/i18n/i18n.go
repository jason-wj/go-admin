package i18n

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

type I18n struct {
	Path   string
	Lang   string
	Source [][]string
	Data   map[string]string
}

var (
	EnI18nClient = &I18n{}
)

//用于支持现有系统,当前支持英文
func InitLang(){
	var err error
	EnI18nClient, err = newI18n("config/lang", "en")
	if err != nil {
		log.Fatalln(err)
	}
}

func newI18n(path string, lang string) (*I18n, error) {
	fileName := fmt.Sprintf("%s/%s.csv", path, lang)
	fs1, err := os.Open(fileName)

	if err != nil {
		return nil, err
	}

	r1 := csv.NewReader(fs1)
	r1.Comma = ','
	r1.FieldsPerRecord = -1
	content, err := r1.ReadAll()
	if err != nil {
		return nil, err
	}

	i18n := I18n{}
	i18n.Path = path
	i18n.Lang = lang
	i18n.Source = content
	i18n.Data = make(map[string]string)
	return &i18n, nil
}

func (i *I18n) T(key string, args ...interface{}) string {
	format := key

	if _, ok := i.Data[key]; ok {
		format = i.Data[key]
	} else {
		for _, row := range i.Source {
			if row[0] == key {
				i.Data[key] = row[1]
				format = row[1]
				break
			}
		}
	}
	format = i.preArgs(format, args...)
	return format
}

//Choose language translation
func (i *I18n) TOption(key string, lang string, args ...interface{}) string {
	i18nClient, err := newI18n(i.Path, lang)
	format := key

	if err != nil {
		return i.preArgs(format, args...)
	}

	for _, row := range i18nClient.Source {
		if row[0] == key {
			i18nClient.Data[key] = row[1]
			format = row[1]
			break
		}
	}

	format = i.preArgs(format, args...)
	return format
}

func (i *I18n) preArgs(format string, args ...interface{}) string {
	if len(args) > 0 {
		format = fmt.Sprintf(format, args...)
	}
	return format
}
