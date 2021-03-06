// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"flag"
	"fmt"
	"os"

	"io"

	log "github.com/Sirupsen/logrus"
	"github.com/natdm/typewriter/parse"
	"github.com/natdm/typewriter/template"
)

func main() {
	inFlag := flag.String("dir", "./", "dir is to specify what folder to parse types from")
	fileFlag := flag.String("file", "", "file is to parse a single file. Will override a directory")
	langFlag := flag.String("lang", "", "determine the language. One of 'flow', 'ts")
	outFlag := flag.String("out", "", "file and path to save output to")
	vFlag := flag.Bool("v", false, "verbose logging")
	recursiveFlag := flag.Bool("r", true, "to recursively ascend all folders in dir")
	flag.Usage = usage
	flag.Parse()

	var lang template.Language
	switch *langFlag {
	case "flow":
		lang = template.Flow
	case "elm":
		lang = template.Elm
	case "ts":
		lang = template.Typescript
	default:
		log.Fatalln("Please pick a proper language ['elm', 'flow', 'ts']")
	}

	var out io.Writer
	if *outFlag != "" {
		f, err := os.Create(*outFlag)
		if err != nil {
			log.Fatalln(err)
		}
		defer f.Close()
		out = f
	} else {
		out = os.Stdout
	}

	var (
		files []string
		types map[string]*template.PackageType
		err   error
	)

	if *fileFlag != "" {
		types, err = parse.Files([]string{*fileFlag}, *vFlag)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		if err = parse.Directory(*inFlag, *recursiveFlag, &files, *vFlag); err != nil {
			log.Fatalln(err)
		}
		types, err = parse.Files(files, *vFlag)
		if err != nil {
			log.Fatalln(err)
		}
	}

	if err := template.Draw(types, out, lang, *vFlag); err != nil {
		log.Fatalln(err)
	}
	log.Info("Done")
}

func usage() {
	fmt.Print(`
	Typewriter
	Convert Go types to other languages
	
	Flags:
		-dir	Parse a complete directory 
			example: 	-dir= ../src/appname/models/
			default: 	./

		-file	Parse a single go file 
			example: 	-file= ../src/appname/models/app.go
			overrides 	-dir and -recursive

		-out	Saves content to folder
			example: 	-out= ../src/appname/models/
						-out= ../src/appname/models/customname.js
			default: 	./models. 

		-r		Transcends directories
			example:	-recursive= false
			default:	true

		-v		Verbose logging, detailing every skipped type, file, or field.
			default: 	false

		-lang 	Language to parse to. One of ["flow"]
			example:	-lang flow
			default:	will not parse
`)
}
