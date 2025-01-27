package main

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func Test_model_Init(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		want tea.Cmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: construct the receiver type.
			var m model
			got := m.Init()
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("Init() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_model_Init(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		want tea.Cmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: construct the receiver type.
			var m model
			got := m.Init()
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("Init() = %v, want %v", got, tt.want)
			}
		})
	}
}
