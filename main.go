package main

import (
    "context"
    "errors"
    "fmt"
    "github.com/cheggaaa/pb/v3"
    "github.com/emirpasic/gods/sets/treeset"
    "github.com/google/go-github/github"
    "github.com/urfave/cli"
    "golang.org/x/oauth2"
    "log"
    "os"
    "reflect"
    "strconv"
    "text/tabwriter"
    "time"
)

type Repo struct {
    name  string
    count int
}

var app = cli.NewApp()
var ghc = getGithubClient()
var ghubT = os.Getenv("GH_TOKEN")
var w = new(tabwriter.Writer)

func getGithubClient() *github.Client {
    ctx := context.Background()
    ts := oauth2.StaticTokenSource(
        &oauth2.Token{AccessToken: ghubT},
    )
    tc := oauth2.NewClient(ctx, ts)
    return github.NewClient(tc)
}

func info() {
    app.Name = "Bruno"
    app.Author = "David Davtian"
    app.Version = "1.0.0"
    app.Compiled = time.Now()
    app.Usage = "GitHub N Repos Demo"
}

func checkArgs(c *cli.Context) (string, int, error) {
    org := c.Args().Get(0)
    s := reflect.TypeOf(org).Kind()
    if s != reflect.String {
        return "", 0, errors.New("organization must be a string")
    }

    mr, err := strconv.Atoi(c.Args().Get(1))
    if err != nil {
        return "", 0, err
    }
    return org, mr, nil
}

func commands() {
    app.Commands = [] cli.Command{
        {
            Name:        "stars",
            Aliases:     []string{"s"},
            Usage:       "Top-N repos by number of stars from a GitHub Organization",
            UsageText:   "bruno stars <org> <max-results> [example: bruno stars Netflix 20]",
            Description: "Top-N repos by number of stars from a GitHub Organization",
            ArgsUsage:   "bruno stars Netflix 20",
            Action: func(c *cli.Context) error {
                // Top-N repos by number of stars //
                fmt.Fprintf(w, "\n\n%s\n\n", "--- Top N by Stars ---")
                org, mr, err := checkArgs(c)
                if err != nil {
                    return errors.New("bruno stars <org> <max-results> [example: bruno stars Netflix 20]")
                }
                ShowTopNReposBySortType(*ghc, org, "stars", mr)
                return nil
            },
        },
        {
            Name:        "forks",
            Aliases:     []string{"f"},
            Usage:       "Top-N repos by number of forks from a GitHub Organization",
            UsageText:   "bruno forks <org> <max-results> [example: bruno forks Netflix 20]",
            Description: "Top-N repos by number of forks from a GitHub Organization",
            ArgsUsage:   "bruno forks Netflix 20",
            Action: func(c *cli.Context) error {
                // Top-N repos by number of forks //
                fmt.Fprintf(w, "\n\n%s\n\n", "--- Top N by Forks ---")
                org, mr, err := checkArgs(c)
                if err != nil {
                    return errors.New("bruno forks <org> <max-results> [example: bruno forks Netflix 20]")
                }
                ShowTopNReposBySortType(*ghc, org, "forks", mr)
                return nil
            },
        },
        {
            Name:        "pulls",
            Aliases:     []string{"p"},
            Usage:       "Top-N repos by number of pulls from a GitHub Organization",
            UsageText:   "bruno pulls <org> <max-results> [example: bruno pulls Netflix 20]",
            Description: "Top-N repos by number of pulls from a GitHub Organization",
            ArgsUsage:   "bruno pulls Netflix 20",
            Action: func(c *cli.Context) error {
                // Top-N repos by number of Pull Requests (PRs) //
                fmt.Fprintf(w, "\n\n%s\n\n", "--- Top N by Pulls ---")
                org, mr, err := checkArgs(c)
                if err != nil {
                    return errors.New("bruno pulls <org> <max-results> [example: bruno pulls Netflix 20]")
                }
                ShowTopNReposByPulls(*ghc, org, mr)
                return nil
            },
        },
        {
            Name:        "contributions",
            Aliases:     []string{"c"},
            Usage:       "Top-N repos by number of contributions % from a GitHub Organization",
            UsageText:   "bruno contributions <org> <max-results> [example: bruno contributions Netflix 20]",
            Description: "Top-N repos by number of contributions from a GitHub Organization",
            ArgsUsage:   "bruno contributions Netflix 20",
            Action: func(c *cli.Context) error {
                // Top-N repos by contribution percentage (PRs/Forks) //
                fmt.Fprintf(w, "\n\n%s\n\n", "--- Top N by Contributions %% ---")
                org, mr, err := checkArgs(c)
                if err != nil {
                    return errors.New("bruno contributions <org> <max-results> [example: bruno contributions Netflix 20]")
                }
                ShowTopNReposByContributionPercentage(*ghc, org, mr)
                return nil
            },
        },
    }
}

func main() {

    if ghubT == "" {
        fmt.Print("\n\nPlease set your GH_TOKEN environment variable to point to a GitHub Token\n\n")
        os.Exit(1)
    }

    w.Init(os.Stdout, 5, 8, 10, '\t', 0)
    defer w.Flush()

    info()
    commands()

    err := app.Run(os.Args)
    if err != nil {
        log.Fatal(err)
    }
}

// custom comparator (sort by count) on Repo
func byCount(a, b interface{}) int {
    // Type assertion will panic if this is not respected
    r1 := a.(Repo)
    r2 := b.(Repo)

    // negative , if a < b
    // zero     , if a == b
    // positive , if a > b

    switch {
    case r1.count < r2.count:
        return -1
    case r1.count > r2.count:
        return 1
    default:
        return 0
    }
}

func ShowTopNReposByContributionPercentage(ghc github.Client, organization string, maxResults int) {
    repos, err := FetchRepositories(ghc, organization, 10000)
    if err != nil {
        fmt.Println(err)
    }

    defer printFooter()

    var counter int

    inverseByCountComparator := func(a, b interface{}) int {
        return -byCount(a, b)
    }

    treeSet := treeset.NewWith(inverseByCountComparator)

    fmt.Print("\n!! -- Fetching Pull Requests from Each Repo, Costly Operation, Multiple Repos with 0 PRs Ignored -- !!\n\n")

    bar := pb.StartNew(len(repos))

    for _, r := range repos {
        name := r.Name
        pullRequests, _ := FetchPullRequests(ghc, organization, *name, 10000)
        treeSet.Add(Repo{name: *name, count: len(pullRequests)})
        bar.Increment()
    }

    bar.Finish()

    reps, err := SearchRepositories(ghc, organization, "forks", 10000)

    printHeader()

    it := treeSet.Iterator()
    for it.Next() {
        if counter >= maxResults {
            return
        }
        counter++
        v := it.Value()
        for _, searchResults := range reps {
            for _, repo := range searchResults.Repositories {
                if *repo.Name == v.(Repo).name {
                    pullCount := v.(Repo).count
                    forkCount := *repo.ForksCount
                    percent := float64(pullCount) / float64(forkCount) * 100
                    if forkCount > 0 {
                        fmt.Fprintf(w, "\n%d.\t%s\t%.2f %%\t", counter, v.(Repo).name, percent)
                    }
                }
            }
        }
    }
}

func ShowTopNReposBySortType(ghc github.Client, organization string, sortType string, maxResults int) {
    repos, err := SearchRepositories(ghc, organization, sortType, 10000)
    if err != nil {
        fmt.Println(err)
    }

    defer printFooter()

    var counter int

    printHeader()

    for _, searchResults := range repos {
        for i, repo := range searchResults.Repositories {
            if i >= maxResults {
                return
            }
            counter++
            switch sortType {
            case "stars":
                fmt.Fprintf(w, "\n%d.\t%s\t%d\t", counter, *repo.Name, *repo.StargazersCount)
            case "forks":
                fmt.Fprintf(w, "\n%d.\t%s\t%d\t", counter, *repo.Name, *repo.ForksCount)
            }
        }
    }
}

func ShowTopNReposByPulls(ghc github.Client, organization string, maxResults int) {
    repos, err := FetchRepositories(ghc, organization, 10000)
    if err != nil {
        fmt.Println(err)
    }

    defer printFooter()

    var counter int

    inverseByCountComparator := func(a, b interface{}) int {
        return -byCount(a, b)
    }

    treeSet := treeset.NewWith(inverseByCountComparator)

    fmt.Print("\n!! -- Fetching Pull Requests from Each Repo, Costly Operation, Multiple Repos with 0 PRs Ignored -- !!\n\n")

    bar := pb.StartNew(len(repos))

    for _, repo := range repos {
        name := repo.Name
        pullRequests, _ := FetchPullRequests(ghc, organization, *name, 10000)
        treeSet.Add(Repo{name: *name, count: len(pullRequests)})
        bar.Increment()
    }
    bar.Finish()

    printHeader()
    it := treeSet.Iterator()
    for it.Next() {
        if counter >= maxResults {
            return
        }
        counter++
        v := it.Value()
        fmt.Fprintf(w, "\n%d.\t%s\t%d\t", counter, v.(Repo).name, v.(Repo).count)
    }
}

func SearchRepositories(ghc github.Client, organization string, sortType string, perPage int) ([]*github.RepositoriesSearchResult, error) {
    options := &github.SearchOptions{Sort: sortType, Order: "desc", ListOptions: github.ListOptions{PerPage: perPage}}

    var allSearchResults []*github.RepositoriesSearchResult

    for {
        searchResults, response, err := ghc.Search.Repositories(context.Background(), fmt.Sprintf("org:%s", organization), options)
        if err != nil {
            return []*github.RepositoriesSearchResult{}, err
        }

        allSearchResults = append(allSearchResults, searchResults)
        if response.NextPage == 0 {
            break
        }
        options.Page = response.NextPage
    }
    return allSearchResults, nil
}

func FetchRepositories(ghc github.Client, organization string, perPage int) ([]*github.Repository, error) {
    options := &github.RepositoryListByOrgOptions{Type: "all", ListOptions: github.ListOptions{PerPage: perPage}}

    var allRepositories []*github.Repository

    for {
        repositories, response, err := ghc.Repositories.ListByOrg(context.Background(), organization, options)
        if err != nil {
            return []*github.Repository{}, err
        }
        allRepositories = append(allRepositories, repositories...)
        if response.NextPage == 0 {
            break
        }
        options.Page = response.NextPage
    }
    return allRepositories, nil
}

func FetchPullRequests(ghc github.Client, owner string, repo string, perPage int) ([]*github.PullRequest, error) {
    options := &github.PullRequestListOptions{State: "open", ListOptions: github.ListOptions{PerPage: perPage}}

    var allPullRequests []*github.PullRequest

    for {
        pullRequests, response, err := ghc.PullRequests.List(context.Background(), owner, repo, options)
        if err != nil {
            return []*github.PullRequest{}, nil
        }
        allPullRequests = append(allPullRequests, pullRequests...)

        if response.NextPage == 0 {
            break
        }
        options.Page = response.NextPage
    }
    return allPullRequests, nil
}

func printHeader() {
    fmt.Fprintf(w, "\n%s\t%s\t%s\t", "Num", "Repo", "Count")
    fmt.Fprintf(w, "\n%s\t%s\t%s\t", "----", "-----", "-----")
}

func printFooter() {
    fmt.Fprintf(w, "\n%s\t%s\t%s\n\n", "----", "-----", "-----")
}
