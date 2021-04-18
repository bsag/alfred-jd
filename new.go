package main

import (
	"flag"
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
)


func parseCatOrID(folderName string, cat bool) (num int64) {
	// if cat is true, we want to get only the AC component
	// TODO: return an error if there's no match
	if cat {
		r := regexp.MustCompile(`^([0-9]{2})`)
		matches := r.FindAllStringSubmatch(folderName, -1)

		for _, v := range matches {
			num, _ := strconv.ParseInt(v[1], 10, 64)
			return num
		}
	} else {
		// if cat is false, we want only the ID component
		// to which we add 1.
		r := regexp.MustCompile(`^([0-9]{2})\.([0-9]{2})`)
		matches := r.FindAllStringSubmatch(folderName, -1)

		for _, v := range matches {
			num, _ := strconv.ParseInt(v[2], 10, 64)
			return num + 1
		}
	}
	return num
}

func getNextIdx(files []file, folderName string) (jdx string) {

	ac := parseCatOrID(folderName, true)

	// We don't have any ID files yet in this category
	// so we construct an AC.ID from the enclosing category folderName
	// and add an ID of 01.
	if len(files) == 0 {
		jid := fmt.Sprintf("%02d.%02d ", ac, 1)
		return jid
	} else {
	last := filepath.Base(files[len(files)-1].Path) // get the path of the last file, numerically sorted
		lastid := parseCatOrID(last, false) // find ID component only

		if lastid >= 99 {
			wf.Warn("You already have 99 IDs", "You may want to split up this category.")
		}

		jid := fmt.Sprintf("%02d.%02d ", ac, lastid)
		return jid
	}

}

func labelFolder(catFolder, query string) {

	// get the foldername without path for display
	folderName := filepath.Base(catFolder)

	// get a slice of files in the category folder
	files := readDir(catFolder)
	// get the next index number: returns an AC.ID string
	idx := getNextIdx(files, folderName)
	newFolderName := idx + query

	wf.NewItem("Enter name for item folder within "+folderName).
		Subtitle("New folder name: " + newFolderName).
		Arg(filepath.Join(catFolder, newFolderName)).
		Valid(true).
		Var("fpath", catFolder).
		Var("dirname", newFolderName).
		Icon(newIdIcon)

	wf.SendFeedback()
}

// Create a new ID item in a category folder
func makeNew() {
	wf.Args()

	query := flag.Arg(0)

	if catFolder != "" {
		labelFolder(catFolder, query)
		return
	}

	startDir := setup()

	// ----------------------------------------------------------------
	// Load data and create Alfred items
	// ----------------------------------------------------------------

	for _, file := range readDir(startDir + "/*") {
		// NewFileItem usually passes the file path as Arg()
		// but this prepopulates the text entry at the next step
		// which is not what we want. Passing Arg("") here clears
		// the field and we put file.Path in Var() instead.
		wf.NewFileItem(file.Path).
			Arg("").
			Subtitle("Select Category in which to create new ID folder.").
			Var("cat", file.Path).
			Icon(categoryIcon)
	}

	if query != "" {
		wf.Filter(query)
	}

	wf.WarnEmpty("No matching folders found", "Try a different query?")

	wf.SendFeedback()
}
