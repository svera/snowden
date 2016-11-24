package mocks

import "github.com/nlopes/slack"

type Sender struct {
	FakePostMessageResponseChannel   string
	FakePostMessageResponseTimestamp string
	FakePostMessageResponseError     error
	Calls                            map[string]int
}

func (s *Sender) PostMessage(channel, text string, params slack.PostMessageParameters) (string, string, error) {
	if s.Calls == nil {
		s.Calls = map[string]int{}
	}
	s.Calls["PostMessage"]++
	return s.FakePostMessageResponseChannel, s.FakePostMessageResponseTimestamp, s.FakePostMessageResponseError
}
