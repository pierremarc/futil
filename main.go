/*
 *  Copyright (C) 2018 Pierre Marchand <pierre.m@atelier-cartographique.be>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE
 */

package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"go/build"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"text/template"
	"time"
)

var (
	funcMap = template.FuncMap{
		"first": func(s string) string {
			return strings.ToLower(string(s[0]))
		},
	}
)

type stringSlice []string

func (s *stringSlice) String() string {
	return strings.Join(*s, ",")
}

func (s *stringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

type typeMap map[string]string

var basicTypes = typeMap{
	"Bool":       "bool",
	"String":     "string",
	"Int":        "int",
	"Int8":       "int8",
	"Int16":      "int16",
	"Int32":      "int32",
	"Int64":      "int64",
	"UInt":       "uint",
	"UInt8":      "uint8",
	"UInt16":     "uint16",
	"UInt32":     "uint32",
	"UInt64":     "uint64",
	"UintPtr":    "uintptr",
	"Byte":       "byte",
	"Rune":       "rune",
	"Float32":    "float32",
	"Float64":    "float64",
	"Complex64":  "complex64",
	"Complex128": "complex128",
}

type generator struct {
	packageName string
	types       typeMap
	imports     []string
}

func (g *generator) generate(templatePath string) ([]byte, error) {
	bs, _ := ioutil.ReadFile(templatePath)
	t := template.Must(template.New("option").Parse(string(bs)))

	data := struct {
		PackageName string
		Timestamp   time.Time
		Types       typeMap
		Imports     []string
	}{
		g.packageName,
		time.Now().UTC(),
		g.types,
		g.imports,
	}

	var buf bytes.Buffer

	err := t.Execute(&buf, data)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("option: ")

	var imports stringSlice

	templateType := flag.String("type", "none", "type to generate [ func | option | result | array], (required)")
	basics := flag.Bool("basics", false, "generate for basic types")
	flag.Var(&imports, "import", "a package to import, can be repeated")
	outputName := flag.String("output", "", "output file name, default is <type>.go")

	flag.Parse()

	if "none" == *templateType {
		log.Fatal(errors.New("type argument is required"))
	}

	types := make(map[string]string)
	args := flag.Args()
	for _, pair := range args {
		parts := strings.Split(pair, "=")
		label := parts[0]
		typ := parts[1]
		types[label] = typ
	}

	if true == *basics {
		for k, v := range basicTypes {
			types[k] = v
		}
	}

	pkg, err := build.Default.ImportDir(".", 0)
	if err != nil {
		log.Fatal(err)
	}

	var (
		g generator
	)

	g.types = types
	g.packageName = pkg.Name
	g.imports = imports

	templatepath := path.Join(os.Getenv("GOPATH"),
		"src/github.com/pierremarc/futil", *templateType)

	src, err := g.generate(templatepath)
	if err != nil {
		log.Fatal(err)
	}

	outPath := fmt.Sprintf("%s.go", *templateType)
	if "" != *outputName {
		outPath = *outputName
	}

	if err = ioutil.WriteFile(outPath, src, 0644); err != nil {
		log.Fatalf("writing output: %s", err)
	}
}
