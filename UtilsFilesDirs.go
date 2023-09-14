/*******************************************************************************
 * Copyright 2023-2023 Edw590
 *
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 ******************************************************************************/

package Utils

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

/*
GPath (GoodPath) is sort of a copy of the string type but that represents a *surely* valid and correct path, also
according to the project conventions as described in the Path() function.

It's a "good path" because it's only given by Path(), which corrects the paths, and because the string component is
private to the package and only requested when absolutely necessary, like to communicate with Go's official functions
that require a string.
*/
type GPath struct {
	// p is the string that represents the path.
	p string
	// s just maps to the OS path separator.
	s string
}

/*
PathFILESDIRS combines a path from the given subpaths of type string or GPath (ONLY), always ending a directory path
with the OS path separator.

Note: the final path separator is the first one found in the subpaths, or the OS path separator if none is found.

To end a path to a directory, ALWAYS end with a path separator (important project convention) and is the only
way to know if the path is a directory or not in case: there are no permissions to access it; it doesn't exist;
it's a relative path - various functions here depend on this convention in these listed cases!

-----------------------------------------------------------

– Params:
  - sub_paths – the subpaths to combine

– Returns:
  - the final path as a GPath
*/
func PathFILESDIRS(sub_paths ...any) GPath {
	var sub_paths_str []string = nil
	for _, sub_path := range sub_paths {
		val_str, ok := sub_path.(string)
		if ok {
			sub_paths_str = append(sub_paths_str, val_str)

			continue
		}

		val_GPath, ok := sub_path.(GPath)
		if ok {
			sub_paths_str = append(sub_paths_str, val_GPath.p)

			continue
		}

		// If it's not a string or GPath, it's an error.
		PanicGENERAL("pathFILESDIRS() received an invalid type of parameter. " + getVariableInfoGENERAL(sub_path))
	}

	if len(sub_paths_str) == 0 {
		return GPath{}
	}

	// Replace all the path separators with the OS path separator.
	for i, sub_path := range sub_paths_str {
		// Replace all the path separators with the OS path separator for Join() to work (only works with the OS one).
		sub_path = strings.Replace(sub_path, "/", string(os.PathSeparator), -1)
		sub_path = strings.Replace(sub_path, "\\", string(os.PathSeparator), -1)

		sub_paths_str[i] = sub_path
	}

	// Check if the last subpath ends in a path separator before calling Join() which will remove it if it's there.
	var ends_in_separator bool = false
	if strings.HasSuffix(sub_paths_str[len(sub_paths_str)-1], string(os.PathSeparator)) {
		ends_in_separator = true
	}

	// The call to Join() is on purpose - it correctly joins *and cleans* the final path string.
	var gPath GPath = GPath{
		p:   filepath.Join(sub_paths_str...),
		s:   string(os.PathSeparator),
	}

	// Check if the path represents a directory and if it does, make sure the path separator is at the end (especially
	// since Join() removes it if it's there).
	if gPath.Exists() {
		if gPath.DescribesDir() && !strings.HasSuffix(gPath.p, gPath.s) {
			gPath.p += gPath.s
		}
	} else {
		// As last resort, check if it's a directory through the last character (project convention).
		if ends_in_separator && !strings.HasSuffix(gPath.p, gPath.s) {
			gPath.p += gPath.s
		}
	}

	return gPath
}

/*
Add adds subpaths to a path.

-----------------------------------------------------------

– Params:
  - sub_paths – the subpaths to add

– Returns:
  - the final path as a GPath
 */
func (gPath GPath) Add(sub_paths ...any) GPath {
	// Create a temporary slice with the first element + the subpaths, all in a 1D slice and all as the 1st parameter
	// of the Path function.
	var temp []any = append([]any{gPath}, sub_paths...)

	return PathFILESDIRS(temp...)
}

/*
GPathToStringConversion converts a GPath to a string.

The function has a big name on purpose to discourage its use - only use it when absolutely necessary (for the reason
written in the GPath type), like to call Go official file/directory functions.

-----------------------------------------------------------

– Params:
  - path – the path to convert

– Returns:
  - the GPath
*/
func (gPath GPath) GPathToStringConversion() string {
	return gPath.p
}

/*
ReadFile reads the contents of a file.

Note: all line breaks are replaced by "\n" for internal use, just like Python does.

-----------------------------------------------------------

– Returns:
  - the contents of the file or nil if an error occurs (including if the path describes a directory)
 */
func (gPath GPath) ReadFile() *string {
	if gPath.DescribesDir() || !gPath.Exists() {
		return nil
	}

	data, err := os.ReadFile(gPath.p)
	if nil != err {
		return nil
	}
	var ret string = string(data)

	ret = strings.ReplaceAll(ret, "\r\n", "\n")
	ret = strings.ReplaceAll(ret, "\r", "\n")

	return &ret
}

/*
WriteTextFile writes the contents of a text file, creating it and any directories if necessary.

Note: all line breaks are replaced with the OS line break(s). So for Windows, "\r" and "\n" are replaced with "\r\n" and
for any other, "\r\n" and "\r" are replaced by "\n".

-----------------------------------------------------------

– Params:
  - content – the contents to write

– Returns:
  - nil if the file was written successfully, an error otherwise (including if the path describes a directory)
*/
func (gPath GPath) WriteTextFile(content string) error {
	var new_content string = content
	if "windows" == runtime.GOOS {
		new_content = strings.ReplaceAll(new_content, "\r\n", "\n")
		new_content = strings.ReplaceAll(new_content, "\r", "\n")
		new_content = strings.ReplaceAll(new_content, "\n", "\r\n")
	} else {
		new_content = strings.ReplaceAll(new_content, "\r\n", "\n")
		new_content = strings.ReplaceAll(new_content, "\r", "\n")
	}

	return gPath.WriteFile([]byte(new_content))
}

/*
WriteFile writes the raw contents of a file, creating it and any directories if necessary.

-----------------------------------------------------------

– Params:
  - content – the contents to write

– Returns:
  - nil if the file was written successfully, an error otherwise (including if the path describes a directory)
 */
func (gPath GPath) WriteFile(content []byte) error {
	if gPath.DescribesDir() || nil != gPath.Create(true) {
		return nil
	}

	var err = os.WriteFile(gPath.p, content, 0o777)
	_ = os.Chmod(gPath.p, 0o777)

	return err
}

/*
DescribesDir checks if a path DESCRIBES a directory or a file - means no matter if it exists or we have permissions to
see it or not.

It first checks if the path exists and if it does, checks if it's a directory or not - else it resorts to the path
string only, using the project convention in which a path that ends in a path separator is a directory and one that
doesn't is a file.

-----------------------------------------------------------

– Returns:
  - true if the path describes a directory, false if it describes a file
 */
func (gPath GPath) DescribesDir() bool {
	file_info, err := os.Stat(gPath.p)
	if nil == err {
		return file_info.IsDir()
	}

	return strings.HasSuffix(gPath.GPathToStringConversion(), gPath.s)
}

/*
Exists checks if a path exists.

-----------------------------------------------------------

– Returns:
  - true if the path exists (meaning the program also has permissions to *see* the file), false otherwise
*/
func (gPath GPath) Exists() bool {
	if nil != gPath.IsSupported() {
		return false
	}

	_, err := os.Stat(gPath.p)
	return nil == err
}

/*
Create creates a path and any necessary subdirectories in case they don't exist already.

-----------------------------------------------------------

– Params:
  - also_file – if true, creates the file *too* if the path represents a file

– Returns:
  - nil if the path was created successfully, an error otherwise
*/
func (gPath GPath) Create(create_file bool) error {
	if err := gPath.IsSupported(); nil != err {
		return err
	}

	var path_list []string = strings.Split(gPath.p, gPath.s)
	var describes_file bool = false
	if !gPath.DescribesDir() {
		// If the path is a file, remove the file part of the file from the list so that it describes a directory only,
		// but memorize if it describes a file.
		describes_file = true
		path_list = path_list[:len(path_list) - 1]
	}

	if !PathFILESDIRS(gPath.p[ : FindAllIndexesGENERAL(gPath.p, gPath.s)[len(path_list) - 1] + 1]).Exists() {
		var current_path GPath = GPath{}
		if strings.HasPrefix(gPath.p, gPath.s) {
			current_path.p = gPath.s
		}
		for _, sub_path := range path_list {
			if "" == sub_path {
				continue
			}

			// Keep adding the subpaths until we reach the file part of the path, where the loop stops.
			current_path.p += sub_path + gPath.s

			if !current_path.Exists() {
				if err := os.Mkdir(current_path.p, 0o777); nil == err {
					_ = os.Chmod(current_path.p, 0o777)
				} else {
					return err
				}
			}
		}
	}

	// Create the file if the path represents a file.
	if create_file && describes_file && !gPath.Exists() {
		_, err := os.Create(gPath.p)
		_ = os.Chmod(gPath.p, 0o777)

		if nil != err {
			return err
		}
	}

	return nil
}

/*
Remove removes a file or directory.

-----------------------------------------------------------

– Returns:
  - nil if the file or directory was removed successfully, an error otherwise
 */
func (gPath GPath) Remove() error {
	if err := gPath.IsSupported(); nil != err {
		return err
	}

	return os.Remove(gPath.p)
}

/*
IsSupported checks if the path is supported by the current OS.

-----------------------------------------------------------

– Returns:
  - nil if the path is supported by the current OS, false otherwise
*/
func (gPath GPath) IsSupported() error {
	// If the path is relative, it works everywhere (it's not specific to any OS). If it's absolute, it's supported
	// if it's an absolute path for the current OS.

	// Note: can't check with filepath.IsAbs() because it returns false for paths that are not supported by the current
	// OS, but are absolute for another OS. So the check must be made manually.

	// Check if the path is the wrong absolute type for the current OS. Else it's supported.
	// Don't forget the separators are changed to the current OS ones, so the checks are "inverted".
	if "windows" == runtime.GOOS {
		// Then check if it's a Linux absolute path.
		if strings.HasPrefix(gPath.p, "\\") {
			return errors.New("the path is not supported by the current OS")
		}
	} else {
		// Then check if it's a Windows absolute path.
		if len(gPath.p) >= 2 && (((gPath.p[0] >= 'a' && gPath.p[0] <= 'z' || gPath.p[0] >= 'A' && gPath.p[0] <= 'Z') && gPath.p[1] == ':') ||
					strings.HasPrefix(gPath.p, "//")) {
			return errors.New("the path is not supported by the current OS")
		}
	}
	// Else it's relative or absolute for the current OS.

	return nil
}
