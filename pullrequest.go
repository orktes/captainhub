package main

import "github.com/google/go-github/github"

func createPullRequestComment(owner string, repo string, prNumber int, body string) (err error) {
	client := getGithubClient()
	_, _, err = client.PullRequests.CreateComment(owner, repo, prNumber, &github.PullRequestComment{
		Body: &body,
	})

	return
}

func createStatus(owner string, repo string, sha string, state string, targetURL string, description string, context string) (err error) {
	client := getGithubClient()
	_, _, err = client.Repositories.CreateStatus(owner, repo, sha, &github.RepoStatus{
		State:       &state,
		TargetURL:   &targetURL,
		Description: &description,
		Context:     &context,
	})

	return
}

func createIssueComment(owner string, repo string, issueNumber int, body string) (err error) {
	client := getGithubClient()

	_, _, err = client.Issues.CreateComment(owner, repo, issueNumber, &github.IssueComment{
		Body: &body,
	})
	return
}

func getPullRequestDetails(owner string, repo string, prNumber int) (pullRequest *github.PullRequest, err error) {
	client := getGithubClient()
	pullRequest, _, err = client.PullRequests.Get(owner, repo, prNumber)
	return
}

func listPullRequestFiles(owner string, repo string, prNumber int) (allFiles []github.CommitFile, err error) {
	client := getGithubClient()
	opts := &github.ListOptions{PerPage: 10}
	var resp *github.Response
	var files []github.CommitFile

	for {
		files, resp, err = client.PullRequests.ListFiles(
			owner,
			repo,
			prNumber, opts)
		if err != nil {
			return
		}
		allFiles = append(allFiles, files...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return
}
