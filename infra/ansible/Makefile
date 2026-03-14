
init:
	uv init
	uv venv

run:
	./.venv/bin/ansible-playbook -i inventory.ini site.yml

k3d:
	./.venv/bin/ansible-playbook -i inventory.ini site.yml --tags k3d

nginx:
	./.venv/bin/ansible-playbook -i inventory.ini site.yml --tags nginx

sync-creds:
	cd $(shell git rev-parse --show-toplevel) && ./scripts/sync-prod-cluster-creds.sh