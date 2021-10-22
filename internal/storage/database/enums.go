package database

const (
	SearchKey = "key"

	//search with key example firstname like '%a%'
	SearchWithKeyCondition = `(FirstName like concat('%%',?,'%%') OR
			MiddleName like concat('%%',?,'%%') OR
			LastName like concat('%%',?,'%%') OR
			FullName like concat('%%',?,'%%'))`
	//search args here we are comparing with four columns hence searchWithKey = 4
	SearchWithKeyArgs = 4
	PageSize          = "limit"
	PageNo            = "page"
	LastId            = "lastid"
)
