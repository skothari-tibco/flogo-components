package github_event

import (
	"context"
	"fmt"
	"io/ioutil"

	"net/http"

	"github.com/google/go-github/github"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/trigger"
)

var triggerMd = trigger.NewMetadata(&Settings{}, &HandlerSettings{}, &Output{}, &Reply{})

func init() {
	trigger.Register(&GitHubEvent{}, &Factory{})
}

type GitHubEvent struct {
	settings *Settings
}

type Factory struct {
}

// Metadata implements trigger.Factory.Metadata
func (*Factory) Metadata() *trigger.Metadata {
	return triggerMd
}

func (t *Factory) New(config *trigger.Config) (trigger.Trigger, error) {
	s := &Settings{}
	err := metadata.MapToStruct(config.Settings, s, true)
	if err != nil {
		return nil, err
	}
	return &GitHubEvent{settings: s}, nil
}

func (t *GitHubEvent) Initialize(ctx trigger.InitContext) error {

	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		payload, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("error reading request body: err=%s\n", err)
			return
		}

		event, err := github.ParseWebHook(github.WebHookType(r), payload)
		if err != nil {
			fmt.Printf("could not parse webhook: err=%s\n", err)
			return
		}

		switch e := event.(type) {

		case *github.PushEvent:

			var user_info CommitInfo
			fmt.Println("Detected Push Event")

			user_info.Set(e)
			RegisterHandlers(ctx, e, user_info)

		case *github.PullRequestEvent:

			fmt.Println("Pull Request Detected")

			var user_info PullInfo

			user_info.Set(e)
			RegisterHandlers(ctx, e, user_info)

		case *github.WatchEvent:

			if e.Action != nil && *e.Action == "starred" {
				fmt.Printf("%s starred repository %s\n",
					*e.Sender.Login, *e.Repo.FullName)
			}
		default:
			fmt.Printf("unknown event type %s\n", github.WebHookType(r))
			return
		}

	})

	return nil

}

func (t *GitHubEvent) Start() error {
	return http.ListenAndServe(t.settings.Port, nil)
}

// Stop implements util.Managed.Stop
func (t *GitHubEvent) Stop() error {
	return nil
}

func RegisterHandlers(ctx trigger.InitContext, e interface{}, user_info interface{}) {

	for _, handler := range ctx.GetHandlers() {

		out := &Output{}
		s := &HandlerSettings{}
		err := metadata.MapToStruct(handler.Settings(), s, true)
		if err != nil {
			return
		}

		if s.GetAll {
			fmt.Println("Sending all data")
			out.Content = e
			_, err := handler.Handle(context.Background(), out)

			if err != nil {
				return
			}

		} else {

			fmt.Println("Sending user data")

			out.Content = user_info

			_, err := handler.Handle(context.Background(), out)

			if err != nil {
				return
			}
		}

	}
}
