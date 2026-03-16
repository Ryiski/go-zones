package main

import (
	"os"
	"path/filepath"
	"strings"
)

// FileModeApplication will run the File Mode application logic bootstrap
func (config CLIOpts) FileModeApplication() {
	logStdOut("Source file: " + config.Source)
	logStdOut("Target directory: " + config.Dir)

	absoluteTargetDirectory, err := filepath.Abs(config.Dir)
	check(err)

	// Initialize the Target Directory
	if err := ValidateConfigDirectory(absoluteTargetDirectory); err != nil {
		logStdErr("Target directory unwritable!")
	} else {
		directoryCheck, err := DirectoryExists(absoluteTargetDirectory)
		check(err)
		if !directoryCheck {
			//Directory doesn't exist, create
			CreateDirectory(absoluteTargetDirectory)
		}
	}

	configPath := absoluteTargetDirectory + "/config"
	CreateDirectory(configPath)
	CreateDirectory(absoluteTargetDirectory + "/zones")

	// Read in Zones file
	server, err := NewDNSServer(config)
	check(err)

	//zones := server.Zones

	revViewPair, err := GenerateBindReverseZoneFiles(&server.DNS, absoluteTargetDirectory)
	check(err)

	_, err = GenerateBindConfig(&server.DNS, absoluteTargetDirectory, revViewPair)
	check(err)

	_, err = GenerateBindZoneConfigFile(&server.DNS, absoluteTargetDirectory)
	check(err)

	_, err = GenerateBindZoneFiles(&server.DNS, absoluteTargetDirectory)
	check(err)

	files, err := os.ReadDir(configPath)
	check(err)

	var includes strings.Builder

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fullPath := filepath.Join(configPath, file.Name())
		includes.WriteString("include \"" + fullPath + "\";\n")
	}

	includeFile := filepath.Join(absoluteTargetDirectory, "includes.conf")

	err = os.WriteFile(includeFile, []byte(includes.String()), 0644)
	check(err)

	//_, err = LoopThroughZonesForBindConfig(server, absoluteTargetDirectory)
	//check(err)

	//_, err = LoopThroughZonesForBindZonesFiles(&server.DNS.Zones, absoluteTargetDirectory)
	//check(err)

	//_, err = LoopThroughZonesForBindReverseV4ZonesFiles(zones, absoluteTargetDirectory)
	//check(err)
}
