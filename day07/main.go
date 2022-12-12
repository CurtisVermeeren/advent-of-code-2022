package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

/*

	"cd x" changes directory to x

	"cd .." moves out one level from the current directory

	"cd /" switches the current directory to the outermost

	"ls" means list. It prints out all of the files and directories.
	files are listed as "filesize filename"
	directories are listed as "dir xyz"

*/

// node represents either a file or directory
type node struct {
	name     string           // name of the file or directory
	size     int              // size of a file
	file     bool             // true if the node is a file
	children map[string]*node // all children directories and files
	parent   *node            // the parent node (one level up directory)
}

func partOne() (int, error) {

	file, err := os.Open("./input.txt")
	defer file.Close()
	if err != nil {
		return -1, err
	}

	fileScanner := bufio.NewScanner(file)

	// Build the filesystem and track the current directory
	fileSystem := []*node{}
	var currentDir *node

	for fileScanner.Scan() {
		command := strings.Fields(fileScanner.Text())
		// $ cd dir
		if len(command) > 2 {
			if command[2] == ".." {
				// Move up one dir to parent
				currentDir = currentDir.parent
			} else if command[2] == "/" {
				// Create the root directory /
				currentDir = &node{"/", 0, false, make(map[string]*node), nil}
				fileSystem = append(fileSystem, currentDir)
			} else {
				// Move to the child directory specified in the command
				currentDir = currentDir.children[command[2]]
			}
		} else if command[0] == "dir" {
			// Add the directory from the command
			currentDir.children[command[1]] = &node{command[1], 0, false, make(map[string]*node), currentDir}
			fileSystem = append(fileSystem, currentDir.children[command[1]])
		} else if command[0] != "$" {
			// Add the file from the command
			fileSize, _ := strconv.Atoi(command[0])
			currentDir.children[command[1]] = &node{command[1], fileSize, true, nil, currentDir}
		}
	}

	// Count the combined size of all directories with size < 100000
	var combined int

	for _, dir := range fileSystem {
		size := calculateDirectorySize(*dir)
		if size <= 100000 {
			combined += size
		}
	}

	return combined, nil
}

func partTwo() (int, error) {

	file, err := os.Open("./input.txt")
	defer file.Close()
	if err != nil {
		return -1, err
	}

	fileScanner := bufio.NewScanner(file)

	// Build the filesystem and track the current directory
	fileSystem := []*node{}
	var currentDir *node

	for fileScanner.Scan() {
		command := strings.Fields(fileScanner.Text())
		// $ cd dir
		if len(command) > 2 {
			if command[2] == ".." {
				// Move up one dir to parent
				currentDir = currentDir.parent
			} else if command[2] == "/" {
				// Create the root directory /
				currentDir = &node{"/", 0, false, make(map[string]*node), nil}
				fileSystem = append(fileSystem, currentDir)
			} else {
				// Move to the child directory specified in the command
				currentDir = currentDir.children[command[2]]
			}
		} else if command[0] == "dir" {
			// Add the directory from the command
			currentDir.children[command[1]] = &node{command[1], 0, false, make(map[string]*node), currentDir}
			fileSystem = append(fileSystem, currentDir.children[command[1]])
		} else if command[0] != "$" {
			// Add the file from the command
			fileSize, _ := strconv.Atoi(command[0])
			currentDir.children[command[1]] = &node{command[1], fileSize, true, nil, currentDir}
		}
	}

	// Total space used is 70000000 subtract the size of the root directory
	// 30000000 is needed so find the difference from the free space had and the desired free space amount
	freeSpace := 30000000 - (70000000 - calculateDirectorySize(*fileSystem[0]))

	var sizeToRemove int = calculateDirectorySize(*fileSystem[0])

	for _, dir := range fileSystem {
		size := calculateDirectorySize(*dir)
		// If the size of the current directory is large enough to create the free space needed
		// AND the size of the current directory is smaller than the previous directory being removed
		if size > freeSpace && size-freeSpace < sizeToRemove-freeSpace {
			sizeToRemove = size

		}
	}
	return sizeToRemove, nil
}

func calculateDirectorySize(rootNode node) (size int) {
	// If the rootNode is a file just return the size
	if rootNode.file {
		return rootNode.size
	}

	// IF the rootNode is a directory calculate the filesize of its children dirs
	for _, dir := range rootNode.children {
		size += calculateDirectorySize(*dir)
	}
	return
}

func main() {
	size, err := partOne()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(size)

	size, err = partTwo()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(size)
}
