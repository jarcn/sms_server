package producer

type Msg struct {
	Topic         string            `json:"topic"`
	Title         string            `json:"title"`
	Content       string            `json:"content"`
	Summary       string            `json:"summary"`
	FromUserId    uint32            `json:"fromUserId"`
	FromUserName  string            `json:"fromUserName"`
	ToUserId      uint64            `json:"toUserId"`
	ToUserName    string            `json:"toUserName"`
	MsgType       uint32            `json:"msgType"`
	KeyPair       map[string]string `json:"keyPair"`
	PushChanneles []string          `json:"pushChanneles"`
}
