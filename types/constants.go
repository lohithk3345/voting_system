package types

const API_KEY string = "x-api-key"
const TOKEN_SECRET string = "TOKEN_SECRET"

const (
	SESSION string = "SESSION"
	JOIN    string = "JOIN"
)

const EmptyString string = ""

const (
	VOTER Role = "VOTER"
	ADMIN Role = "ADMIN"
)

var Roles []Role = []Role{VOTER, ADMIN}
