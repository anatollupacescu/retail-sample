package acceptance

import (
	"fmt"

	faker "github.com/bxcodec/faker/v3"

	"github.com/google/uuid"
)

func Word() string {
	uid := uuid.New()
	word := faker.Word()
	return fmt.Sprintf("%s-%s", word, uid)[0:36]
}
