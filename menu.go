package main

import (
	"encoding/json"
	"log"

	_ "embed"

	g "github.com/AllenDang/giu"
)

// Define structs to represent the menu and submenu structure
type SubmenuItem struct {
	Label  string `json:"label"`
	Action string `json:"action"`
}

type MenuItem struct {
	Label    string        `json:"label"`
	Submenus []SubmenuItem `json:"submenus"`
}

type MenuData struct {
	Menus []MenuItem `json:"menus"`
}

// Actions for submenu items
func submenuAction(actionName string) func() {
	return func() {
		ActionMap[actionName]()
		menu_choice = 0
		g.Update()
	}
}

var menuData MenuData
var menu_choice = 0
var showAboutPopup bool // State variable to control the "About" popup

//go:embed menu.json
var menuJSON []byte

func loadMenu() {
	// Parse the embedded JSON file
	err := json.Unmarshal(menuJSON, &menuData)
	if err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}
}

func loop() {
	// Create a slice for main menu buttons
	var mainMenuItems []g.Widget
	for idx, menu := range menuData.Menus {
		// Capture index and label for each main menu button
		i := idx
		label := menu.Label

		// Append an asterisk to the label of the selected menu item
		if menu_choice == i+1 {
			label = "<" + label + ">"
		}

		// Create the button for the menu item
		mainMenuItems = append(mainMenuItems, g.Button(label).OnClick(func() {
			menu_choice = i + 1
			g.Update() // Force redraw of the UI when menu changes
		}))
	}

	// Layout for the main menus (horizontal)
	g.SingleWindow().Layout(
		g.Row(mainMenuItems...),
		g.Separator(),
		g.Row(
			g.Column(
				// Layout for submenus (vertical)
				func() g.Widget {
					if menu_choice > 0 && menu_choice <= len(menuData.Menus) {
						submenu := menuData.Menus[menu_choice-1].Submenus
						var submenuWidgets []g.Widget
						for _, sub := range submenu {
							// Add each submenu button and its associated action
							submenuWidgets = append(submenuWidgets, g.Button(sub.Label).OnClick(submenuAction(sub.Action)))
						}
						return g.Column(submenuWidgets...)
					}
					return g.Label("")
				}()),
			g.Column(Body())),
	)

	// Check if we should open the "About" popup
	if showAboutPopup {
		showAboutPopup = false // Reset the flag
		g.OpenPopup("About")
		g.Update()
	}

	// Define the "About" popup modal
	g.PopupModal("About").Layout(
		g.Label(AboutText),
		g.Button("Close").OnClick(func() {
			g.CloseCurrentPopup() // Close the popup when "Close" is clicked
		}),
	).Build() // Ensure Build is called on the PopupModal
}
