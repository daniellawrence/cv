package main

import educationv1 "github.com/daniellawrence/cv/gen/go/education/v1"

var educations = []*educationv1.Education{
	{
		Id:          1,
		Institution: "Griffith University",
		Degree:      "Bachelor, Information Technology",
		StartDate:   "2005-12",
		EndDate:     "2008-06",
	},
}
