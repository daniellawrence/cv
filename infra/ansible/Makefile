
init:
	uv init
	uv venv

run:
	./.venv/bin/ansible-playbook -i inventory.ini site.yml

k3d:
	./.venv/bin/ansible-playbook -i inventory.ini site.yml --tags k3d

nginx:
	./.venv/bin/ansible-playbook -i inventory.ini site.yml --tags nginx