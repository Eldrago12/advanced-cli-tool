package github

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/google/go-github/v39/github"
	"golang.org/x/oauth2"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var credentials Credentials

func authenticate() *github.Client {
	loadCredentials()
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: basicAuthToken(credentials.Username, credentials.Password),
		},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

func basicAuthToken(username, password string) string {
	return base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
}

func loadCredentials() {
	configDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Println("Error getting config directory:", err)
		return
	}
	configFile := filepath.Join(configDir, "act-cli", "credentials.json")

	file, err := os.Open(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			promptForCredentials()
			saveCredentials(configFile)
		} else {
			fmt.Println("Error opening credentials file:", err)
		}
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&credentials); err != nil {
		fmt.Println("Error decoding credentials:", err)
	}
}

func saveCredentials(configFile string) {
	configDir := filepath.Dir(configFile)
	if err := os.MkdirAll(configDir, os.ModePerm); err != nil {
		fmt.Println("Error creating config directory:", err)
		return
	}

	file, err := os.Create(configFile)
	if err != nil {
		fmt.Println("Error creating credentials file:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(&credentials); err != nil {
		fmt.Println("Error encoding credentials:", err)
	}
}

func promptForCredentials() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter GitHub username: ")
	username, _ := reader.ReadString('\n')
	fmt.Print("Enter GitHub password: ")
	password, _ := reader.ReadString('\n')

	credentials = Credentials{
		Username: strings.TrimSpace(username),
		Password: strings.TrimSpace(password),
	}
}

func Logout() {
	configDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Println("Error getting config directory:", err)
		return
	}
	configFile := filepath.Join(configDir, "act-cli", "credentials.json")
	os.Remove(configFile)
	fmt.Println("Logged out successfully.")
}

func ListRepos() {
	client := authenticate()
	ctx := context.Background()
	repos, _, err := client.Repositories.List(ctx, "", nil)
	if err != nil {
		fmt.Println("Error listing repos:", err)
		return
	}

	for _, repo := range repos {
		fmt.Println(*repo.Name)
	}
}

func RunGitCommand(args ...string) {
	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running git %s: %v\n", strings.Join(args, " "), err)
	}
}

func Commit(message string) {
	RunGitCommand("commit", "-m", message)
}

func Push(repo, branch string) {
	if branch == "" {
		branch = "main"
	}

	// Initialize git if necessary
	if _, err := os.Stat(".git"); os.IsNotExist(err) {
		RunGitCommand("init")
		RunGitCommand("remote", "add", "origin", fmt.Sprintf("https://github.com/%s/%s.git", credentials.Username, repo))
	}

	RunGitCommand("add", ".")
	Commit("Auto commit")

	if err := handleDivergentBranch(branch); err != nil {
		fmt.Printf("Error handling divergent branch: %v\n", err)
		return
	}

	RunGitCommand("push", "-u", "origin", branch)
}

func Pull(branch string) {
	RunGitCommand("pull", "--rebase", "origin", branch)
}

func CreateBranch(branch string) {
	RunGitCommand("checkout", "-b", branch)
}

func handleDivergentBranch(branch string) error {
	err := exec.Command("git", "pull", "--rebase", "origin", branch).Run()
	if err != nil {
		return fmt.Errorf("error handling divergent branch: %w", err)
	}
	return nil
}

func PushWithLFS(repo, branch, file string) {
	RunGitCommand("lfs", "track", file)
	RunGitCommand("add", ".gitattributes")
	RunGitCommand("add", file)
	Commit("Auto commit with LFS")

	if err := handleDivergentBranch(branch); err != nil {
		fmt.Printf("Error handling divergent branch: %v\n", err)
		return
	}

	RunGitCommand("push", "-u", "origin", branch)
}
