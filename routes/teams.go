package routes

import (
	"strings"
	"net/http"

	"github.com/ggpd/brackets/env"
	"github.com/gin-gonic/gin"

	"fmt"	
)

func (e *Env) GetTeamRoute(c *gin.Context) {
	token, err := c.Cookie("user_session")

	var login *env.User

	if err == nil {
		login, err = e.Db.CheckSession(token)
	}

	team, err := e.Db.GetTeam(c.Param("selector"), true)
	if err != nil {
		e.Log.Println(err)
		c.HTML(http.StatusNotFound, "notfound.html", nil)
		return
	}

	perm := env.Member
	if login != nil {
		perm, err = e.Db.GetRank(*team, *login)
	}

	var teams []env.Team
	teams = append(teams, *team)

	games, err := e.Db.GetUpcomingGames(teams)

	c.HTML(http.StatusOK, "team_index.html", gin.H{
		"login": login,
		"team":  team,
		"games": games,
		"rank": perm,
	})
}

func (e *Env) GetEditTeamRoute(c *gin.Context) {
	token, err := c.Cookie("user_session")
	
	if err != nil {
		c.HTML(http.StatusNotFound, "notfound.html", nil)
		return
	}
	
		login, err := e.Db.CheckSession(token)
	
		team, err := e.Db.GetTeam(c.Param("selector"), true)
		if err != nil {
			e.Log.Println(err)
			c.HTML(http.StatusNotFound, "notfound.html", nil)
			return
		}

		perm, err := e.Db.GetRank(*team, *login)
		if perm != env.Manager {
			e.Log.Println(err)
			c.HTML(http.StatusNotFound, "notfound.html", nil)
			return
		}
	
		c.HTML(http.StatusOK, "team_edit.html", gin.H{
			"login": login,
			"team":  team,
		})
}

func (e *Env) GetDeletePlayerRoute(c *gin.Context) {
	token, err := c.Cookie("user_session")
	
		var login *env.User
	
		if err == nil {
			login, err = e.Db.CheckSession(token)
		}
	
		team, err := e.Db.GetTeam(c.Param("selector"), true)
		if err != nil {
			e.Log.Println(err)
			c.HTML(http.StatusNotFound, "notfound.html", nil)
			return
		}

		perm, err := e.Db.GetRank(*team, *login)
		if perm != env.Manager {
			e.Log.Println(err)
			c.HTML(http.StatusNotFound, "notfound.html", nil)
			return
		}


		delSt := c.Param("user")

		var del *env.Player
		index := 0
		for _, p := range team.Players {
			if strings.Compare(delSt, p.Selector.String()) == 0 {
				del = p
				break
			}
			index++
		}

		team.Players[len(team.Players)-1], team.Players[index] = team.Players[index], team.Players[len(team.Players)-1]
		team.Players = team.Players[:len(team.Players)-1]

		team, err = e.Db.DeletePlayer(*team, *del)
		if err != nil {
			e.Log.Println(err)
		}

		c.Redirect(http.StatusFound, fmt.Sprintf("/team/%s/edit", team.Selector.String()))
		c.HTML(http.StatusOK, "team_edit.html", gin.H{
			"login": login,
			"team":  team,
		})
}

func (e *Env) PostEditTeamRoute(c *gin.Context) {
	token, err := c.Cookie("user_session")
	
		var login *env.User
	
		if err == nil {
			login, err = e.Db.CheckSession(token)
		}
	
		team, err := e.Db.GetTeam(c.Param("selector"), true)
		if err != nil {
			e.Log.Println(err)
			c.HTML(http.StatusNotFound, "notfound.html", nil)
			return
		}

		perm, err := e.Db.GetRank(*team, *login)
		if perm != env.Manager {
			e.Log.Println(err)
			c.HTML(http.StatusNotFound, "notfound.html", nil)
			return
		}
		

		teamName, er1 := c.GetPostForm("name")
		if !validField(teamName, er1) {
			//Redo
		}

		if strings.Compare(teamName, team.Name) != 0 {
			team.Name = teamName
			err = e.Db.UpdateTeam(*team)
			if err != nil {
				e.Log.Println(err)
			}
		}

		list := team.Players
		for _, pl := range list {
			rank, er2 := c.GetPostForm(fmt.Sprintf("%s:%s", pl.Selector.String(), "rank"))

			if !validField(rank, er2){
				//ERROR
			}

			r := env.ToRank(rank)

			if r == pl.Rank {
				continue
			}

			pl.Rank = r
			team, err = e.Db.UpdatePlayer(*team, *pl)
			if err != nil {
				e.Log.Println(err)
			}

		}
	
		c.HTML(http.StatusOK, "team_edit.html", gin.H{
			"login": login,
			"team":  team,
		})
}


func (e *Env) PostAddPlayerRoute(c *gin.Context) {
	token, err := c.Cookie("user_session")
	
		var login *env.User
	
		if err == nil {
			login, err = e.Db.CheckSession(token)
		}
	
		team, err := e.Db.GetTeam(c.Param("selector"), true)
		if err != nil {
			e.Log.Println(err)
			c.HTML(http.StatusNotFound, "notfound.html", nil)
			return
		}

		perm, err := e.Db.GetRank(*team, *login)
		if perm != env.Manager {
			e.Log.Println(err)
			c.HTML(http.StatusNotFound, "notfound.html", nil)
			return
		}

		uString, er1 := c.GetPostForm("new_sel")
		rank, er2 := c.GetPostForm("new_rank")
		if !validField(uString, er1) || !validField(rank, er2) {
			return
		}

		usr, err := e.Db.GetUser(uString)
		if err != nil {
			return
		}

		r := env.ToRank(rank)

		pl := env.Player {
			User: *usr,
			Rank: r,
		}


		team, err = e.Db.AddPlayer(*team, pl)

		c.Redirect(http.StatusFound, fmt.Sprintf("/team/%s/edit", team.Selector.String()))
		c.HTML(http.StatusOK, "team_edit.html", gin.H{
			"login": login,
			"team":  team,
		})	
	}

	func (e *Env) PostCreateTeamRoute(c *gin.Context){
		token, err := c.Cookie("user_session")
		
			var login *env.User
		
			if err == nil {
				login, err = e.Db.CheckSession(token)
			}
		
			tour, err := e.Db.GetTournament(c.Param("selector"), true)
			if err != nil {
				e.Log.Println(err)
				c.HTML(http.StatusNotFound, "notfound.html", nil)
				return
			}

			teamName, er1 := c.GetPostForm("new_team")
			if !validField(teamName, er1)  {
				return
			}

			team := env.Team {
				TournamentID: tour.ID,
				Name: teamName,
			}

			t, err := e.Db.CreateTeam(team)
			if err != nil {
				e.Log.Println(err)
			}

			pl := env.Player{
				User: *login,
				Rank: env.Manager,
			}

			t, err = e.Db.AddPlayer(*t, pl)
			if err != nil {
				e.Log.Println(err)
			}

			e.Log.Println(t.Selector.String())
			tour.Teams = append(tour.Teams, t)

			c.Redirect(http.StatusFound, fmt.Sprintf("/tournament/%s", tour.Selector.String()))
			c.HTML(http.StatusOK, "tournament_index.html", gin.H{
				"login": login,
				"tour": tour,
			})
	}