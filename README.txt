# brackets
Git: https://github.com/ggpd/brackets
CI: https://ci.amviego.org
    Provides updated binaries meant to be run on Mac OS, but I can chage to any OS you need
    User: brown
    Password: qwerty

config-example.yml needs to be updated with a different database and renamed to config.yml

If you would like to compile it yourself install at least go1.8.

go get github.com/ggpd/brackets
go get -u github.com/golang/dep/cmd/dep
cd $GOPATH/src/github.com/ggpd/brackets
dep ensure
go run main.go

Best way to access the site is http://71.176.106.101:8090/ which will be up to date.

# Example Users

gdumbrallkk@behance.net
9wipovgF5bJl

rlebbernkl@cornell.edu
vN9g3T6D89KG

# Implemented Features

## Users
- Register a User
    - signup button in top right corner
- Login a previously registered User
    - login button in top right corner
- Logout
    - login button in top right corner
- Edit
    - click on user's name in top right
- View an overview of user information on main page after login
- My profile page 
    - on navigation bar after login
- Pagenated users page

## Tournaments
- Pagenated tournaments page
    - on navigation bar
    - lets you create a new tournament by entering name into box while logged in
- Tournament index page lists organizers and current teams
    - click on any tournament after visiting tournaments page
    - create a new team by putting the name in the box while logged in

## Teams
- Teams index  
    - go to the pagenated tournaments page then go to any random
    tournament and click on a team to view upcoming games and players
- Teams Edit
    - go to teams index page and click on the Edit Team button on the navigation
    bar