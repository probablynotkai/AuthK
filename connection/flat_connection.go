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

type storedGroup struct {
	Id          int
	Name        string
	IsDefault   bool
	Permissions []string
	Inheritance []int
}

func (f *FlatConnection) Connect() {
	log.Println("Attempting to connect to flat file data source...")

	f.userPath = f.Directory + "\\users.json"
	f.groupPath = f.Directory + "\\groups.json"

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

func (f *FlatConnection) LoadGroups() {
	log.Println("Loading groups...")

	groupData, err := os.ReadFile(f.groupPath)
	if err != nil {
		log.Fatal(err)
		return
	}

	var storedGroups = []storedGroup{}
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

	var loadedGroups = []objects.Group{}

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

		loadedGroups = append(loadedGroups, *group)
	}

	groups = loadedGroups

	log.Println("Loaded groups and permissions...")
}

func (f *FlatConnection) Save() error {
	userData, err := json.MarshalIndent(users, "", "    ")
	if err != nil {
		return err
	}

	groupData, err := json.MarshalIndent(groups, "", "    ")
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

func (f *FlatConnection) CreateUser(name string) (*objects.User, error) {
	if name == "" {
		return nil, errors.New("identifier cannot be empty")
	}

	defaultGroup := f.GetDefaultGroup()
	if defaultGroup == nil {
		defaultGroup = &objects.Group{
			Id: -1,
		}
	}

	user := objects.User{
		Id:          getNextUserId(),
		Name:        name,
		Permissions: []string{},
		GroupId:     defaultGroup.Id,
	}

	users = append(users, user)

	return &user, f.Save()
}

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

func (f *FlatConnection) GetUsers() *[]objects.User {
	return &users
}

func (f *FlatConnection) GetGroups() *[]objects.Group {
	return &groups
}

func (f *FlatConnection) GetUser(id int) *objects.User {
	for _, v := range *f.GetUsers() {
		if v.Id == id {
			return &v
		}
	}

	return nil
}

func (f *FlatConnection) GetGroup(id int) *objects.Group {
	for _, v := range *f.GetGroups() {
		if v.Id == id {
			return &v
		}
	}
	return nil
}

func (f *FlatConnection) GetDefaultGroup() *objects.Group {
	for _, v := range *f.GetGroups() {
		if v.IsDefault {
			return &v
		}
	}

	return nil
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

// func isCyclicalInheritance(g objects.Group) bool {
// 	allGroups = []objects.Group{}
// 	getNestedGroups(g.Inheritance, &allGroups)
// }

// func getNestedGroups(g []objects.Group, a *[]objects.Group) {
// 	if len(g.Inheritance) == 0 {
// 		return
// 	} else {

// 	}
// }

/*
1. load all stored groups (inheritance is IDs)

2. iterate stored groups, load into normal storage excluding inheritance

3. iterate stored groups, get normal group by id (ref), iterate normal groups and ref inheritance ids, if normal group id = ref inheritance id, add to normal group inheritance (groups not ids)

4. for range normal group 1 inheritance groups, recurse each inheritance group's inheritance groups for normal group 1's id, if exists, panic, no cyclical inheritance allowed
*/
