package server

//MemoryStruct contains detailed userInfo
type MemoryStruct struct {
	SlackID      string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	DeleteStatus bool   `json:"status"`
}
