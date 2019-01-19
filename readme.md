# Announce command

# What does this do?
- Send event announcements to multiple slack channels:
    Polyhack channel
    TorontoJS
    Techmasters
- Pin event announcement in slack for polyhack group
- Post event announcement to twitter

# How does this work?
- Uses AWS Lambda to execute commands
- Code in golang
- Interacts with the twitter and slack APIs