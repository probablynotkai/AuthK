package connection

import (
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/probablynotkai/objects"
)

var users []objects.User
var groups []objects.Group

type FlatConnection struct {
	Directory string

	userPath  string
	groupPath string
}

type StoredGroup struct {
	Id          int
	Name        string
	IsDefault   bool
	Permissions []string
	Inheritance []int
}

type StoredUser struct {
	Id          int
	Name        string
	GroupId     int
	Permissions []string
}

/*
Function to initialise connection.
*/
func (f *FlatConnection) Connect() {
	log.Println("Attempting to connect to flat file data source...")

	// Leave blank for current directory
	if f.Directory == "" {
		f.userPath = "users.json"
		f.groupPath = "groups.json"
	} else {
		f.userPath = f.Directory + "\\users.json"
		f.groupPath = f.Directory + "\\groups.json"
	}

	if _, err := os.Stat(f.userPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Println("users.json doesn't exist, creating now...")
			err := os.WriteFile(f.userPath, []byte{}, os.ModePerm)
			if err != nil {
				log.Fatal(err)
				return
			}
		} else {
			log.Fatal(err)
			return
		}
	}

	if _, err := os.Stat(f.groupPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Println("groups.json doesn't exist, creating now...")
			err := os.WriteFile(f.groupPath, []byte{}, os.ModePerm)
			if err != nil {
				log.Fatal(err)
				return
			}
		} else {
			log.Fatal(err)
			return
		}
	}

	log.Println("Connected to flat file data source.")

	f.LoadGroups()
	f.LoadUsers()
}

/*
Loads stored users.
*/
func (f *FlatConnection) LoadUsers() {
	log.Println("Loading users...")

	userData, err := os.ReadFile(f.userPath)
	if err != nil {
		log.Fatal(err)
		return
	}

	var storedUsers []objects.User
	json.Unmarshal(userData, &storedUsers)

	users = storedUsers

	log.Println("Loaded users and permissions.")
}

/*
Loads stored groups.
*/
func (f *FlatConnection) LoadGroups() {
	log.Println("Loading groups...")

	groupData, err := os.ReadFile(f.groupPath)
	if err != nil {
		log.Fatal(err)
		return
	}

	var storedGroups = []StoredGroup{}
	json.Unmarshal(groupData, &storedGroups)

	// Load excl. inheritance
	for _, v := range storedGroups {
		groups = append(groups, objects.Group{
			Id:          v.Id,
			Name:        v.Name,
			IsDefault:   v.IsDefault,
			Permissions: v.Permissions,
		})
	}

	// Insert loaded groups to respective inheritance
	for _, v := range storedGroups {
		group := f.GetGroup(v.Id)

		if group == nil {
			continue
		}

		for _, w := range v.Inheritance {
			if w == v.Id {
				continue
			}

			parent := f.GetGroup(w)

			group.Inheritance = append(group.Inheritance, *parent)
		}
	}

	log.Println("Loaded groups and permissions...")
}

/*
Saves application data. Execute after granting or revoking permissions.
*/
func (f *FlatConnection) Save() error {
	userData, err := json.MarshalIndent(formatUsers(users), "", "    ")
	if err != nil {
		return err
	}

	groupData, err := json.MarshalIndent(formatGroups(groups), "", "    ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(f.userPath, userData, os.ModePerm); err != nil {
		return err
	}
	if err := os.WriteFile(f.groupPath, groupData, os.ModePerm); err != nil {
		return err
	}

	return nil
}

/*
Creates a user with the specified name, will attempt to add to default group.
*/
func (f *FlatConnection) CreateUser(name string) (*objects.User, error) {
	var (
		user objects.User
	)

	if name == "" {
		return nil, errors.New("identifier cannot be empty")
	}

	defaultGroup := f.GetDefaultGroup()

	if defaultGroup == nil {
		user = objects.User{
			Id:          getNextUserId(),
			Name:        name,
			Permissions: []string{},
		}
	} else {
		user = objects.User{
			Id:          getNextUserId(),
			Name:        name,
			Permissions: []string{},
			Group:       *defaultGroup,
		}
	}

	users = append(users, user)

	return &user, f.Save()
}

/*
Creates a group with the specified name.
*/
func (f *FlatConnection) CreateGroup(name string) (*objects.Group, error) {
	if name == "" {
		return nil, errors.New("identifier cannot be empty")
	}

	group := objects.Group{
		Id:          getNextGroupId(),
		Name:        name,
		IsDefault:   false,
		Permissions: []string{},
		Inheritance: []objects.Group{},
	}

	groups = append(groups, group)

	return &group, f.Save()
}

/*
Returns all users.
*/
func (f *FlatConnection) GetUsers() *[]objects.User {
	return &users
}

/*
Returns all groups.
*/
func (f *FlatConnection) GetGroups() *[]objects.Group {
	return &groups
}

/*
Returns user with ID, returns nil if none.
*/
func (f *FlatConnection) GetUser(id int) *objects.User {
	for i := range *f.GetUsers() {
		user := &(*f.GetUsers())[i]
		if user.Id == id {
			return user
		}
	}

	return nil
}

/*
Returns group with specified ID, returns nil if none.
*/
func (f *FlatConnection) GetGroup(id int) *objects.Group {
	for i := range *f.GetGroups() {
		group := &(*f.GetGroups())[i]
		if group.Id == id {
			return group
		}
	}
	return nil
}

/*
Returns group with specified name, returns first in list and nil if none.
*/
func (f *FlatConnection) GetGroupByName(name string) *objects.Group {
	for i := range *f.GetGroups() {
		group := &(*f.GetGroups())[i]
		if group.Name == name {
			return group
		}
	}
	return nil
}

/*
Returns the default stored group, returns nil if none.
*/
func (f *FlatConnection) GetDefaultGroup() *objects.Group {
	for _, v := range *f.GetGroups() {
		if v.IsDefault {
			return &v
		}
	}

	return nil
}

func formatUsers(users []objects.User) []StoredUser {
	formattedUsers := []StoredUser{}

	for _, v := range users {
		formattedUsers = append(formattedUsers, StoredUser{
			Id:          v.Id,
			Name:        v.Name,
			GroupId:     v.Group.Id,
			Permissions: v.Permissions,
		})
	}

	return formattedUsers
}

func formatGroups(groups []objects.Group) []StoredGroup {
	formattedGroups := []StoredGroup{}

	for _, v := range groups {
		inheritanceIds := []int{}

		for _, i := range v.Inheritance {
			inheritanceIds = append(inheritanceIds, i.Id)
		}

		formattedGroups = append(formattedGroups, StoredGroup{
			Id:          v.Id,
			Name:        v.Name,
			IsDefault:   v.IsDefault,
			Inheritance: inheritanceIds,
			Permissions: v.Permissions,
		})
	}

	return formattedGroups
}

func getNextUserId() int {
	var (
		last = 0
	)
	for _, v := range users {
		if v.Id > last {
			last = v.Id
		}
	}
	return last + 1
}

func getNextGroupId() int {
	var (
		last = 0
	)
	for _, v := range groups {
		if v.Id > last {
			last = v.Id
		}
	}
	return last + 1
}
