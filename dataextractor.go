package dataextractor01

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

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

type Piscine struct {
	Amount  int
	IsBonus bool
	Path    string
	Object  struct {
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

func ExtractData() {
	log.Println("Extracting data...")
	log.Println("Endpoint: ", os.Getenv("ENDPOINT"))
	client, err := api01.NewClient(os.Getenv("ENDPOINT"))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Client created")
	resp := client.GraphqlQuery(PiscineQuery, api01.Vars{"userId": 792, "eventIds": []int{3}})

	if resp.HasErrors() {
		log.Fatal(resp.Errors)
	}

	log.Println("Query executed")

	if resp.Data["score"] == nil {
		log.Fatal("No data returned")
	}

	var quest []Piscine
	body, err := json.Marshal(resp.Data["score"].([]interface{}))
	if err != nil {
		log.Fatal(err)
	}
	// log.Println(quest)
	err = json.Unmarshal(body, &quest)
	if err != nil {
		log.Fatal(err)
	}

	// sum of all amounts
	var sum int = 0
	for _, q := range quest {
		sum += q.Amount
	}
	log.Println("Total XP: ", sum)
	// log.Println("Level: ", getLevelFromXp(sum, 0))

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

	// // group by difficulty
	type Difficulty struct {
		Name  string
		diff  float32
		count int
	}
	difficultyMap := make(map[string]Difficulty)
	for _, q := range quest {
		// take the name the quest if the path equals the parent path
		for _, p := range q.Object.Parents {
			for _, path := range p.Paths {
				if path.Path == q.Path {
					if difficultyMap[p.Parent.QuestName].diff == 0 {
						difficultyMap[p.Parent.QuestName] = Difficulty{Name: p.Parent.QuestName, diff: p.Difficulty, count: 1}
					} else {
						difficultyMap[p.Parent.QuestName] = Difficulty{Name: p.Parent.QuestName, diff: ((difficultyMap[p.Parent.QuestName].diff) / float32(difficultyMap[p.Parent.QuestName].count)) + p.Difficulty, count: difficultyMap[p.Parent.QuestName].count + 1}
					}
				}
			}
		}
	}

	skill := getSkills()
	type Skill struct {
		skill    float32
		MaxLevel int
		current  int
	}

	// skillMap := Skill{Alogrithms: 0, Math: 0, Unix: 0, Golang: 0, ProblemSolving: 0}
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

	// add percentage
	for k, v := range skillMap {
		skillMap[k] = Skill{skill: v.skill / float32(len(questMap)), MaxLevel: v.MaxLevel, current: v.current}
	}

	for k, v := range skillMap {
		fmt.Printf("%s : %f \n", k, v.skill*100)
	}
}
