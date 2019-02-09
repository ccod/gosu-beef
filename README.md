## Gosu Settle The Beef

currently a work in progress

### What we have:
- A way to login, so I can verify who the user is, and if they are part of the guild
- Means to get summary information about the player, namely their solo ladder record

### What next
- Persistence layer, save the initially loaded information to Postgres so subsequent logins are quicker
- Set up the actual tiers, so the app becomes more than an empty blizzard login service
- Establish a group of admins that will let the app know if the challengers succeeded
- Set up a cron job to keep the information on the guild members up to date.

### Wishes
- Set up a notification service for challenges?
- A means of setting up events and queues for those challenges
- investigate api to see if there is a way to automate changes in ranking after challanges are carried out
  - maybe track match history of challengers for custom games that conform to a naming convention?

### Tools used to make the app
- Use Kubernetes and docker for creating deployments on google cloud
- Using Golang on the backend to keep the services fast and small
- Using Nginx and React-App for the frontend
- Needed to fork https://github.com/mitchellh/go-bnet, needed to tweak endpoints for use as blizzard login tool
  - https://github.com/ccod/go-bnet is where you can find some of my tweaks
- Use https://github.com/FuzzyStatic/blizzard for api calls after I determine who the user is.