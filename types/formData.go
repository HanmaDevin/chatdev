package types

type FormData struct {
	Username string
	Password string
	Errors   FormDataError
}

type FormDataError struct {
	UsernameError string
	PasswordError string
}
