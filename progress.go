package dataextractor01

import "encoding/json"

type Skills struct {
	MaxLevel int
	Skill    []string
	Name     string
}

var ProgressData string = `
[
	{
		"maxLevel": 13,
		"skill": [
			"Unix",
			"Git"
		],
		"name": "Quest 01"
	},
	{
		"maxLevel": 8,
		"skill": [
			"Algorithms",
			"Golang",
			"Git"
		],
		"name": "Quest 02"
	},
	{
		"maxLevel": 12,
		"skill": [
			"Algorithms",
			"Golang"
		],
		"name": "Quest 03"
	},
	{
		"maxLevel": 9,
		"skill": [
			"Algorithms",
			"Golang",
			"Math"
		],
		"name": "Quest 04"
	},
	{
		"maxLevel": 18,
		"skill": [
			"Algorithms",
			"Golang"
		],
		"name": "Quest 05"
	},
	{
		"maxLevel": 4,
		"skill": [
			"Unix",
			"Golang"
		],
		"name": "Quest 06"
	},
	{
		"maxLevel": 7,
		"skill": [
			"Algorithms",
			"Golang"
		],
		"name": "Quest 07"
	},
	{
		"maxLevel": 5,
		"skill": [
			"Unix",
			"Golang"
		],
		"name": "Quest 08"
	},
	{
		"maxLevel": 8,
		"skill": [
			"Golang"
		],
		"name": "Quest 09"
	},
	{
		"maxLevel": 12,
		"skill": [
			"Problem Solving",
			"Golang"
		],
		"name": "Hackathon"
	},
	{
		"maxLevel": 15,
		"skill": [
			"Algorithms",
			"Golang"
		],
		"name": "Quest 11"
	},
	{
		"maxLevel": 11,
		"skill": [
			"Algorithms",
			"Golang"
		],
		"name": "Quest 12"
	},
	{
		"maxLevel": 7,
		"skill": [
			"Algorithms",
			"Golang"
		],
		"name": "Checkpoint 01"
	},
	{
		"maxLevel": 7,
		"skill": [
			"Algorithms",
			"Golang"
		],
		"name": "Checkpoint 02"
	},
	{
		"maxLevel": 8,
		"skill": [
			"Algorithms",
			"Golang",
			"Math"
		],
		"name": "Checkpoint 03"
	},
	{
		"maxLevel": 9,
		"skill": [
			"Algorithms",
			"Golang",
			"Math"
		],
		"name": "Checkpoint 04"
	}
]`

func getSkills() []Skills {
	var skills []Skills
	json.Unmarshal([]byte(ProgressData), &skills)
	return skills
}
