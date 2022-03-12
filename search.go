package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"regexp"

	aw "github.com/deanishe/awgo"
)

type file struct {
	Path  string
	IsDir bool
}

func expandPath(path string) string {
	re := regexp.MustCompile(`^~`)
	return re.ReplaceAllString(path, os.Getenv("HOME"))
}

// setup gets the JD_DIR workflow variable
//   Throws an error and stops if the variable has no value or the path
//   does not exist.
func setup() (startDir string) {

	// get the value from Alfred's workflow environment variables
	startDir = wf.Config.Get("JD_DIR")
	startDir = expandPath(startDir)

	// check if the variable has a value set
	if startDir == "" {
		wf.Fatal("No path set for J.D directory.")
	}

	// check whether the path exists if it has a value set
	if _, err := os.Stat(startDir); os.IsNotExist(err) {
		wf.Fatal("Path set for J.D directory does not exist.")
	}

	return (startDir)
}

// get the folders matching the glob pattern
func readDir(pattern string) (files []file) {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}

	for _, m := range matches {
		infos, _ := ioutil.ReadDir(m)
		for _, fi := range infos {
			// ignore hidden files
			if strings.HasPrefix(fi.Name(), ".") {
				continue
			}

			// skip non-directory items
			// TODO: make this a parameter so we can include it
			// when calling from new item and not otherwise
			if !(fi.IsDir()) {
				continue
			}

			files = append(files, file{filepath.Join(m, fi.Name()), fi.IsDir()})
		}
	}

	return files
}

func runSearch(searchlvl, query string) {

	// ----------------------------------------------------------------
	// Get path to J.D directory and construct glob pattern
	// ----------------------------------------------------------------

	startDir := setup()

	var searchPattern string

	switch searchlvl {
	case "A":
		searchPattern = startDir // No glob needed - just search the J.D root
	case "C":
		searchPattern = startDir + "/*"
	case "ID":
		searchPattern = startDir + "/*/*"
	default:
		searchPattern = ""
	}

	// ----------------------------------------------------------------
	// Load data and create Alfred items
	// ----------------------------------------------------------------

	for _, file := range readDir(searchPattern) {
		// Convenience method. Sets Item title to filename, subtitle
		// to shortened path, arg to full path, and icon to file icon.
		it := wf.NewFileItem(file.Path)

		// Alternate actions
		// Default is to open in default application
		it.NewModifier(aw.ModCtrl).
			Subtitle("Open in Terminal").
			Var("action", "terminal")

		if file.IsDir {
			it.NewModifier(aw.ModAlt).
				Subtitle("Browse in Alfred").
				Var("action", "browse")
		}
	}

	if query != "" {
		wf.Filter(query)
	}

	// Show a warning in Alfred if there are no items
	wf.WarnEmpty("No matching folders found", "Try a different query?")
	wf.SendFeedback()
}

// doSearch does the main work of showing the user options for level to be searched
// then passes on the user query and finally filters and sends the results back to Alfred.
func doSearch() {
	// ----------------------------------------------------------------
	// Handle CLI arguments
	// ----------------------------------------------------------------

	// You should always use wf.Args() in Script Filters. It contains the
	// same as os.Args[1:], but the arguments are first parsed for AwGo's
	// magic actions (i.e. "workflow:*" to allow the user to easily open
	// the log or data/cache directory).
	wf.Args()

	query := flag.Arg(0)

	// first time around, all these vars are empty,
	// so we just show the default list of options
	// then run this function when the option
	// is chosen.
	if getLevel != "" {
		runSearch(getLevel, query)
		return
	}

	log.Printf("first query=%s", query)

	// Create search options as Alfred items
	wf.NewItem("Areas").
		Subtitle("Search Areas").
		Valid(true).
		Var("lvl", "A").
		Icon(areaIcon)

	wf.NewItem("Categories").
		Subtitle("Search Categories").
		Valid(true).
		Var("lvl", "C").
		Icon(categoryIcon)

	wf.NewItem("IDs (Items)").
		Subtitle("Search Items").
		Valid(true).
		Var("lvl", "ID").
		Icon(idIcon)

	// ----------------------------------------------------------------
	// Filter items based on user query
	// ----------------------------------------------------------------

	if query != "" {
		wf.Filter(query)
	}

	// ----------------------------------------------------------------
	// Send results to Alfred
	// ----------------------------------------------------------------

	// Show a warning in Alfred if there are no items
	wf.WarnEmpty("No matching folders found", "Try a different query?")

	// Send JSON to Alfred. After calling this function, you can't send
	// any more results to Alfred.
	wf.SendFeedback()
}

