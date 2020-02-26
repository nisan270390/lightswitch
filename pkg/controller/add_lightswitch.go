package controller

import (
	"LightSwitch/pkg/controller/lightswitch"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, lightswitch.Add)
}
