package objects

import "errors"

type User struct {
	Id          int
	Name        string
	Group       Group
	Permissions []string
}

func (u *User) SetGroup(g Group) error {
	if u == nil {
		return errors.New("user cannot be nil")
	}

	u.Group = g

	return nil
}

func (u *User) Grant(permission string) error {
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

	return nil
}

func (u *User) Can(permission string) (bool, error) {
	if u == nil || permission == "" {
		return false, errors.New("user or permission is nil")
	}

	for _, v := range u.Permissions {
		if v == permission {
			return true, nil
		}
	}

	for _, v := range u.Group.GetAllPermissions() {
		if v == permission {
			return true, nil
		}
	}

	return false, nil
}

func (u *User) Revoke(permission string) error {
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

	return nil
}
