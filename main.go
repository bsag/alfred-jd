// Copyright (c) 2018 Dean Jackson <deanishe@deanishe.net>
// MIT Licence - http://opensource.org/licenses/MIT

/*
Workflow fuzzy is a basic demonstration of AwGo's fuzzy filtering.

It displays and filters the contents of your Downloads directory in Alfred,
and allows you to open files, reveal in Finder or browse in Alfred.
*/
package main

import (
	"flag"

	aw "github.com/deanishe/awgo"
	"github.com/deanishe/awgo/update"
)

var (
	query string
	// Icons
	updateAvailable = &aw.Icon{Value: "icons/update-available.png"}
	areaIcon        = &aw.Icon{Value: "icons/area.png"}
	categoryIcon    = &aw.Icon{Value: "icons/category.png"}
	idIcon          = &aw.Icon{Value: "icons/id.png"}
	newIdIcon       = &aw.Icon{Value: "icons/newid.png"}

	repo = "bsag/alfred-jd"

	// command line arguments
	doAction, getLevel, catFolder string

	// Our Workflow object
	wf *aw.Workflow
)

func init() {
	// Initialise workflow and set up update mechanism
	wf = aw.New(update.GitHub(repo), aw.HelpURL(repo+"/issues"))

	// command line flags
	flag.StringVar(&doAction, "action", "", "choose action")
	flag.StringVar(&getLevel, "level", "", "choose level")
	flag.StringVar(&catFolder, "cat", "", "category folder")
}

// run executes the Script Filter.
func run() {
	wf.Args() // call to handle magic actions

	// ----------------------------------------------------------------
	// Parse command-line flags and decide what to do

	flag.Parse()

	// choose the entry point to the different workflow branches
	// search levels and act on candidates
	if doAction == "search" {
		doSearch()
	}

	// search categories and create new ID folder
	if doAction == "new" {
		makeNew()
	}

}

func main() {

	// Call workflow via `Run` wrapper to catch any errors, log them
	// and display an error message in Alfred.
	wf.Run(run)
}
