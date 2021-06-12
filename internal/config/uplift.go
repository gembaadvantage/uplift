/*
Copyright (c) 2021 Gemba Advantage

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package config

import (
	"io/ioutil"
	"os"

	"github.com/go-yaml/yaml"
)

// Uplift defines the root configuration of the application
type Uplift struct {
	FirstVersion string `yaml:"firstVersion"`
	Bumps        []Bump `yaml:"bumps"`
}

// Bump defines configuration for bumping indvidual files based
// on the new calculated semantic version number
type Bump struct {
	File  string `yaml:"file"`
	Regex string `yaml:"regex"`
	Count int    `yaml:"count"`
}

// Load the YAML config file
func Load(f string) (Uplift, error) {
	fh, err := os.Open(f)
	if err != nil {
		return Uplift{}, err
	}
	defer fh.Close()

	// Read the contents of the file in one go
	data, err := ioutil.ReadAll(fh)
	if err != nil {
		return Uplift{}, err
	}

	var cfg Uplift
	err = yaml.UnmarshalStrict(data, &cfg)
	return cfg, err
}
