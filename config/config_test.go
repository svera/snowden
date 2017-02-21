package config

import "testing"

func TestLoad(t *testing.T) {
	var err error

	if _, err = Load("non-existing-file.yaml"); err == nil {
		t.Errorf("Trying to load a non existing configuration file must return an error.")
	}
}

func TestParseNotValidYAML(t *testing.T) {
	var err error
	data := []byte(`@`)

	if _, err = Parse(data); err == nil {
		t.Errorf("Trying to parse a non valid YAML must return an error.")
	}
}

func TestValidateToken(t *testing.T) {
	var err error
	data := []byte(``)

	_, err = Parse(data)
	if err.Error() != "Invalid token." {
		t.Errorf("Configuration file must return an error if no Slack token is defined.")
	}
}

func TestValidateRepos(t *testing.T) {
	var err error
	// IMPORTANT: REMEMBER TO INDENT YAML DATA WITH SPACES, NOT TABS
	data := []byte(
		`
        slack_token: TEST
        `,
	)

	_, err = Parse(data)
	if err.Error() != "There must be at least one watched repository." {
		t.Errorf("Configuration file must return an error if there are no watched repositories defined.")
	}
}

func TestValidateRules(t *testing.T) {
	var err error
	// IMPORTANT: REMEMBER TO INDENT YAML DATA WITH SPACES, NOT TABS
	data := []byte(
		`
        debug: false
        slack_token: TEST
        watched:
            test-repo:
        `,
	)

	_, err = Parse(data)
	if err.Error() != "There are no rules defined." {
		t.Errorf("Configuration file must return an error if there are no rules defined.")
	}
}

func TestValidateWatchedFiles(t *testing.T) {
	var err error
	// IMPORTANT: REMEMBER TO INDENT YAML DATA WITH SPACES, NOT TABS
	data := []byte(
		`
        debug: false
        slack_token: TEST
        watched:
            test-repo:
                -
                    subscribers:
                        ["Fulanito"]
        `,
	)

	_, err = Parse(data)
	if err.Error() != "Rule 0 for repository test-repo has no watched files." {
		t.Errorf("Configuration file must return an error if a rule has no watched files.")
	}
}

func TestValidateSubscribers(t *testing.T) {
	var err error
	// IMPORTANT: REMEMBER TO INDENT YAML DATA WITH SPACES, NOT TABS
	data := []byte(
		`
        debug: false
        slack_token: TEST
        watched:
            test-repo:
                -
                    names:
                        - folder_b/folder_b1
        `,
	)

	_, err = Parse(data)
	if err.Error() != "Rule 0 for repository test-repo has no subscribers." {
		t.Errorf("Configuration file must return an error if a rule has no subscribers.")
	}
}
