package enum

type Lang string

const (
	LangUz Lang = "uz"
	LangRu Lang = "ru"
	LangEn Lang = "en"
)

const DefaultLang = LangUz

func (this Lang) String() string { return string(this) }
