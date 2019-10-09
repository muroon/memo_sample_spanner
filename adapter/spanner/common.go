package spanner

import (
	"memo_sample_spanner/domain/model"
	"memo_sample_spanner/infra/cloudspanner"

	"github.com/google/uuid"
)

// yoRODB get YORODB instance
func yoRODB() model.YORODB {
	return cloudspanner.DB().(model.YORODB)
}

// generateID generate Key
func generateID() (string, error) {
	u4, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	return u4.String(), nil
}
