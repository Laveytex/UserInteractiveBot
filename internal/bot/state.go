package bot

import "UserInteractiveBot/internal/models"

type State struct {
	AdminStates   map[int64]models.AdminState
	CreatingPosts map[int64]*models.CreatingPost
}

func NewState() *State {
	return &State{
		AdminStates:   make(map[int64]models.AdminState),
		CreatingPosts: make(map[int64]*models.CreatingPost),
	}
}

func (s *State) GetAdminState(chatID int64) models.AdminState {
	if state, exists := s.AdminStates[chatID]; exists {
		return state
	}
	return models.StateIdle
}

func (s *State) SetAdminState(chatID int64, state models.AdminState) {
	s.AdminStates[chatID] = state
}

func (s *State) GetCreatingPost(chatID int64) *models.CreatingPost {
	if post, exists := s.CreatingPosts[chatID]; exists {
		return post
	}
	return nil
}

func (s *State) SetCreatingPost(chatID int64, post *models.CreatingPost) {
	s.CreatingPosts[chatID] = post
}

func (s *State) ClearState(chatID int64) {
	delete(s.AdminStates, chatID)
	delete(s.CreatingPosts, chatID)
}
