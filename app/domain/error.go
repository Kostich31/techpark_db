package domain

const (
	NoSlug       = "Can't find slug\n"
	NoUser       = "Can't find user\n"
	BadParentPost = "Parent post was created in another thread\n"
	ConflictData = "Conflict data\n"
)

const (
	PgxBadParentErrorCode = "77777"
	PgxNoFoundFieldErrorCode = "23503"
	PgxUniqErrorCode         = "23505"
)

type CustomError struct {
	Message string `json:"message"`
}
