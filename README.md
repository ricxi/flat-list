# flat-list
A to-do list application written as a Go microservice. I've reading a lot lately about the different ways you can structure a Go application, so I decided to experiment a bit with the classic *flat* structure.

## organization
```
.
├── dev_scripts
├── frontend-client
├── mailer
├── migrations
├── runner
├── shared
├── task
├── token
└── user
```
* dev_scripts: a few automation scripts to help with automating the development of this project (a lot of throw aways)
* user: microservice for handling interactions
* `dev.bash` - a small cli utility tool; run `source bash` if you want to use it in your terminal (requires bash)