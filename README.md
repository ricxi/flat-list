# flat-list
A to-do list application written as a Go microservice. I've reading a lot lately about the different ways you can structure a Go application, so I decided to experiment a bit with the classic *flat* structure.

## organization
```
.
├── dev_scripts
├── frontend-client
├── mailer
├── migrations
├── shared
├── task
├── token
└── user
```
* dev_scripts: a few automation scripts to help with automating the development of this project (a lot of throw aways)
* user: microservice for handling interactions
* `dev.bash` - a small cli utility tool; run `source bash` if you want to use it in your terminal (requires bash)
* mailer - this service was much smaller, so I stored the mocks with their tests.

## Notes
If running outside of docker container:
* golang-migrate must be installed (I could dl this with the go toolchain or write a package to handle migrations?)
* .env and config files must be set up (I might switch this to config files)
* mailer service must be disabled or SMTP server credentials must be provided