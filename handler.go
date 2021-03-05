package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func displayMainPage(c *gin.Context) {
	c.HTML(http.StatusOK, "main", nil)
}

func displayLanding(c *gin.Context) {
	c.HTML(http.StatusOK, "landing", nil)
}

func displayRegister(c *gin.Context) {
	c.HTML(http.StatusOK, "registration", "Create Your Account")
}

func displayLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "login", nil)
}

func processRegistration(c *gin.Context) {
	user := c.PostForm("user")
	pass := c.PostForm("pass")
	if userNameExists(user) {
		log.Println("Invalid User")
		c.HTML(http.StatusTemporaryRedirect, "registration", "User Name is not Available")
		return
	}
	if pass == "" {
		log.Println("Bad Password")
		c.HTML(http.StatusTemporaryRedirect, "registration", "Password cannot be blank")
		return
	}
	log.Println("Processed registartion, saving new user ", user)
	saveNewUser(user, pass)
	displayLogin(c)
}

func processLogin(c *gin.Context) {

	user := c.PostForm("user")
	pass := c.PostForm("pass")
	if validateUser(user, pass) {
		session := sessions.Default(c)
		session.Set("loggedIn", true)
		session.Save()
		fmt.Println("saving session after login")
		fmt.Println(session)
		displayMainPage(c)
	} else {
		c.HTML(http.StatusUnauthorized, "login", nil)
	}
}

func displayAddPage(c *gin.Context) {
	c.HTML(http.StatusOK, "add", nil)
}

func processAddTask(c *gin.Context) {
	var tasklist []Task
	taskName := c.PostForm("taskName")
	saveTask(taskName)
	found, err := Tasks.Get(&tasklist)
	for _, task := range tasklist {
		fmt.Println(task.ID, ":  ", task.Name)
	}
	if err != nil {
		log.Println(err, found, "error reading tasklist")
	}
	fmt.Println("task list was found", found)
	fmt.Println("sending task list to main page", tasklist)
	c.HTML(http.StatusCreated, "main", tasklist)
}
