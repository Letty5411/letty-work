package main

import (
	"context"
	"os/exec"
	"time"

	"github.com/google/go-github/github"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

const (
	DefaultRepo  = "pouch"
	DefaultOwner = "alibaba"
)

type Client struct {
	client *github.Client
	cfg    Config
}

// Config refers
type Config struct {
	Owner       string
	Repo        string
	AccessToken string
}

func main() {
	var cfg Config
	var cmdServe = &cobra.Command{
		Use:  "",
		Args: cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			Run(cfg)
		},
	}

	flagSet := cmdServe.Flags()
	flagSet.StringVarP(&cfg.Owner, "owner", "o", DefaultOwner, "github ID to which connect in GitHub")
	flagSet.StringVarP(&cfg.Repo, "repo", "r", DefaultRepo, "github repo to which connect in GitHub")
	flagSet.StringVarP(&cfg.AccessToken, "token", "t", "", "access token to have some control on resources")

	cmdServe.Execute()
}

func Run(cfg Config) error {

	var c Client
	// Create an authenticated Client
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.AccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	c.client = github.NewClient(tc)

	c.cfg = cfg

	//t := time.Date(2018, time.May, 14, 00, 00, 0, 0, time.UTC)
	t := time.Now()
	for {
		//commit := make([]*github.RepositoryCommit, 0)
		commit, _ := c.GetFilterCommit(t)

		logrus.Println(t)
		logrus.Println(len(commit))

		if len(commit) != 0 {
			c.RunCI(commit)
		}

		// Get the current time and check if there is any update

		t = time.Now()
		time.Sleep(600 * time.Second)
	}

	//baseUrl := *pr[0].GetBase().GetRepo().URL

	return nil
}

func (c *Client) RunCI(commit []*github.RepositoryCommit) {
	logrus.Println("In ci")
	for _, v := range commit {
		logrus.Printf("%s", v.GetSHA())

		{
			// CI on AliOS
			cmd := exec.Command("java", "-jar", "/home/sit/letty/src/github.com/letty-work/pouch-ci/jenkins-cli.jar", "-s",
				"http://tester:tester@11.160.112.29:8080/", "build", "-p", "commit="+v.GetSHA(), "OpenSourcePouchOnInternalAlios49")
			logrus.Println(cmd)
			err := cmd.Start()
			if err != nil {
				logrus.Errorf("%s", err)
			}
			logrus.Printf("Waiting for command to finish...")
			err = cmd.Wait()
			logrus.Printf("Command finished with error: %v", err)
		}

		{
			// CI on opensource AliOS
			cmd := exec.Command("java", "-jar", "/home/sit/letty/src/github.com/letty-work/pouch-ci/jenkins-cli.jar", "-s",
				"http://tester:tester@11.160.112.29:8080/", "build", "-p", "commit="+v.GetSHA(), "OpenSourcePouch4.9")
			logrus.Println(cmd)
			err := cmd.Start()
			if err != nil {
				logrus.Errorf("%s", err)
			}
			logrus.Printf("Waiting for command to finish...")
			err = cmd.Wait()
			logrus.Printf("Command finished with error: %v", err)
		}
		{
			// CI on opensource AliOS 3.10
			cmd := exec.Command("java", "-jar", "/home/sit/letty/src/github.com/letty-work/pouch-ci/jenkins-cli.jar", "-s",
				"http://tester:tester@11.160.112.29:8080/", "build", "-p", "commit="+v.GetSHA(), "OpenSourcePouchOnInternalAlios3.10")
			logrus.Println(cmd)
			err := cmd.Start()
			if err != nil {
				logrus.Errorf("%s", err)
			}
			logrus.Printf("Waiting for command to finish...")
			err = cmd.Wait()
			logrus.Printf("Command finished with error: %v", err)
		}
		{
			// CI on opensource performance test
			cmd := exec.Command("java", "-jar", "/home/sit/letty/src/github.com/letty-work/pouch-ci/jenkins-cli.jar", "-s",
				"http://tester:tester@11.160.112.29:8080/", "build", "-p", "commit="+v.GetSHA(), "PerformanceOpensourcePouch")
			logrus.Println(cmd)
			err := cmd.Start()
			if err != nil {
				logrus.Errorf("%s", err)
			}
			logrus.Printf("Waiting for command to finish...")
			err = cmd.Wait()
			logrus.Printf("Command finished with error: %v", err)
		}
	}
}
