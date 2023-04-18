package DataExtractor01

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"regexp"

	"github.com/Hamzaelkhatri/ImageBuilder/v2"
	"github.com/pemora/api01"
)

var PiscineQuery = `
query Piscine($userId:Int!,$eventIds:[Int!]!){
	score:transaction(
	  where: {userId: {_eq: $userId}, eventId: {_in: $eventIds}, type: {_eq: "xp"},object:{parents:{parent:{name:{_neq:"Piscine Go"}}}}}
	) {
	  amount
	  isBonus
	  path
	  userLogin
	  object{
		name
		parents{
			paths{
				path
			}
		  difficulty:attrs(path:"difficulty")
		  parent{
			questName:name
			parents{
			  eventName:parent{
				name
			  }
			}
		  }
		}
	  }
	}
  }
`

var RaidQuery = `
query Raids($userid:Int!){
	event_user(where:{user:{id:{_eq:$userid}},event:{object:{name:{_in:["quad","sudoku","quadchecker"]}}}}){
		level
		xp{
			amount
		}
		event{
		status
		progresses(where:{user:{id:{_eq:$userid}}}){
			grade
		}
		groups_aggregate{
			nodes{
				members{
				 userLogin
				}
			}
		}
		  path
		  object{
			name
		  }
		}
	}
}
`

type Piscine struct {
	Amount    int
	IsBonus   bool
	Path      string
	UserLogin string
	Object    struct {
		Name    string
		Parents []struct {
			Paths []struct {
				Path string
			}
			Difficulty float32
			Parent     struct {
				QuestName string
				Parents   []struct {
					EventName struct {
						Name string
					}
				}
			}
		}
	}
}

type Raids []struct {
	Level int `default:"0"`
	Xp    struct {
		Amount int `default:"0"`
	}
	Event struct {
		Status     string `default:"incomplete"`
		Progresses []struct {
			Grade float32 `default:"0"`
		}
		GroupsAggregate struct {
			Nodes []struct {
				Members []struct {
					UserLogin string
				}
			}
		} `json:"groups_aggregate"`
		Path   string `default:"raid"`
		Object struct {
			Name string `default:"raid"`
		}
	}
}

func Init(client api01.Client, idUser int) ([]Piscine, Raids) {
	var quest []Piscine
	var raids Raids
	resp := client.GraphqlQuery(PiscineQuery, api01.Vars{"userId": idUser, "eventIds": []int{3}})

	if resp.HasErrors() {
		log.Fatal(resp.Errors)
	}

	if resp.Data["score"] == nil {
		log.Fatal("No data returned")
	}

	body, err := json.Marshal(resp.Data["score"].([]interface{}))
	if err != nil {
		log.Fatal(err)
	}
	// log.Println(quest)
	err = json.Unmarshal(body, &quest)
	if err != nil {
		log.Fatal(err)
	}

	resp = client.GraphqlQuery(RaidQuery, api01.Vars{"userid": idUser})
	if resp.HasErrors() {
		log.Fatal(resp.Errors)
	}

	if resp.Data["event_user"] == nil {
		log.Fatal("No data returned")
	}

	body, err = json.Marshal(resp.Data["event_user"].([]interface{}))
	if err != nil {
		log.Fatal(err)
	}
	// log.Println(quest)
	err = json.Unmarshal(body, &raids)
	if err != nil {
		log.Fatal(err)
	}
	return quest, raids
}

func getAvatar(username string) string {
	client := &http.Client{}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://%s/git/api/v1/users/%s", os.Getenv("ENDPOINT"), username), nil)
	if err != nil {
		log.Println(err)
		return ""
	}
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return data["avatar_url"].(string)
}

func ExtractData(idUser int) (string, error) {
	client, err := api01.NewClient(os.Getenv("ENDPOINT"))
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	quest, raids := Init(client, idUser)
	// group by quest
	type Quest struct {
		Name  string
		Xp    int
		count int
	}

	questMap := make(map[string]Quest)
	for _, q := range quest {
		for _, p := range q.Object.Parents {
			for _, path := range p.Paths {
				if path.Path == q.Path {
					questMap[p.Parent.QuestName] = Quest{Name: p.Parent.QuestName, Xp: questMap[p.Parent.QuestName].Xp + q.Amount, count: questMap[p.Parent.QuestName].count + 1}
				}
			}
		}
	}
	for _, r := range raids {
		questMap[r.Event.Object.Name] = Quest{Name: r.Event.Object.Name, Xp: questMap[r.Event.Object.Name].Xp + r.Xp.Amount, count: questMap[r.Event.Object.Name].count + 1}
	}

	sumQuestExercises := 0
	sumCheckpointExercises := 0
	regex := regexp.MustCompile(`^Quest.*`)
	for _, q := range questMap {
		if regex.MatchString(q.Name) {
			sumQuestExercises += q.count
		} else {
			sumCheckpointExercises += q.count
		}
	}

	skill := getSkills()
	type Skill struct {
		skill    float32
		MaxLevel int
		current  int
	}

	skillMap := make(map[string]Skill)
	for _, q := range questMap {
		for _, s := range skill {
			if s.Name == q.Name {
				for _, d := range s.Skill {
					skillMap[d] = Skill{skill: skillMap[d].skill + float32(float32(q.count)/float32(s.MaxLevel)), MaxLevel: s.MaxLevel, current: q.count}
				}
			}
		}
	}
	for k, v := range skillMap {
		skillMap[k] = Skill{skill: v.skill / float32(len(questMap)), MaxLevel: v.MaxLevel, current: v.current}
	}
	xps := 0
	for _, q := range quest {
		xps += q.Amount
	}
	if len(quest) == 0 {
		return "", errors.New("error no data")
	}
	// log.Println(raids)
	raidsCounts := 0
	raidGenrated := func() []ImageBuilder.Raid {
		var r []ImageBuilder.Raid
		for _, raid := range raids {
			r = append(r, ImageBuilder.Raid{
				Name:   raid.Event.Object.Name,
				Grade:  raid.Event.Progresses[0].Grade,
				Status: raid.Event.Status,
			})
			if raid.Event.Progresses[0].Grade >= 1 {
				raidsCounts++
			}
		}
		return r
	}
	/*
			Name   string
		Status string
		Grade  float32
	*/
	// var r Raids = nil
	return ImageBuilder.Init(
		ImageBuilder.CardData{
			Name:              quest[0].UserLogin,
			Avatar:            getAvatar(quest[0].UserLogin),
			Level:             int(getLevel(float64(xps))),
			NumberOfExercises: sumQuestExercises,
			Raids:             raidGenrated(),
			Skills: [][]float32{
				{
					float32(sumQuestExercises),
					float32(sumCheckpointExercises), float32(raidsCounts),
				},
			},
		},
	), nil
}

func getLevel(xp float64) float64 {
	squareRoot := math.Sqrt(math.Pow(-(9*xp)/11-778042/1331, 2) + 11698628938101/28344976)
	cubicRoot := math.Pow(-(9*xp)/22+squareRoot/2-389021/1331, 1.0/3.0)
	return -cubicRoot/3 - 83.0/66 + 7567/(484*cubicRoot)
}
