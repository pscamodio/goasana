package goasana

import (
	"encoding/json"
	"errors"
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
	Created_at string
	Followers  []User
	Color      string
	Notes      string
	Workspace  Workspace
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
	Completed_at    string
	Due_on          string
	Followers       []User
	Modified_at     string
	Projects        []Project
	Parent          BaseData
	Workspace       Workspace
}

type AsanaConn struct {
	connInfo
	Connected bool
	Me        User
}

type ExtraParams struct {
	PrettyPrint    bool
	RequiredFields []string
	ExpandFields   []string
}

const (
	main_uri          string = "https://app.asana.com/api/1.0"
	users_uri         string = "/users"
	workspaces_uri    string = "/workspaces"
	me_uri            string = "/me"
	tasks_uri         string = "/tasks"
	subtasks_uri      string = "/subtasks"
	tags_uri          string = "/tags"
	projects_uri      string = "/projects"
	organizations_uri string = "/organizations"
	teams_uri         string = "/teams"
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

func (params *ExtraParams) ToFilters() (m map[string]string) {
	m = make(map[string]string)
	if params == nil {
		return
	}
	if params.PrettyPrint {
		m["opt_pretty"] = "true"
	}
	if params.RequiredFields != nil {
		fields := make([]string, 0)
		for _, f := range params.RequiredFields {
			fields = append(fields, f)
		}
		m["opt_fields"] = strings.Join(fields, ",")
	}
	if params.ExpandFields != nil {
		fields := make([]string, 0)
		for _, f := range params.ExpandFields {
			fields = append(fields, f)
		}
		m["opt_expand"] = strings.Join(fields, ",")
	}
	return
}

func NewSimpleConnection(apiKey string) (conn *AsanaConn, err error) {
	if apiKey == "" {
		return nil, errors.New("Empty Key")
	}
	cnInfo := connInfo{true, apiKey, ""}
	data, err := cnInfo.SendRequest("GET", main_uri+users_uri+me_uri)
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
	return &AsanaConn{connInfo: cnInfo, Connected: true, Me: temp.Data}, nil
}

func (conn AsanaConn) GetMe(params *ExtraParams) (me *User, err error) {
	data, err := conn.SendRequestWithFilters("GET", main_uri+users_uri+me_uri, params.ToFilters())
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

func (conn AsanaConn) GetUserData(user_id int, params *ExtraParams) (user *User, err error) {
	idstring := "/" + strconv.FormatInt(int64(user_id), 10)
	data, err := conn.SendRequestWithFilters("GET", main_uri+users_uri+idstring, params.ToFilters())
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

func (conn AsanaConn) GetUsers(params *ExtraParams) (users []User, err error) {
	data, err := conn.SendRequestWithFilters("GET", main_uri+users_uri, params.ToFilters())
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

func (conn AsanaConn) GetUsersFromWorkspace(workspace_id int, params *ExtraParams) (users []User, err error) {
	idstring := "/" + strconv.FormatInt(int64(workspace_id), 10)
	data, err := conn.SendRequestWithFilters("GET", main_uri+workspaces_uri+idstring+users_uri, params.ToFilters())
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

func (conn AsanaConn) GetTaskFromUser(workspace, userid int, archivied bool, params *ExtraParams) (tasks []Task, err error) {
	archivied_str := "false"
	if archivied {
		archivied_str = "true"
	}
	filters := params.ToFilters()
	filters["assignee"] = strconv.FormatInt(int64(userid), 10)
	filters["workspace"] = strconv.FormatInt(int64(workspace), 10)
	filters["include_archivied"] = archivied_str
	data, err := conn.SendRequestWithFilters("GET", main_uri+tasks_uri, filters)
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

func (conn AsanaConn) GetTaskData(taskid int, params *ExtraParams) (task *Task, err error) {
	taskstr := "/" + strconv.FormatInt(int64(taskid), 10)
	data, err := conn.SendRequestWithFilters("GET", main_uri+tasks_uri+taskstr, params.ToFilters())
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
	return &temp.Data, nil
}

func (conn AsanaConn) GetSubTask(taskid int, params *ExtraParams) (stasks []Task, err error) {
	idstr := "/" + strconv.FormatInt(int64(taskid), 10)
	data, err := conn.SendRequestWithFilters("GET", main_uri+tasks_uri+idstr+subtasks_uri, params.ToFilters())
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

func (conn AsanaConn) GetTagsFromTask(taskid int, params *ExtraParams) (tags []Tag, err error) {
	idstr := "/" + strconv.FormatInt(int64(taskid), 10)
	data, err := conn.SendRequestWithFilters("GET", main_uri+tasks_uri+idstr+tags_uri, params.ToFilters())
	if err != nil {
		return
	}
	var temp struct {
		Data   []Tag
		Errors []Error
	}
	json.Unmarshal(data, &temp)
	err = checkForErrors(temp.Errors)
	if err != nil {
		return
	}
	return temp.Data, nil
}

func (conn AsanaConn) GetWorkspaces(params *ExtraParams) (workspaces []Workspace, err error) {
	data, err := conn.SendRequestWithFilters("GET", main_uri+workspaces_uri, params.ToFilters())
	if err != nil {
		return
	}
	var temp struct {
		Data   []Workspace
		Errors []Error
	}
	json.Unmarshal(data, &temp)
	err = checkForErrors(temp.Errors)
	if err != nil {
		return
	}
	return temp.Data, nil
}

func (conn AsanaConn) GetWorkspaceData(workspace_id int, params *ExtraParams) (workspace *Workspace, err error) {
	idstr := "/" + strconv.FormatInt(int64(workspace_id), 10)
	data, err := conn.SendRequestWithFilters("GET", main_uri+workspaces_uri+idstr, params.ToFilters())
	if err != nil {
		return
	}
	var temp struct {
		Data   Workspace
		Errors []Error
	}
	json.Unmarshal(data, &temp)
	err = checkForErrors(temp.Errors)
	if err != nil {
		return
	}
	return &temp.Data, nil
}

func (conn AsanaConn) GetProjects(workspace_id int, params *ExtraParams) (projects []Project, err error) {
	filters := params.ToFilters()
	filters["workspace"] = strconv.FormatInt(int64(workspace_id), 10)
	data, err := conn.SendRequestWithFilters("GET", main_uri+projects_uri, filters)
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

func (conn AsanaConn) GetProjectData(project_id int, params *ExtraParams) (project *Project, err error) {
	idstr := "/" + strconv.FormatInt(int64(project_id), 10)
	data, err := conn.SendRequestWithFilters("GET", main_uri+projects_uri+idstr, params.ToFilters())
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
	return &temp.Data, nil
}

func (conn AsanaConn) GetTaskFromProject(project_id int, archivied bool, params *ExtraParams) (tasks []Task, err error) {
	archivied_str := "false"
	if archivied {
		archivied_str = "true"
	}
	filters := params.ToFilters()
	filters["project"] = strconv.FormatInt(int64(project_id), 10)
	filters["include_archivied"] = archivied_str
	data, err := conn.SendRequestWithFilters("GET", main_uri+tasks_uri, filters)
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

func (conn AsanaConn) GetTagData(tag_id int, params *ExtraParams) (tag *Tag, err error) {
	idstr := "/" + strconv.FormatInt(int64(tag_id), 10)
	data, err := conn.SendRequestWithFilters("GET", main_uri+tags_uri+idstr, params.ToFilters())
	if err != nil {
		return
	}
	var temp struct {
		Data   Tag
		Errors []Error
	}
	json.Unmarshal(data, &temp)
	err = checkForErrors(temp.Errors)
	if err != nil {
		return
	}
	return &temp.Data, nil
}

func (conn AsanaConn) GetTasksFromTag(tag_id int, params *ExtraParams) (tasks []Task, err error) {
	idstr := "/" + strconv.FormatInt(int64(tag_id), 10)
	data, err := conn.SendRequestWithFilters("GET", main_uri+tags_uri+idstr+tasks_uri, params.ToFilters())
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

func (conn AsanaConn) GetTags(params *ExtraParams) (tags []Tag, err error) {
	data, err := conn.SendRequestWithFilters("GET", main_uri+tags_uri, params.ToFilters())
	if err != nil {
		return
	}
	var temp struct {
		Data   []Tag
		Errors []Error
	}
	json.Unmarshal(data, &temp)
	err = checkForErrors(temp.Errors)
	if err != nil {
		return
	}
	return temp.Data, nil
}

func (conn AsanaConn) GetTagsFromWorkspace(workspace_id int, params *ExtraParams) (tags []Tag, err error) {
	idstr := "/" + strconv.FormatInt(int64(workspace_id), 10)
	data, err := conn.SendRequestWithFilters("GET", main_uri+workspaces_uri+idstr+tags_uri, params.ToFilters())
	if err != nil {
		return
	}
	var temp struct {
		Data   []Tag
		Errors []Error
	}
	json.Unmarshal(data, &temp)
	err = checkForErrors(temp.Errors)
	if err != nil {
		return
	}
	return temp.Data, nil
}

func (conn AsanaConn) GetTeams(organization_id int, params *ExtraParams) (teams []Team, err error) {
	idstr := "/" + strconv.FormatInt(int64(organization_id), 10)
	data, err := conn.SendRequestWithFilters("GET", main_uri+organizations_uri+idstr+teams_uri, params.ToFilters())
	if err != nil {
		return
	}
	var temp struct {
		Data   []Team
		Errors []Error
	}
	json.Unmarshal(data, &temp)
	err = checkForErrors(temp.Errors)
	if err != nil {
		return
	}
	return temp.Data, nil
}
