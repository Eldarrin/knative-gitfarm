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

func HandleWorkflowJob(payload github.WorkflowJobPayload) {
	log.Print("Handling Workflow Job")

	githubUrl := "https://github.com/" + payload.Repository.FullName

	// do enterprise stuff
	if !strings.Contains(payload.WorkflowJob.URL, "api.github.com") {
		githubUrl = "this will break"
	}

	// do organisatino stuff
	//if payload.Organization != (github.WorkflowJobPayload{}.Organization) {
	//	organization := payload.Organization.Login
	//}

	// do runner groups stuff
	runnerGroup := ""

	// clear and mash up labels
	labels := "docker"

	workDir := "/runner"

	runnerName := labels + "-" + RandString(8)

	configApp := "config.sh"

	flag.Parse()
	personalAccessToken := os.Getenv(personalAccessTokenKey)
	if personalAccessToken == "" {
		log.Fatal("Unauthorized: No token present")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: personalAccessToken})
	tc := oauth2.NewClient(ctx, ts)
	client := ghclient.NewClient(tc)

	runnerToken, _, err := client.Actions.CreateRegistrationToken(ctx, "eldarrin", "knative-gitfarm")
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

	cmdRun := exec.Command("run.sh")

	stdout, err = cmdRun.Output()

	if err != nil {
		log.Print(err.Error())
		return
	}

	log.Print(string(stdout))

}

func main() {
	flag.Parse()
	log.Print("gitwebhook sample started.")
	secretToken := os.Getenv(webhookSecretKey)

	hook, _ := github.New(github.Options.Secret(secretToken))

	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		log.Print("in handle")
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
			HandleWorkflowJob(job)

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
