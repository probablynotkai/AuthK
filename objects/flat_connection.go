package objects

import (
	"encoding/json"
	"errors"
	"log"
	"os"
)

var users []User

type FlatConnection struct {
	FileLocation string
}

func (f FlatConnection) Connect() {
	log.Println("Attempting to connect to flat file data source...")

	_, err := os.Stat(f.FileLocation)
	if err != nil {
		log.Fatal(err)
		return
	}

	data, err := os.ReadFile(f.FileLocation)
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Println("Connected to flat file data source.")

	var permissions Permissions
	json.Unmarshal(data, &permissions)

	log.Println("Loading users and permissions...")

	for k, v := range permissions.PermissionsMap {
		permissionArray := []string{}
		for _, va := range v {
			permissionArray = append(permissionArray, va)
		}

		users = append(users, User{
			Name:        k,
			Permissions: permissionArray,
		})
	}

	log.Println("Loaded users and permissions.")
}

func (f *FlatConnection) Save() error {
	data, err := json.Marshal(users)
	if err != nil {
		return err
	}

	os.WriteFile(f.FileLocation, data, os.ModePerm)

	return nil
}

func (f *FlatConnection) CreateUser(identifier string) (*User, error) {
	if identifier == "" {
		return nil, errors.New("identifier cannot be nil")
	}

	user := User{
		Name:        identifier,
		Permissions: []string{},
	}

	users = append(users, user)

	return &user, f.Save()
}

func (f *FlatConnection) Grant(u *User, permission string) error {
	if u == nil || permission == "" {
		return errors.New("user and permission cannot be nil")
	}

	alreadyHas := false
	for _, v := range u.Permissions {
		if v == permission {
			alreadyHas = true
		}
	}

	if !alreadyHas {
		u.Permissions = append(u.Permissions, permission)
	}

	return f.Save()
}

func (f *FlatConnection) Can(u *User, permission string) (bool, error) {
	if u == nil || permission == "" {
		return false, errors.New("user or permission is nil")
	}

	for _, v := range u.Permissions {
		if v == permission {
			return true, nil
		}
	}

	return false, nil
}

func (f *FlatConnection) Revoke(u *User, permission string) error {
	if u == nil || permission == "" {
		return errors.New("user and permission cannot be nil")
	}

	index := -1
	for i := 0; i < len(u.Permissions); i++ {
		if u.Permissions[i] == permission {
			index = i
		}
	}

	if index == -1 {
		return nil
	}

	// some fucking voodoo magic
	u.Permissions[index] = u.Permissions[len(u.Permissions)-1] // replace target with end element
	u.Permissions = u.Permissions[:len(u.Permissions)-1]       // replace with array -1 length

	return f.Save()
}

func (f *FlatConnection) GetUsers() []User {
	return users
}
