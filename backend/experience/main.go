package main

import (
	"log"
	"net/http"

	"github.com/daniellawrence/cv/backend/common"
	experiencev1 "github.com/daniellawrence/cv/gen/go/experience/v1"

	"google.golang.org/protobuf/encoding/protojson"
)

var experiences = []*experiencev1.Experience{
	{
		Id:        "block-sre",
		Company:   "Block",
		Title:     "Site Reliability Engineer, Observability",
		StartDate: "2023-05-01",
		EndDate:   "",
		Location:  "Melbourne, Australia",
		Highlights: []string{
			"Led foundational development of Block's observability and telemetry platform",
			"Architected and executed Datadog migrations from Sumologic and NewRelic",
			"Managed full lifecycle of logging, metrics, and tracing for Afterpay, Cash, and Square",
			"Kubernetes fleet management and telemetry pipeline design",
		},
		Skills: []string{
			"AWS", "Datadog", "Vector", "OpenTelemetry", "Kubernetes", "Go",
		},
	},
	{
		Id:        "okta-principal-sre",
		Company:   "Okta",
		Title:     "Principal Site Reliability Engineer, Observability",
		StartDate: "2022-04-01",
		EndDate:   "2023-05-01",
		Location:  "Australia",
		Highlights: []string{
			"Designed automated alert routing solution for on-call response optimization",
			"Managed Wavefront configuration and operational administration",
			"Provided production support and 24/7 on-call responsibilities",
		},
		Skills: []string{
			"AWS", "Go", "Wavefront", "Prometheus", "AlertManager",
		},
	},
	{
		Id:        "paloalto-techlead",
		Company:   "Palo Alto Networks",
		Title:     "Technical Lead, SCM, CI/CD & Observability",
		StartDate: "2018-06-01",
		EndDate:   "2022-04-01",
		Location:  "Melbourne, Australia",
		Highlights: []string{
			"Led Tools, Monitoring, Metrics, and SCM/CI/CD SRE teams scaling from 3 to 12 engineers",
			"Architected and maintained production infrastructure at massive scale",
			"Platform scaling across configuration management, metrics, logs, CI/CD, and Kubernetes",
		},
		Skills: []string{
			"Saltstack", "Prometheus", "ELK", "GitLab", "Terraform",
			"GCP", "Kubernetes", "Python", "PostgreSQL", "Go",
		},
	},
	{
		Id:        "linkedin-staff-sre",
		Company:   "LinkedIn",
		Title:     "Staff Site Reliability Engineer, Jobs & Recruiter",
		StartDate: "2016-02-01",
		EndDate:   "2017-11-01",
		Location:  "San Francisco Bay Area",
		Highlights: []string{
			"SRE for 100+ applications across thousands of hosts in a service-oriented architecture",
			"Mentored SRE team growth from 3 to 12 engineers",
			"Overhauled monitoring philosophy and implementation across multiple teams",
		},
		Skills: []string{
			"Python", "Java", "Go", "Kafka",
		},
	},
	{
		Id:        "medibank-devops",
		Company:   "Medibank Health Solutions",
		Title:     "Linux & Cloud Services Engineer (DevOps)",
		StartDate: "2014-10-01",
		EndDate:   "2015-02-01",
		Location:  "Melbourne, Australia",
		Highlights: []string{
			"Founding member of Cloud Services team scaling from 2 to 5 engineers",
			"Built Medibank's first AWS environment and PCI compliance automation",
		},
		Skills: []string{
			"AWS", "Puppet", "Python",
		},
	},
	{
		Id:        "toll-sysadmin",
		Company:   "Toll Holdings Limited",
		Title:     "Senior Linux & Unix Admin",
		StartDate: "2012-03-01",
		EndDate:   "2013-03-01",
		Location:  "Australia",
	},
	{
		Id:        "ibm-unix-admin",
		Company:   "IBM",
		Title:     "Unix Admin",
		StartDate: "2006-12-01",
		EndDate:   "2008-12-01",
		Location:  "Australia",
	},
}

func listexperience(w http.ResponseWriter, r *http.Request) {
	m := protojson.MarshalOptions{
		UseProtoNames:   false,
		EmitUnpopulated: true,
	}

	w.Header().Set("Content-Type", "application/json")

	data, err := m.Marshal(&experiencev1.ListExperienceResponse{
		Experience: experiences,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(data)
}

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/experience", listexperience)

	log.Println("experience service listening on :8080")

	err := http.ListenAndServe(":8080", common.CorsMiddleware(mux))
	if err != nil {
		log.Fatal(err)
	}
}
