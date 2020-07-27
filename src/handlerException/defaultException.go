package handlerException

type DefaultError struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"mensagem"`
}
