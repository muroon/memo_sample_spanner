package spanner

import (
	"memo_sample_spanner/domain/app"
	"memo_sample_spanner/domain/model"
	"memo_sample_spanner/infra/cloudspanner"
	"memo_sample_spanner/infra/error"

	"github.com/google/uuid"
)

var errm apperror.ErrorManager

func init() {
	errm = apperror.NewErrorManager()
}

// yoRODB get YORODB instance
func yoRODB() model.YORODB {
	return cloudspanner.DB().(model.YORODB)
}

// generateID generate Key
func generateID() (string, error) {
	u4, err := uuid.NewRandom()
	if err != nil {
		return "", errm.Wrap(err, app.DBError)
	}

	return u4.String(), nil
}
