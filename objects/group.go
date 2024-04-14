package objects

import (
	"errors"
)

type Group struct {
	Id          int
	Name        string
	IsDefault   bool
	Permissions []string
	Inheritance []Group
}

func (g *Group) IsChild(parentId int) (bool, error) {
	if parentId < 0 {
		return false, errors.New("parent id cannot be below 0")
	}

	for _, v := range g.Inheritance {
		if v.Id == parentId {
			return true, nil
		}
	}

	return false, nil
}

func (g *Group) AddParent(parent *Group) error {
	if parent == nil {
		return errors.New("parent cannot be nil")
	}

	for _, v := range g.Inheritance {
		if v.Id == parent.Id {
			return nil
		}
	}

	g.Inheritance = append(g.Inheritance, *parent)
	return nil
}

func (g *Group) RevokeParent(parentId int) error {
	if parentId < 0 {
		return errors.New("parent id cannot be below 0")
	}

	index := -1
	for i := 0; i < len(g.Inheritance); i++ {
		if g.Inheritance[i].Id == parentId {
			index = i
		}
	}

	if index == -1 {
		return nil
	}

	// some fucking voodoo magic
	g.Permissions[index] = g.Permissions[len(g.Permissions)-1] // replace target with end element
	g.Permissions = g.Permissions[:len(g.Permissions)-1]       // replace with array -1 length

	return nil
}

func (g *Group) SetDefault(isDefault bool) error {
	g.IsDefault = isDefault

	return nil
}

func (g *Group) Grant(permission string) error {
	if permission == "" {
		return errors.New("permission cannot be empty")
	}

	alreadyHas := false
	for _, v := range g.Permissions {
		if v == permission {
			alreadyHas = true
		}
	}

	if !alreadyHas {
		g.Permissions = append(g.Permissions, permission)
	}

	return nil
}

func (g *Group) Can(permission string) (bool, error) {
	if g == nil || permission == "" {
		return false, errors.New("permission cannot be empty")
	}

	for _, v := range g.GetAllPermissions() {
		if v == permission {
			return true, nil
		}
	}

	return false, nil
}

func (g *Group) Revoke(permission string) error {
	if g == nil || permission == "" {
		return errors.New("permission cannot be empty")
	}

	index := -1
	for i := 0; i < len(g.Permissions); i++ {
		if g.Permissions[i] == permission {
			index = i
		}
	}

	if index == -1 {
		return nil
	}

	// some fucking voodoo magic
	g.Permissions[index] = g.Permissions[len(g.Permissions)-1] // replace target with end element
	g.Permissions = g.Permissions[:len(g.Permissions)-1]       // replace with array -1 length

	return nil
}

func (g *Group) GetAllPermissions() []string {
	all := []string{}
	all = append(all, g.Permissions...)
	for _, i := range g.Inheritance {
		getNestedPermissions(g, &i, &all)
	}

	return all
}

func getNestedPermissions(base *Group, g *Group, permissions *[]string) {
	if g == nil || base.Name == g.Name {
		return
	}

	*permissions = append(*permissions, g.Permissions...)

	if len(g.Inheritance) > 0 {
		for _, v := range g.Inheritance {
			go getNestedPermissions(g, &v, permissions)
		}
	}
}
