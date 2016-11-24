package logic

import (
	"testing"

	"github.com/google/go-github/github"
	"github.com/svera/snowden/config"
	"github.com/svera/snowden/mocks"
)

var (
	sender *mocks.Sender
	gh     *mocks.Github
	cfg    *config.Config
)

func setup() {
	fileName1 := "test_folder_a/test_file"
	fileName2 := "test_folder_b/test_file"
	sender = &mocks.Sender{}
	gh = &mocks.Github{
		FakeCommitFiles: []*github.CommitFile{
			&github.CommitFile{
				Filename: &fileName1,
			},
			&github.CommitFile{
				Filename: &fileName2,
			},
		},
	}
	cfg = &config.Config{
		SlackToken: "TEST",
		Watched: map[string][]config.Rule{
			"test_repo": {
				config.Rule{
					Names:       []string{"test_folder_a", "test_folder_b"},
					Subscribers: []string{"fulanito"},
				},
			},
		},
	}
}

func TestProcess(t *testing.T) {
	setup()

	lgc := New(sender, gh, cfg)
	lgc.Process("opened", "owner", "test_repo", 1, "title", "description")
	if sender.Calls["PostMessage"] != 1 {
		t.Errorf("Snowden should have sent one message, %d messages delivered", sender.Calls["PostMessage"])
	}
}

func TestProcessNotSubscribedToFolder(t *testing.T) {
	setup()

	cfg.Watched["test_repo"][0].Names = []string{"test_folder_c"}
	lgc := New(sender, gh, cfg)
	lgc.Process("opened", "owner", "test_repo", 1, "title", "description")
	if sender.Calls["PostMessage"] != 0 {
		t.Errorf("Snowden should not have sent any message, %d messages delivered", sender.Calls["PostMessage"])
	}
}

func TestProcessNotSubscribedToRepository(t *testing.T) {
	setup()

	lgc := New(sender, gh, cfg)
	lgc.Process("opened", "owner", "no_subscribed_repo", 1, "title", "description")
	if sender.Calls["PostMessage"] != 0 {
		t.Errorf("Snowden should not have sent any message, %d messages delivered", sender.Calls["PostMessage"])
	}
}

func TestProcessDoesNotNotifyIfPROwnerIsInExceptionsList(t *testing.T) {
	setup()

	cfg.Watched["test_repo"][0].Exceptions = []string{"owner"}
	lgc := New(sender, gh, cfg)
	lgc.Process("opened", "owner", "test_repo", 1, "title", "description")
	if sender.Calls["PostMessage"] != 0 {
		t.Errorf("Snowden should not have sent any message, %d messages delivered", sender.Calls["PostMessage"])
	}
}

func TestProcessOnlyIfOpenedOrReopened(t *testing.T) {
	setup()

	lgc := New(sender, gh, cfg)
	lgc.Process("closed", "owner", "test_repo", 1, "title", "description")
	if sender.Calls["PostMessage"] != 0 {
		t.Errorf("Snowden should not have sent any message, %d messages delivered", sender.Calls["PostMessage"])
	}
}
