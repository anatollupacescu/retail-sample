package helper

import (
	"fmt"

	faker "github.com/bxcodec/faker/v3"

	"github.com/google/uuid"
)

func Name() string {
	uid := uuid.New()
	word := faker.Word()

	return fmt.Sprintf("%s-%s", word, uid)[0:36]
}
