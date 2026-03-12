package main

import interestv1 "github.com/daniellawrence/cv/gen/go/interest/v1"

var interests = []*interestv1.Interest{
	{
		Id:   "tech",
		Type: "Technical Interests",
		Names: []string{
			"Site Reliability Engineering",
			"Observability Architecture",
			"Continuous Deployment",
			"Developer Tooling",
			"Process Improvement",
		},
	},
	{
		Id:   "skills",
		Type: "Languages / Skills",
		Names: []string{
			"go",
			"python",
			"observability",
			"ai-skills",
			"prometheus-stack",
		},
	},
	{
		Id:   "hobbies",
		Type: "Hobbies",
		Names: []string{
			"3D printing / All things maker",
			"Home Automation",
			"Wildlife Rescue",
		},
	},
}
