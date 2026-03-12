package main

import experiencev1 "github.com/daniellawrence/cv/gen/go/experience/v1"

var experiences = []*experiencev1.Experience{
	{
		Id:        "block-sre",
		Company:   "Block",
		Title:     "Site Reliability Engineer, Observability",
		StartDate: "2023-05",
		EndDate:   "current",
		Location:  "Melbourne, Australia",
		Highlights: []string{
			"Led reliability for Block's telemetry platform",
			"Platform migrations from Sumologic and NewRelic into Datadog",
			"Managed full lifecycle of logging, metrics, and tracing for Afterpay, Cash, and Square",
			"Agent skills developement for including sentry management, k8s rightsizing, and datadog cost control",
			"MCP development for internal log platform",
			"telemetry pipeline design and implementation",
			"Mentoring, technical guidance and project planning",
		},
		Skills: []string{
			"AWS", "Datadog", "Vector", "OpenTelemetry", "Kubernetes", "Go", "Sentry",
		},
	},
	{
		Id:        "okta-principal-sre",
		Company:   "Okta",
		Title:     "Principal Site Reliability Engineer, Observability",
		StartDate: "2022-04",
		EndDate:   "2023-05",
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
		StartDate: "2017-11",
		EndDate:   "2022-04",
		Location:  "Melbourne, Australia",
		Highlights: []string{
			"Led Tools, Monitoring, Metrics, and SCM/CI/CD SRE teams scaling",
			"Architected and maintained CI/CD infrastructure",
			"Platform scaling across configuration management, metrics, logs, CI/CD, and Kubernetes",
			"Designed and implmeneted standard build templates",
			"Core member of building out first GCP environment",
		},
		Skills: []string{
			"GCP", "Salt", "Prometheus", "GitLab", "Terraform",
			"Kubernetes", "Python", "Go",
		},
	},
	{
		Id:        "linkedin-staff-sre",
		Company:   "LinkedIn",
		Title:     "Staff Site Reliability Engineer, Jobs & Recruiter",
		StartDate: "2015-02",
		EndDate:   "2017-11",
		Location:  "San Francisco Bay Area",
		Highlights: []string{
			"SRE for 100+ applications across thousands of hosts in a service-oriented architecture",
			"Mentored SRE team growth from 3 to 12 engineers",
			"Overhauled monitoring philosophy and implementation across multiple teams",
			"Created and lead many training & knowledge sharing efforts, speaking at external conferences.",
			"Started SRE scorecard for service improvement across all of linkedin",
			"Mentoring and training my growing SRE team (3 to 12), while preserving a site-up and members first culture",
		},
		Skills: []string{
			"Python", "Java", "Go", "Kafka", "containerd", "couchdb",
		},
	},
	{
		Id:        "medibank-devops",
		Company:   "Medibank Health Solutions",
		Title:     "Linux & Cloud Services Engineer (DevOps)",
		StartDate: "2013-03",
		EndDate:   "2015-02",
		Location:  "Melbourne, Australia",
		Highlights: []string{
			"Founding member of Cloud Services team scaling from 2 to 5 engineers",
			"Built Medibank's first AWS environment and PCI compliance automation",
			"Improved build times and realibility with docker & buildbot",
			"Contributed to core code-base and implemenated product and realibility features",
		},
		Skills: []string{
			"AWS", "Puppet", "Python", "django", "docker",
		},
	},
	{
		Id:        "toll-sysadmin",
		Company:   "Toll Holdings Limited",
		Title:     "Senior Linux & Unix Admin",
		StartDate: "2008-12",
		EndDate:   "2013-03",
		Location:  "Australia",
		Skills: []string{
			"solaris", "redhat", "Python", "django", "nagios", "graphite",
		},
	},
	{
		Id:        "ibm-unix-admin",
		Company:   "IBM",
		Title:     "Unix Admin",
		StartDate: "2005-12",
		EndDate:   "2008-12",
		Location:  "Australia",
		Skills: []string{
			"solaris", "redhat", "hpux", "aix",
		},
	},
}
