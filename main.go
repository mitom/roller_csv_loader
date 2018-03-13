// Copyright © 2018 Tamas Millian <tamas.millian@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package main

import (
	"encoding/csv"
	"fmt"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"roller/pkg"
)

type loader string

var Loader loader

func main() {}
func (l loader) Load(config *pkg.CacheConfig) []pkg.CachedProfile {
	//get the current user for the home directory to support tilde~ paths
	usr, err := user.Current()
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}

	path := config.GetOptions()["path"].(string)

	// handle home paths with tilde
	if path[:2] == "~/" {
		path = filepath.Join(usr.HomeDir, path[2:])
	}

	// run it through abspath to support ..paths
	path, _ = filepath.Abs(path)
	f, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}
	// read the file
	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}

	skipFirst := false
	_, exists := config.GetOptions()["skip_first"]
	if exists {
		skipFirst = config.GetOptions()["skip_first"].(bool)
	}
	_, exists = config.GetOptions()["mapping"]
	var mapping []string
	if !exists {
		mapping = []string{"name", "account_id", "role"}
	} else {
		var ok bool
		mapping, ok = convertStringSlice(config.GetOptions()["mapping"].([]interface{}))

		if !ok {
			fmt.Sprintf("%s is missing the `mapping` attribute", config.GetName())
			os.Exit(3)
		}
	}

	results := make([]pkg.CachedProfile, len(records))
	for i, row := range records {
		if i == 0  && skipFirst {
			continue
		}
		results[i] = parseRow(row, mapping)
	}

	return results
}

func parseRow(row []string, mapping []string) pkg.CachedProfile {
	var result pkg.CachedProfile
	for k, cell := range row {
		if k >= len(mapping) {
			break
		}

		switch mapping[k] {
		case "account_name":
			result.Name = strings.TrimSpace(cell)
			break
		case "role":
			result.Parameters.Role = strings.TrimSpace(cell)
			break
		case "account_id":
			result.Parameters.AccountId = strings.TrimSpace(cell)
			break
		case "switch_url":
			parsed, _ := url.Parse(cell)
			v, e := parsed.Query()["roleName"]
			if e && result.Parameters.Role == "" {
				result.Parameters.Role = v[0]
			}
			v, e = parsed.Query()["account"]
			if e && result.Parameters.AccountId == "" {
				result.Parameters.AccountId = v[0]
			}
			break
		default:
			continue

		}
	}
	return result
}


func convertStringSlice(data []interface{}) ([]string, bool) {
	r := make([]string, len(data))

	for k, v := range data {
		vs, ok := v.(string)
		if !ok {
			return nil, false
		}
		r[k] = vs
	}

	return r, true
}