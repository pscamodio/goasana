package goasana

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Error struct {
	Message string
}

type BaseData struct {
	Id   int
	Name string
}

type Tag struct {
	BaseData
}

type Team struct {
	BaseData
}

type Project struct {
	BaseData
	Archivied  bool
	Created_at string
	Followers  []User
	Color      string
	Notes      string
	Workspace  Workspace
	Team       Team
}

type Workspace struct {
	BaseData
	Is_Organization bool
}

type User struct {
	BaseData
	Email      string
	Photo      map[string]string
	Workspaces []Workspace
}

type Task struct {
	BaseData
	Assignee        User
	Assignee_status string
	Created_at      string
	Completed       bool
	Completed_ad    string
	Due_on          string
	Followers       []User
	Modified_at     string
	Projects        []Project
	Parent          BaseData
	Workspace       Workspace
}

const (
	main_uri          string = "https://app.asana.com/api/1.0"
	users_uri         string = "/users"
	workspaces_uri    string = "/workspaces"
	me_uri            string = "/me"
	tasks_uri         string = "/tasks"
	projects_uri      string = "/projects"
	organizations_uri string = "/organizations"
)

func checkForErrors(err []Error) error {
	if len(err) == 0 {
		return nil
	}
	lines := make([]string, len(err))
	for i, e := range err {
		lines[i] = e.Message
	}
	return errors.New(strings.Join(lines, "\n"))
}

func GetMe() (me *User, err error) {
	data, err := SendRequest("GET", main_uri+users_uri+me_uri)
	if err != nil {
		return
	}
	var temp struct {
		Data   User
		Errors []Error
	}
	json.Unmarshal(data, &temp)
	err = checkForErrors(temp.Errors)
	if err != nil {
		return
	}
	return &temp.Data, nil
}

func GetUsers() (users []User, err error) {
	data, err := SendRequest("GET", main_uri+users_uri)
	if err != nil {
		return
	}
	var temp struct {
		Data   []User
		Errors []Error
	}
	json.Unmarshal(data, &temp)
	err = checkForErrors(temp.Errors)
	if err != nil {
		return
	}
	return temp.Data, nil
}

func GetUsersFromWorkspace(workspace_id int) (users []User, err error) {
	idstring := "/" + strconv.FormatInt(int64(workspace_id), 10)
	data, err := SendRequest("GET", main_uri+workspaces_uri+idstring+users_uri)
	if err != nil {
		return
	}
	var temp struct {
		Data   []User
		Errors []Error
	}
	json.Unmarshal(data, &temp)
	err = checkForErrors(temp.Errors)
	if err != nil {
		return
	}
	return temp.Data, nil
}

func GetTaskFromUser(workspace, userid int) (tasks []Task, err error) {
	filters := map[string]string{
		"assignee":  strconv.FormatInt(int64(userid), 10),
		"workspace": strconv.FormatInt(int64(workspace), 10)}
	data, err := SendRequestWithFilters("GET", main_uri+tasks_uri, filters)
	if err != nil {
		return
	}
	var temp struct {
		Data   []Task
		Errors []Error
	}
	json.Unmarshal(data, &temp)
	err = checkForErrors(temp.Errors)
	if err != nil {
		return
	}
	return temp.Data, nil
}

func GetTaskData(taskid int) (task Task, err error) {
	taskstr := "/" + strconv.FormatInt(int64(taskid), 10)
	data, err := SendRequest("GET", main_uri+tasks_uri+taskstr)
	if err != nil {
		return
	}
	var temp struct {
		Data   Task
		Errors []Error
	}
	json.Unmarshal(data, &temp)
	err = checkForErrors(temp.Errors)
	if err != nil {
		return
	}
	return temp.Data, nil
}

func GetWorkspaces() (workspaces []Workspace, err error) {
	data, err := SendRequest("GET", main_uri+workspaces_uri)
	if err != nil {
		return
	}
	var temp struct {
		Data   []Workspace
		Errors []Error
	}
	fmt.Println(string(data))
	json.Unmarshal(data, &temp)
	err = checkForErrors(temp.Errors)
	if err != nil {
		return
	}
	return temp.Data, nil
}

func GetProjects(workspace_id int) (projects []Project, err error) {
	filters := map[string]string{
		"workspace": strconv.FormatInt(int64(workspace_id), 10)}
	data, err := SendRequestWithFilters("GET", main_uri+projects_uri, filters)
	if err != nil {
		return
	}
	var temp struct {
		Data   []Project
		Errors []Error
	}
	json.Unmarshal(data, &temp)
	err = checkForErrors(temp.Errors)
	if err != nil {
		return
	}
	return temp.Data, nil
}

func GetProjectData(project_id int) (project Project, err error) {
	idstr := "/" + strconv.FormatInt(int64(project_id), 10)
	data, err := SendRequest("GET", main_uri+projects_uri+idstr)
	fmt.Println(main_uri + projects_uri + idstr)
	if err != nil {
		return
	}
	var temp struct {
		Data   Project
		Errors []Error
	}
	json.Unmarshal(data, &temp)
	err = checkForErrors(temp.Errors)
	if err != nil {
		return
	}
	return temp.Data, nil
}
