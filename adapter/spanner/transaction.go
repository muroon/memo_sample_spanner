package spanner

import (
	"memo_sample_spanner/domain/transaction"
	"memo_sample_spanner/infra/cloudspanner"
)

func NewTransaction() transaction.ITransaction {
	return cloudspanner.DB()
}
