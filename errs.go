package main

import "errors"

var levelComplete = errors.New("level complete")
var levelFailed = errors.New("level failed")

type SelectedLevel struct {
	name string
}

func NewLevelSelected(name string) *SelectedLevel {
	return &SelectedLevel{name: name}
}

func (ls *SelectedLevel) Error() string {
	return ls.name
}
