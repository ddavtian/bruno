# bruno

Bruno is a simple demo CLI tool that interacts with a GitHub Organization and its respective repositories and 
performs the following actions:

* Top-N repos by number of stars.
* Top-N repos by number of forks.
* Top-N repos by number of Pull Requests (PRs).
* Top-N repos by contribution percentage (PRs/Forks).

### Install

`go get github.com/ddavtian/bruno`

### Token

Since Bruno interacts with GitHub's API it requires a GitHub Access Token. GitHub allows anonymous access but 
throttles and limits quite quickly.

`export GH_TOKEN=<YOUR GITHUB PERSONAL TOKEN>`

### Options

```
bruno -h

NAME:
   Bruno - GitHub N Repos Demo

USAGE:
   bruno [global options] command [command options] [arguments...]

VERSION:
   1.0.0

AUTHOR:
   David Davtian

COMMANDS:
   stars, s          Top-N repos by number of stars from a GitHub Organization
   forks, f          Top-N repos by number of forks from a GitHub Organization
   pulls, p          Top-N repos by number of pulls from a GitHub Organization
   contributions, c  Top-N repos by number of contributions % from a GitHub Organization
   help, h           Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

### Examples

#### Top 10 Netflix Repositories Ranked by Stars

```
bruno stars Netflix 10

--- Top N by Stars ---

Num		Repo				Count
----		-----				-----
1.		Hystrix				18424
2.		falcor				9471
3.		eureka				8471
4.		zuul				8277
5.		pollyjs				8101
6.		chaosmonkey			7385
7.		SimianArmy			7357
8.		fast_jsonapi			4611
9.		security_monkey			3843
10.		vizceral			3280
----		-----				-----
```
#### Top 10 Netflix Repositories Ranked by Forks

```
bruno forks Netflix 10

--- Top N by Forks ---

Num		Repo				Count
----		-----				-----
1.		Hystrix				3799
2.		eureka				2374
3.		zuul				1643
4.		SimianArmy			1069
5.		ribbon				786
6.		security_monkey			736
7.		conductor			663
8.		chaosmonkey			548
9.		Cloud-Prize			483
10.		falcor				454
----		-----				-----
```
#### Top 10 Netflix Repositories Ranked by Pull Requests
```
bruno pulls Netflix 10

--- Top N by Pulls ---

!! -- Fetching Pull Requests from Each Repo, Costly Operation, Multiple Repos with 0 PRs Ignored -- !!

171 / 171 [----------------------------------------------------------------------------] 100.00% 4 p/s

Num		Repo			Count
----		-----			-----
1.		astyanax		43
2.		Hystrix			42
3.		archaius		31
4.		ribbon			30
5.		conductor		21
6.		fast_jsonapi		16
7.		zuul			14
8.		dynomite		12
9.		hollow			11
10.		governator		8
----		-----			-----
```
#### Top 10 Netflix Repositories Ranked by Contribution Percentage (PRs/Forks)

```
bruno contributions Netflix 10

--- Top N by Contributions %% ---

!! -- Fetching Pull Requests from Each Repo, Costly Operation, Multiple Repos with 0 PRs Ignored -- !!

171 / 171 [----------------------------------------------------------------------------] 100.00% 4 p/s

Num		Repo			Count
----		-----			-----
1.		astyanax		11.56 %
2.		Hystrix			1.11 %
3.		archaius		7.13 %
4.		ribbon			3.82 %
5.		conductor		3.17 %
6.		fast_jsonapi		4.48 %
7.		zuul			0.85 %
8.		dynomite		2.89 %
9.		hollow			8.80 %
10.		governator		4.91 %
----		-----			-----
```
### Manual Tests

* Test 1: Test to ensure an error is displayed if `GH_TOKEN` is not exported as an environment variable.
* Test 2: Test usage to ensure the user is passing in proper set of arguments for each of the commands.
* Test 3: Test with an organization that has a small set of repos.
* Test 4: Test with an organization that has a large set of repos.
* Test 5: Test to ensure stdout formatting is proper.
* Test 6: Test each of the supported commands, stars, forks, pulls and contributions.
* Test 7: Test each of the supported short hand commands, s, f, p and c.
