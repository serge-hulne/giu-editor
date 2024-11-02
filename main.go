package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	g "github.com/AllenDang/giu"
	"github.com/sqweek/dialog"
)

//-----------------------------------------

//************
// App state :
//************

var text string
var filePath string
var AboutText string
var SizeX int = 800
var SizeY int = 600
var wnd *g.MasterWindow
var status string = "OK."

//-----------------------------------------

//**************
// Main
//**************

// Running the app :
func main() {
	// Load the menu data from the embedded JSON file
	loadMenu()

	// Create the window and start the GUI loop
	wnd = g.NewMasterWindow("Simple text editor", SizeX, SizeY, 0)
	wnd.Run(loop)
}

//-----------------------------------------

//*******************
// App visual layout :
//*******************

func Body() g.Widget {
	X, Y := wnd.GetSize()
	return g.Layout{
		g.Column(
			g.InputTextMultiline(&text).Size(float32(X-50), float32(Y-70)),
			g.Label("Status : "+status),
		),
	}
}

//-----------------------------------------

//##############
// App actions :
//##############

//

//**************
// Quit :
//**************

func quit() {
	os.Exit(0)
}

//**************
// Open File :
//**************

func OpenFile() {
	// Open a file dialog to select a file
	var err error
	filePath, err = dialog.File().Title("Open File").Load() // Set filePath globally
	if err != nil {
		log.Println("Failed to open file:", err)
		status = fmt.Sprintf("Failed to open file: %s", err)
		return
	}

	// Read the contents of the selected file
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Println("Failed to read file:", err)
		status = fmt.Sprintf("Failed to read file: %s", err)
		return
	}

	// Load the content into the text editor
	text = string(content)
	fmt.Println("File opened:", filePath)
	status = "File opened: " + filePath
	g.Update() // Refresh the UI with the new text
}

//**************
// Save File :
//**************

func SaveFile() {
	// Check if filePath is set
	if filePath == "" {
		log.Println("No file path specified. Use 'Save As' to specify a file name.")
		return
	}

	fmt.Println("Saving file:", filePath)
	f, err := os.Create(filePath)
	if err != nil {
		log.Println("Failed to create file:", err)
		status = fmt.Sprintf("Failed to create file: %s", err)
		return
	}

	_, err = f.WriteString(text)
	if err != nil {
		log.Println("Failed to write to file:", err)
		status = fmt.Sprintf("Failed to write to file: %s", err)
		f.Close()
		return
	}

	err = f.Close()
	if err != nil {
		log.Println("Failed to close file:", err)
		status = fmt.Sprintf("Failed to close file: %s", err)
	}
	fmt.Println("File saved successfully!")
	status = "File saved successfully!"
}

//**************
// Save File As :
//**************

func SaveFileAs() {
	// Open a "Save As" dialog to select a new file name
	newFilePath, err := dialog.File().Title("Save File As").Save()
	if err != nil {
		log.Println("Failed to select file path:", err)
		return
	}

	// Save the content to the new file
	fmt.Println("Saving file as:", newFilePath)
	f, err := os.Create(newFilePath)
	if err != nil {
		log.Println("Failed to create file:", err)
		return
	}

	_, err = f.WriteString(text)
	if err != nil {
		log.Println("Failed to write to file:", err)
		f.Close()
		return
	}

	err = f.Close()
	if err != nil {
		log.Println("Failed to close file:", err)
		return
	}

	// Update the global filePath to the new file path
	filePath = newFilePath
	fmt.Println("File saved successfully as:", filePath)
	status = "File" + filePath + " saved successfully!"
}

//**************
// About :
//**************

func About() {
	// Set the flag to show the "About" popup
	AboutText = "This is a simple text editor.\nVersion 1.0\n(c) 2024 Serge Hulne."
	showAboutPopup = true
}
