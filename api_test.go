package goasana

import (
	"testing"
)

func init() {

}

const (
	myname           string = "Amodio Pesce"
	myid             int    = 8541370844451
	myfirstworkspace int    = 8541391811820
)

func TestGetMe(t *testing.T) {
	me, err := GetMe()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if me.Name != "Amodio Pesce" {
		t.Log("Name not match")
		t.FailNow()
	}
	if me.Workspaces[0].Id != myfirstworkspace {
		t.Log("Workspace not match")
		t.FailNow()
	}
	t.Log(me.Id)
}

func TestGetUsers(t *testing.T) {
	users, err := GetUsers()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	for _, user := range users {
		t.Log(user.Name)
	}
}

func TestGetUsersFromWorkspace(t *testing.T) {
	users, err := GetUsersFromWorkspace(myfirstworkspace)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	for _, user := range users {
		t.Log(user.Name)
	}
}

func TestGetAllTasks(t *testing.T) {
	_, err := GetTaskFromUser(myfirstworkspace, myid)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
}

func TestGetAllTaskWithInfo(t *testing.T) {
	tasks, err := GetTaskFromUser(myfirstworkspace, myid)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	for idx, task := range tasks {
		tasks[idx], err = GetTaskData(task.Id)
		if err != nil {
			t.Log(err)
			t.Fail()
			err = nil
		}
		t.Log(tasks[idx])
	}
}
