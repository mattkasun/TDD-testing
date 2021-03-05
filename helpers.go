package main

import (
	"fmt"
	"log"

	"github.com/lithammer/shortuuid"

	"golang.org/x/crypto/bcrypt"
)

var Users, Tasks *Storage

func init() {
	Users = CreateStore("data", "users")
	Tasks = CreateStore("data", "tasks")
	//options := file.DefaultOptions
	//options.Directory = "data"
	// Create client
	//client, err := file.NewStore(options)
	//if err != nil {
	//	panic(err)
	//}
	//	defer client.Close()

}

func userNameExists(user string) bool {
	var userlist []User
	found, err := Users.Get(&userlist)

	if err != nil {
		fmt.Println("error retrieving uselist", err)
		return false
	}
	if !found {
		return false
	}
	for _, existing := range userlist {
		if user == existing.Name {
			return true
		}
	}
	return false
}
func saveTask(taskName string) {
	id := shortuuid.New()
	task := Task{ID: id, Name: taskName}
	fmt.Println("tasklist to be saved", id, taskName, task)
	var tasklist []Task
	found, err := Tasks.Get(&tasklist)
	if err != nil {
		log.Println(err, found, "error reading tasklist")
	}
	tasklist = append(tasklist, task)
	if Tasks.Save(tasklist) != nil {
		panic("tasklist was not saved")
	}
	fmt.Println("saved tasklist ", tasklist)
}

func saveNewUser(newuser, newpass string) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(newpass), 4)
	if err != nil {
		log.Fatal("error encoding password", err)
	}
	user := User{Name: newuser, Pass: string(bytes)}
	var userlist []User
	found, err := Users.Get(&userlist)
	if err != nil {
		log.Println(err, found, "error reading userlist")
	}
	userlist = append(userlist, user)
	if Users.Save(userlist) != nil {
		panic("userlist was not saved")
	}
}

func validateUser(u, p string) bool {
	var userlist []User
	found, err := Users.Get(&userlist)
	if (err != nil) || (!found) {
		return false
	}
	for _, user := range userlist {
		if user.Name == u && checkPassword(p, user.Pass) {
			return true
		}
	}
	log.Println("no such user", u, p)
	return false
}

func checkPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
