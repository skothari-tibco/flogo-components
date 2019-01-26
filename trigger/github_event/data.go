package github_event

import (
	
	"github.com/google/go-github/github"
	"time"
)


type CommitInfo struct{
	
	Ref string 

	Author struct{
		Name string 
		Email string

	}
	
	Added []string
	Removed  []string
	Modified []string
	Message string
	CloneURL string
	UpdatedAt interface{}
	
}

type PullInfo struct{
	Action string
	User struct{
		Login string
	}
	Body string
	CreatedAt time.Time
	Assignees          []interface{}
	RequestedReviewers []interface{}
}

func (c *CommitInfo) Set(e *github.PushEvent) {

	c.Ref = *e.Ref
	c.Author.Name = *e.Commits[0].Author.Name

	c.Author.Email = *e.Commits[0].Author.Email

	c.Added = e.HeadCommit.Added

	c.Removed = e.HeadCommit.Removed

	c.Modified = e.HeadCommit.Modified

	c.Message = *e.Commits[0].Message

	c.CloneURL = *e.Repo.CloneURL

	c.UpdatedAt = e.Repo.UpdatedAt

	
}

func (c *PullInfo) Set(e *github.PullRequestEvent){
	
	c.Action = *e.Action

	c.User.Login = *e.PullRequest.User.Login

	c.Body = *e.PullRequest.Body

	c.CreatedAt = *e.PullRequest.CreatedAt

	//c.Assignees = e.PullRequest.Assignees

	//c.RequestedReviewers = *e.PullRequest.RequestedReviewers

}