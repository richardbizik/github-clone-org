package cmd

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/spf13/cobra"
)

var organization string

type repoPage struct {
	repositories []string
	currentPage  int
	maxPage      int
}

func init() {
	cloneCmd.PersistentFlags().StringVarP(&organization, "organization", "o", "", "organization directory eg. github.com/richardbizik/ (required)")
	_ = cloneCmd.MarkPersistentFlagRequired("organization")

}

var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "clone will clone the repositories from the entire organization",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cloneGithubOrg(organization)
	},
}

func cloneGithubOrg(ghOrg string) error {
	if !strings.HasPrefix(ghOrg, "github.com") {
		return fmt.Errorf("organization string has to start with github.com")
	}
	pages := make([]repoPage, 0)
	page, err := getRepoPage(ghOrg, 1)
	if err != nil {
		return err
	}
	pages = append(pages, page)
	maxPage := page.maxPage
	totalRepos := len(page.repositories)
	if page.currentPage != maxPage {
		for i := 2; i <= maxPage; i++ {
			page, err = getRepoPage(ghOrg, i)
			if err != nil {
				return err
			}
			pages = append(pages, page)
			totalRepos += len(page.repositories)
		}
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	fmt.Printf("Cloning %d repositories into: %s\n", totalRepos, wd)
	for _, rp := range pages {
		for _, v := range rp.repositories {
			path := strings.TrimPrefix(v, "/")
			gitClone := exec.Command("git", "clone", fmt.Sprintf("git@github.com:%s.git", path))
			var stdBuffer bytes.Buffer
			mw := io.MultiWriter(os.Stdout, &stdBuffer)
			gitClone.Stdout = mw
			gitClone.Stderr = mw
			err := gitClone.Run()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func getRepoPage(ghOrg string, page int) (repoPage, error) {
	c := http.Client{}
	resp := repoPage{}
	parts := strings.SplitN(ghOrg, "/", 2)
	if len(parts) < 2 {
		return resp, fmt.Errorf("expected organization after github.com")
	}
	response, err := c.Get(fmt.Sprintf("https://github.com/orgs/%s/repositories?type=all&page=%d", parts[1], page))
	if err != nil {
		return resp, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return resp, fmt.Errorf("received code %d from github.com", response.StatusCode)
	}
	node, err := htmlquery.Parse(response.Body)
	if err != nil {
		return resp, err
	}
	repositories := htmlquery.FindOne(node, "//div[@id=\"org-repositories\"]/div/div/div[@class=\"Box\"]/ul")
	pagination := htmlquery.FindOne(node, "//div[@id=\"org-repositories\"]/div/div/div[2]/div/em")

	if pagination == nil || pagination.LastChild == nil {
		return resp, fmt.Errorf("could not find a pagination on github org page")
	}
	currentPage := pagination.LastChild.Data
	resp.currentPage, err = strconv.Atoi(currentPage)
	if err != nil {
		return resp, fmt.Errorf("failed to convert page to integer")
	}
	for _, a := range pagination.Attr {
		if a.Key == "data-total-pages" {
			maxPage, err := strconv.Atoi(a.Val)
			if err != nil {
				return resp, fmt.Errorf("failed to convert max page to integer")
			}
			resp.maxPage = maxPage
			break
		}

	}
	if repositories == nil || repositories.LastChild == nil || repositories.LastChild.PrevSibling == nil {
		return resp, fmt.Errorf("could not find a repositories on github org page")
	}
	if repositories.LastChild.PrevSibling == nil {
		return resp, fmt.Errorf("could not find a repositories on github org page")
	}
	lis := htmlquery.Find(repositories, "//li")
	for _, n := range lis {
		ahref := htmlquery.FindOne(n, "//div/div/div/h3/a")
		for _, a := range ahref.Attr {
			if a.Key == "href" {
				resp.repositories = append(resp.repositories, a.Val)
				break
			}
		}
	}

	return resp, nil
}
