package env

type session_tokens struct {
	id        int
	selector  string
	validator string
	userId    int
	exp       int64
}
