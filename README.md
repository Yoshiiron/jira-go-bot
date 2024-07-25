# Jira-Bot
Jira-Go-Bot is a part of Jira-Bot project.
Jira-Bot is a project, that allow users to recieve a notifications in telegram.
For users it works like:
They changed something in issue and user who is assignee or creater(depends on who send webhook) gets message from bot with changes in telegram.

## Jira-Go-Bot
Jira-Go-Bot can register new user and send them messages from Jira.

### New User Registration
User activate the bot and send him "Зарегистрироваться". When this happen, bot starting a registration process.
Activation occurs in three stages:
1) User sends bot his login on Jira
2) Bot parsing the email in Jira and send code on it
3) User send bot a code from his email
4) If code is valid - user succecly registred
