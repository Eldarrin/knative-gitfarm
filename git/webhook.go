package main

import (
	"context"
	"flag"
	"github.com/go-playground/webhooks/v6/github"
	ghclient "github.com/google/go-github/v45/github"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

const (
	path = "/webhooks"
	// Secret given to github. Used for verifying the incoming objects.
	personalAccessTokenKey = "GITHUB_PERSONAL_TOKEN"
	// Personal Access Token created in github that allows us to make
	// calls into github.
	webhookSecretKey = "WEBHOOK_SECRET"
)

type JobInfo struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	CallingURL  string `json:"calling_url"`
	Labels      string `json:"labels"`
	RepoName    string `json:"repo_name"`
	RunnerGroup string `json:"runner_group"`
	Owner       string `json:"owner"`
}

func HandleWorkflowJob(jobInfo *JobInfo) {
	log.Print("Handling Workflow Job")

	githubUrl := "https://github.com/" + jobInfo.Owner + "/" + jobInfo.Name

	// do enterprise stuff
	if !strings.Contains(jobInfo.CallingURL, "api.github.com") {
		githubUrl = "this will break"
	}

	// do organisatino stuff
	//if payload.Organization != (github.WorkflowJobPayload{}.Organization) {
	//	organization := payload.Organization.Login
	//}

	// do runner groups stuff
	runnerGroup := jobInfo.RunnerGroup

	// clear and mash up labels
	// strip non-essentials
	labels := jobInfo.Labels

	workDir := "/runner"

	runnerName := labels + "-" + RandString(8)

	configApp := workDir + "/config.sh"

	flag.Parse()
	personalAccessToken := os.Getenv(personalAccessTokenKey)
	if personalAccessToken == "" {
		log.Fatal("Unauthorized: No token present")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: personalAccessToken})
	tc := oauth2.NewClient(ctx, ts)
	client := ghclient.NewClient(tc)

	runnerToken, _, err := client.Actions.CreateRegistrationToken(ctx, jobInfo.Owner, jobInfo.RepoName)
	if err != nil {
		log.Fatal(err)
	}

	runnerTokenValue := *runnerToken.Token

	runnerName = "--name " + runnerName
	githubUrl = "--url " + githubUrl
	runnerTokenValue = "--token " + runnerTokenValue
	if runnerGroup != "" {
		runnerGroup = "--runnergroup " + runnerGroup
	}
	labels = "--labels " + labels
	workDir = "--work " + workDir

	cmd := exec.Command(configApp, "--unattended", "--replace", runnerName, githubUrl, runnerTokenValue, runnerGroup, labels, workDir, "--ephemeral", "--disableupdate")

	stdout, err := cmd.Output()

	if err != nil {
		log.Print(err.Error())
		return
	}

	log.Print(string(stdout))

	cmdRun := exec.Command(workDir + "/run.sh")

	stdout, err = cmdRun.Output()

	if err != nil {
		log.Print(err.Error())
		return
	}

	log.Print(string(stdout))

}

func newJob(name string) *JobInfo {
	j := JobInfo{Name: name}
	j.ID = 1
	j.Labels = "main"
	j.CallingURL = "https://api.github.com"
	j.Owner = "eldarrin"
	j.RepoName = "knative-gitfarm"
	return &j
}

func main() {
	flag.Parse()
	log.Print("gitwebhook sample started.")
	secretToken := os.Getenv(webhookSecretKey)

	hook, _ := github.New(github.Options.Secret(secretToken))

	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		log.Print("in handle")

		HandleWorkflowJob(newJob("knative-gitfarm"))
		payload, err := hook.Parse(r, github.WorkflowJobEvent, github.ReleaseEvent, github.PullRequestEvent)
		if err != nil {
			if err == github.ErrEventNotFound {
				log.Print("GitHub Event not found.")
				// ok event wasn;t one of the ones asked to be parsed
			}
		}
		switch payload.(type) {

		case github.WorkflowJobPayload:
			job := payload.(github.WorkflowJobPayload)
			log.Print("%+v", job)

		case github.ReleasePayload:
			release := payload.(github.ReleasePayload)
			// Do whatever you want from here...
			log.Print("%+v", release)

		case github.PullRequestPayload:
			pullRequest := payload.(github.PullRequestPayload)
			// Do whatever you want from here...
			log.Print("%+v", pullRequest)
		}
	})
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Cannot open port")
		return
	}
}
