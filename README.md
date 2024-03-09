# Jira-Go-Bot

## For cloning:
docker build -t example --build-arg ssh_prv_key="$(cat ../bot-keys/repo-key)" --build-arg ssh_pub_key="$(cat ../bot-keys/repo-key.pub)" .