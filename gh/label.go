package gh

import (
	"fmt"

	"github.com/google/go-github/github"
)

// Issue Colors
const (
	IssueColorRed    = "b60205"
	IssueColorOrange = "d93f0b"
	IssueColorYellow = "fbca04"
	IssueColorGreen  = "0e8a16"
	IssueColorCyan   = "006b75"
	IssueColorMarine = "1d76db"
	IssueColorBlue   = "0052cc"
	IssueColorPurple = "5319e7"

	IssueColorWhite = "ffffff"
	IssueColorGrey  = "e6e6e6"

	IssueColorLightRed    = "e99695"
	IssueColorLightOrange = "f9d0c4"
	IssueColorLightYellow = "fef2c0"
	IssueColorLightGreen  = "c2e0c6"
	IssueColorLightCyan   = "bfdadc"
	IssueColorLightMarine = "c5def5"
	IssueColorLightBlue   = "bfd4f2"
	IssueColorLightPurple = "d4c5f9"
)

var (
	// DefaultLabels for all repositories held by FlowUp
	DefaultLabels = []*github.Label{
		// estimations
		NewLabel("estimated: 0.5", IssueColorLightCyan),
		NewLabel("estimated: 1", IssueColorLightCyan),
		NewLabel("estimated: 2", IssueColorLightCyan),
		NewLabel("estimated: 3", IssueColorLightCyan),
		NewLabel("estimated: 5", IssueColorLightCyan),
		NewLabel("estimated: 8", IssueColorLightCyan),
		NewLabel("estimated: 13", IssueColorLightCyan),

		// consumptions
		NewLabel("consumed: 0.5", IssueColorCyan),
		NewLabel("consumed: 1", IssueColorCyan),
		NewLabel("consumed: 2", IssueColorCyan),
		NewLabel("consumed: 3", IssueColorCyan),
		NewLabel("consumed: 5", IssueColorCyan),
		NewLabel("consumed: 8", IssueColorCyan),
		NewLabel("consumed: 13", IssueColorCyan),

		// priorities
		NewLabel("priority: low", IssueColorGreen),
		NewLabel("priority: medium", IssueColorYellow),
		NewLabel("priority: high", IssueColorRed),

		// merging
		NewLabel("merge: ready", IssueColorGreen),
		NewLabel("merge: conflicts", IssueColorRed),

		// misc
		NewLabel("enhancement", IssueColorBlue),
		NewLabel("bug", IssueColorRed),
		NewLabel("fix", IssueColorLightGreen),
		NewLabel("test", IssueColorLightRed),
		NewLabel("e2e", IssueColorLightCyan),
		NewLabel("refactor", IssueColorMarine),
		NewLabel("IE", IssueColorBlue),
		NewLabel("screw IE", IssueColorBlue),
		NewLabel("ux", IssueColorLightYellow),
		NewLabel("optimization", IssueColorMarine),
		NewLabel("security", IssueColorRed),
		NewLabel("documentation", IssueColorGreen),
		NewLabel("on hold", IssueColorLightMarine),
		NewLabel("deferred", IssueColorLightRed),
		NewLabel("help wanted", IssueColorGreen),
		NewLabel("duplicate", IssueColorWhite),
		NewLabel("invalid", IssueColorGrey),
		NewLabel("wontfix", IssueColorWhite),
		NewLabel("question", IssueColorPurple),

		// stages
		NewLabel("discussion", IssueColorYellow),
		NewLabel("design", IssueColorYellow),
		NewLabel("testing", IssueColorYellow),
		NewLabel("staging", IssueColorYellow),
		NewLabel("production", IssueColorYellow),
	}
)

// NewLabel creates a new label without setting url
func NewLabel(name, color string) *github.Label {
	return &github.Label{
		Name:  &name,
		Color: &color,
	}
}

// GithubLabelService
type GithubLabelService struct {
	client *github.Client
}

func NewGithubLabelService(client *github.Client) *GithubLabelService {
	return &GithubLabelService{client}
}

func (s *GithubLabelService) CreateLabel(owner, repo string, label *github.Label) error {
	// set url for the label
	issueURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/labels/%s", owner, repo, *label.Name)
	label.URL = &issueURL

	_, _, err := s.client.Issues.CreateLabel(owner, repo, label)
	return err
}

func (s *GithubLabelService) CreateLabels(owner, repo string, labels []*github.Label) []error {
	var errors []error

	for _, label := range labels {
		err := s.CreateLabel(owner, repo, label)
		if err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}

func (s *GithubLabelService) RemoveLabel(owner, repo, name string) error {
	_, err := s.client.Issues.DeleteLabel(owner, repo, name)
	return err
}

func (s *GithubLabelService) RemoveAllLabels(owner, repo string) []error {
	var errors []error
	labels, _, err := s.client.Issues.ListLabels(owner, repo, nil)
	if err != nil {
		errors = append(errors, err)
		return errors
	}

	for _, label := range labels {
		err := s.RemoveLabel(owner, repo, *label.Name)
		if err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}
