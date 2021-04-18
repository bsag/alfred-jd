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
	"log"
	"os"
	"os/exec"

	aw "github.com/deanishe/awgo"
	"github.com/deanishe/awgo/update"
)

// Name of the background job that checks for updates
const updateJobName = "checkForUpdate"

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
	doCheck                       bool

	// Our Workflow object
	wf *aw.Workflow
)

func init() {
	flag.BoolVar(&doCheck, "check", false, "check for a new version")

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

	// Alternate action: Get available releases from remote.
	if doCheck {
		wf.Configure(aw.TextErrors(true))
		log.Println("Checking for updates...")
		if err := wf.CheckForUpdate(); err != nil {
			wf.FatalError(err)
		}
		return
	}

	// Call self with "check" command if an update is due and a check
	// job isn't already running.
	if wf.UpdateCheckDue() && !wf.IsRunning(updateJobName) {
		log.Println("Running update check in background...")

		cmd := exec.Command(os.Args[0], "-check")
		if err := wf.RunInBackground(updateJobName, cmd); err != nil {
			log.Printf("Error starting update check: %s", err)
		}
	}

	// Only show update status if query is empty.
	if query == "" && wf.UpdateAvailable() {
		// Turn off UIDs to force this item to the top.
		// If UIDs are enabled, Alfred will apply its "knowledge"
		// to order the results based on your past usage.
		wf.Configure(aw.SuppressUIDs(true))

		// Notify user of update. As this item is invalid (Valid(false)),
		// actioning it expands the query to the Autocomplete value.
		// "workflow:update" triggers the updater Magic Action that
		// is automatically registered when you configure Workflow with
		// an Updater.
		//
		// If executed, the Magic Action downloads the latest version
		// of the workflow and asks Alfred to install it.
		wf.NewItem("Update available!").
			Subtitle("↩ to install").
			Autocomplete("workflow:update").
			Valid(false).
			Icon(updateAvailable)
	}

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
