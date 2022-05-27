package runnercontrol

import (
	"context"
	"flag"
	ghclient "github.com/google/go-github/v45/github"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

const (
	path = "/jobs"
	health = "/health"
	ready = "/ready"
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

var available = true

func HandleWorkflowJob(ctx context.Context, jobInfo *JobInfo, ch chan<- string) {
	log.Print("Handling Workflow Job")
	log.Print(ch)

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
	//runnerGroup := jobInfo.RunnerGroup

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

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: personalAccessToken})
	tc := oauth2.NewClient(ctx, ts)
	client := ghclient.NewClient(tc)

	runnerToken, _, err := client.Actions.CreateRegistrationToken(ctx, jobInfo.Owner, jobInfo.RepoName)
	if err != nil {
		log.Fatal(err)
	}

	runnerTokenValue := *runnerToken.Token

	cmdConfig := &exec.Cmd{
		Path: configApp,
		Args: []string{configApp,
			"--unattended",
			"--replace",
			"--name", runnerName,
			"--url", githubUrl,
			"--token", runnerTokenValue,
			"--labels", labels,
			"--work", workDir,
			"--ephemeral",
			"--disableupdate"},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	log.Print(cmdConfig.String())

	if err := cmdConfig.Run(); err != nil {
		log.Print(err)
	}

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

	ctx := context.Context(context.Background())

	ch := make(chan string)

	if available {
		available = false
		w.WriteHeader(200)
		w.Write([]byte("Got the request and am spinning up"))

		go HandleWorkflowJob(ctx, newJob("knative-gitfarm"), ch)
	} else {
		w.WriteHeader(503)
		w.Write([]byte("Server is active, send it somewhere else"))
	}

}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("OK")
}

func readyHandler(w http.ResponseWriter, _ *http.Request) {
	if available {
		w.WriteHeader(200)
		w.Write([]byte("OK")
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
