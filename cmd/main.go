package main

import (
	"context"
	"flag"
	ghClient "github.com/google/go-github/v45/github"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

const (
	path    = "/jobs"
	health  = "/health"
	ready   = "/ready"
	workDir = "/runner"
	// Personal Access Token created in GitHub that allows us to make calls into GitHub.
	personalAccessTokenKey = "GITHUB_PERSONAL_TOKEN"
)

type JobInfo struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	CallingURL   string `json:"calling_url"`
	Labels       string `json:"labels"`
	RepoName     string `json:"repo_name"`
	RunnerGroup  string `json:"runner_group"`
	Owner        string `json:"owner"`
	Organization string `json:"organization"`
}

var available = true

func getRunnerToken(ctx context.Context, jobInfo *JobInfo) string {
	personalAccessToken := os.Getenv(personalAccessTokenKey)
	if personalAccessToken == "" {
		log.Fatal("Unauthorized: No token present")
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: personalAccessToken})
	tc := oauth2.NewClient(ctx, ts)
	client := ghClient.NewClient(tc)

	runnerTokenValue := ""

	if jobInfo.Organization != "" {
		runnerToken, _, err := client.Actions.CreateOrganizationRegistrationToken(ctx, jobInfo.Organization)
		if err != nil {
			log.Fatal(err)
		}
		runnerTokenValue = *runnerToken.Token
	} else {
		runnerToken, _, err := client.Actions.CreateRegistrationToken(ctx, jobInfo.Owner, jobInfo.RepoName)
		if err != nil {
			log.Fatal(err)
		}
		runnerTokenValue = *runnerToken.Token
	}

	return runnerTokenValue
}

func configureRunner(ctx context.Context, jobInfo *JobInfo) {
	githubUrl := "https://github.com/" + jobInfo.Owner + "/" + jobInfo.Name

	// do enterprise stuff
	if !strings.Contains(jobInfo.CallingURL, "api.github.com") {
		githubUrl = "this will break"
	}

	// clear and mash up labels
	// strip non-essentials
	labels := jobInfo.Labels

	runnerName := labels + "-" + RandString(8)

	configApp := workDir + "/config.sh"

	// do organisation stuff
	if jobInfo.Organization != "" {
		githubUrl = "https://github.com/" + jobInfo.Organization + "/" + jobInfo.Name
	}

	runnerTokenValue := getRunnerToken(ctx, jobInfo)
	args := []string{configApp,
		"--unattended",
		"--replace",
		"--name", runnerName,
		"--url", githubUrl,
		"--token", runnerTokenValue,
		"--labels", labels,
		"--work", workDir,
		"--ephemeral",
		"--disableupdate"}

	if jobInfo.RunnerGroup != "" {
		args = append(args, "--runnergroups", jobInfo.RunnerGroup)
	}

	cmdConfig := &exec.Cmd{
		Path:   configApp,
		Args:   args,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	log.Print(cmdConfig.String())

	if err := cmdConfig.Run(); err != nil {
		log.Print(err)
	}
}

func HandleWorkflowJob(ctx context.Context, jobInfo *JobInfo, ch chan<- string) {
	log.Print("Handling Workflow Job")
	log.Print(ch)

	configureRunner(ctx, jobInfo)

	cmdRun := &exec.Cmd{
		Path:   workDir + "/run.sh",
		Args:   []string{workDir + "/run.sh"},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	log.Print(ch)
	log.Print(cmdRun.String())

	if err := cmdRun.Run(); err != nil {
		log.Print(err)
	}
	log.Print(ch)
	available = true

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

func handler(w http.ResponseWriter, _ *http.Request) {
	log.Print("in handler")

	ctx := context.Background()

	ch := make(chan string)

	if available {
		// block this call so knative thinks its processing and doesn't kill mid-job
		available = false
		HandleWorkflowJob(ctx, newJob("knative-gitfarm"), ch)
	} else {
		// accept no more requests so it spawns a new agent
		w.WriteHeader(503)
		w.Write([]byte("Server is active, send it somewhere else"))
	}

}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func readyHandler(w http.ResponseWriter, _ *http.Request) {
	if available {
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	} else {
		w.WriteHeader(503)
		w.Write([]byte("KO"))
	}
}

func main() {
	flag.Parse()
	log.Print("gitrunner started.")

	http.HandleFunc(path, handler)
	http.HandleFunc(health, healthHandler)
	http.HandleFunc(ready, readyHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Cannot open port")
		return
	}
}
