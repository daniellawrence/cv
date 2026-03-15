# Agent Guidelines for CV Repository

## Overview
This document provides guidelines for agents working on the CV repository, particularly around environment-specific configurations and production safety.

# Rules

- never break any rule or any reason
- You are never allowed to run git-commit or git-push to master, main or production
- You are allowed to git-commit and push into feat-{short-name} branches, if permission is granted.


## Reference

* binaries (helm, helmfile, kubectl) are kept in ./bin/
* scripts contain helpful wrappers for those binaries, but not always required
* backend/<service> - backends written in golang
* backend/common    - shared code between backends
* frontend          - react frontend that calls into backends
* infra             - deployment and infrastructure realted
* infra/ansible     - production configuraiton
* infra/argocd      - deployment
* infra/crds        - ignore
* deployments/artifacts    - auto-generated
* deployments/charts       - charts used to deploy the backend and frontend
* deployments/env          - values that are different between prod + dev
* deployment/helmfile.yaml - used to generate k8s manifests 
* proto             - protobufs, try not to change this
* references        - ignore 


You might need to inspect generated files, but never edit them directly.

*`infra/deployments/artifacts/ (auto-generated)
* gen (auto-generated)