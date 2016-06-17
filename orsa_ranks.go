package matcher

import (
    "fmt"
    "strings"
    "regexp"
)

type MatchGroup struct {
    Teams [4]TeamRank
    N   int
    AvgPoints float64
}

func (m MatchGroup) ToStringReddit (groupNumber int, weekNumber int, platform string, format string) string {
    const MuutBaseUrl string = "https://muut.com/arlchampionships#!"
    const MuutPCPS4Format string = "[%s](%s/week-%d-%s-pcps4:%s)"
    const MuutXboxFormat string = "[%s](%s/week-%d-xbox-all-formats:%s)"

    noSpecialChars, _ := regexp.Compile("[^A-Za-z0-9 ]+")

    outputStr := fmt.Sprintf("Group %d: ", groupNumber)
    teamMatches := ""
    for i, team := range m.Teams {
        if team.Name == "" { continue }
        //team names for match header
        outputStr += team.Name
        if i+1 < m.N {
            outputStr += ", "
        }

        //team matches with forum links
        nextTeam := (i+1) % m.N
        teamStr := team.Name + " vs " + m.Teams[nextTeam].Name
        //remove special characters
        teamStr = noSpecialChars.ReplaceAllString(teamStr, "")
        //replace all spaces with hyphens
        teamUrl := strings.Replace(strings.ToLower(teamStr), " ", "-", -1)

        if len(teamUrl) > 27 {
            teamUrl = teamUrl[:27]
        }

        teamUrl = strings.Trim(teamUrl, "-")

        if platform == "pcps4" {
            teamMatches += fmt.Sprintf(MuutPCPS4Format, teamStr, MuutBaseUrl, weekNumber, format, teamUrl)
        } else {
            teamMatches += fmt.Sprintf(MuutXboxFormat, teamStr, MuutBaseUrl, weekNumber, teamUrl)
        }
        teamMatches += "\n\n"
    }
    outputStr += "\n\n" + teamMatches + "\n"


    return outputStr
}

func (m MatchGroup) ToStringORSA (groupNumber int, weekNumber int, platform string, format string) string {
    const MuutBaseUrl string = "https://muut.com/arlchampionships#!"
    const MuutPCPS4Format string = "%s/week-%d-%s-pcps4:%s"
    const MuutXboxFormat string = "%s/week-%d-xbox-all-formats:%s"
    const ORSAFormat string = `<button class="orsa-accordion">%s</button>
<div class="orsa-muut-thread">
<a class="muut" href="%s/comments"></a>
<script src="//cdn.muut.com/1/moot.min.js"></script>
</div>`

    noSpecialChars, _ := regexp.Compile("[^A-Za-z0-9 ]+")

    outputStr := ""
    teamMatches := ""
    for i, team := range m.Teams {
        if team.Name == "" { continue }
        //team matches with forum links
        nextTeam := (i+1) % m.N
        teamStr := team.Name + " vs " + m.Teams[nextTeam].Name
        //remove special characters
        teamStr = noSpecialChars.ReplaceAllString(teamStr, "")
        //replace all spaces with hyphens
        teamUrl := strings.Replace(strings.ToLower(teamStr), " ", "-", -1)

        if len(teamUrl) > 27 {
            teamUrl = teamUrl[:27]
        }

        teamUrl = strings.Trim(teamUrl, "-")

        var muutURL string
        if platform == "pcps4" {
            muutURL = fmt.Sprintf(MuutPCPS4Format, MuutBaseUrl, weekNumber, format, teamUrl)
        } else {
            muutURL = fmt.Sprintf(MuutXboxFormat, MuutBaseUrl, weekNumber, teamUrl)
        }

        teamMatches += fmt.Sprintf(ORSAFormat, teamStr, muutURL)
        teamMatches += "\n"
    }
    outputStr += teamMatches

    return outputStr
}

type TeamRank struct {
    Name    string  `csv:"Team"`
    Played  int     `csv:"GP"`
    Wins    int     `csv:"W"`
    Losses  int     `csv:"L"`
    Draws   int     `csv:"D"`
    Points  int     `csv:"Pts"`
    WinPercentage   string `csv:"Win Pct"`
    GoalDiff    int `csv:"GD"`
    GoalsFor    int `csv:"GF"`
    GoalsAgainst    int     `csv:"GA"`
}

type RankedList []TeamRank

func (r RankedList) Len() int { return len(r) }
func (r RankedList) Less(i, j int) bool {
    if r[i].Points < r[j].Points {
        return false
    } else if r[i].Points > r[j].Points {
        return true
    } else {
        if r[i].GoalDiff < r[j].GoalDiff {
            return false
        }
    }
    return true
}
func (r RankedList) Swap(i, j int) { 
    r[i], r[j] = r[j], r[i]
}
func (r *RankedList) Push(x interface{}) {
    *r = append(*r, x.(TeamRank))
}
func (r *RankedList) Pop() interface{} {
    old := *r
    n := len(old)
    x := old[n-1]
    *r = old[0 : n-1]
    return x
}
