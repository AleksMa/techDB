package models

type Chat struct {
	ID            uint64   `json:"id"`
	Name          string   `json:"name"`
	TotalMSGCount int64    `json:"-"`
	Members       []uint64 `json:"members"`
	LastMessage   string   `json:"last_message"`
}

type ResponseChatsArray struct {
	Chats      []Chat      `json:"chats"`
	Workspaces []Workspace `json:"workspaces"`
}

func NewChatModel(Name string, ID1 uint64, ID2 uint64) *Chat {
	return &Chat{
		ID:            0,
		Name:          Name,
		TotalMSGCount: 0,
		Members:       []uint64{ID1, ID2},
		LastMessage:   "",
	}
}
