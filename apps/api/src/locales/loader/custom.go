package loader

import "github.com/sigmaott/gest/package/extension/i18nfx/loader"

var en = []loader.Translation{
	{
		Key:      "cardinal_test",
		Trans:    "You have 1 day left.",
		Type:     "One",
		Rule:     "",
		Override: false,
	},
}

var vi = []loader.Translation{
	{
		Key:      "cardinal_test",
		Trans:    "You have 1 day left.",
		Type:     "One",
		Rule:     "",
		Override: false,
	},
}

type memoryLoader struct {
}

func (m *memoryLoader) LoadData() map[string]loader.ListTranslation {
	return map[string]loader.ListTranslation{
		"en": en,
		"vi": vi,
	}
}

func NewI18nMemoryLoader() loader.II18nLoader {
	return &memoryLoader{}
}
