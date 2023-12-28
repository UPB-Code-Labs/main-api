package definitions

import "github.com/UPB-Code-Labs/main-api/src/languages/domain/entities"

type LanguagesRepository interface {
	GetAll() (languages []*entities.Language, err error)
	GetTemplateUUIDByLanguageUUID(uuid string) (templateUUID string, err error)
	GetTemplateBytes(uuid string) (template []byte, err error)
}
