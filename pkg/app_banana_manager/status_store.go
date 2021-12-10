package app_banana_manager

type Status string

const (
	StatusComposingTweet = "ComposingTweet"
)

type ChatStatusStore struct {
	Status map[int64]Status
}

func (s *ChatStatusStore) GetRaw(chatID int64) (Status, bool) {
	status, ok := s.Status[chatID]
	return status, ok
}

func (s *ChatStatusStore) Get(chatID int64) Status {
	return s.Status[chatID]
}

func (s *ChatStatusStore) Set(chatID int64, status Status) {
	s.Status[chatID] = status
}

func (s *ChatStatusStore) Unset(chatID int64) {
	delete(s.Status, chatID)
}

func NewChatStatusStore() *ChatStatusStore {
	return &ChatStatusStore{
		Status: map[int64]Status{},
	}
}
