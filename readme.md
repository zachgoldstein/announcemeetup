# Announce command

# What does this do?
- Walks through a file path given by a file path and tries to send an announcement for each yaml file it can.
- Send event announcements to multiple endpoints:
  - Slack: Post and pin message up in a channel.
  - Twitter: Post up message on twitter

# To run the command:
```
go build
./announcemeetup ./data/
```