package restapi

type APIResponse struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}
