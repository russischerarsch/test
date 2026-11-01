package e2e

type CreatePairReq struct {
	Pair string `json:"Pair"`
}

type CreatePairResp struct {
	ID int64 `json:"id"`
}

type Quote struct {
	ID   string `json:"id"`
	Pair string `json:"pair"`
}

type ErrorResp struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
