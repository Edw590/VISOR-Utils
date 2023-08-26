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
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"VISOR_S_L/PERSONAL_FILES_EOG"
)

// _PERSONAL_FOLDER is the name of the main directory of the personal files.
const _PERSONAL_FOLDER string = PersonalConsts.VISOR_DIR + "PERSONAL_FILES_EOG/"

/*
GPath (GoodPath) is sort of a copy of the string type but that represents a *surely* valid and correct path, also
according to the project conventions as described in the Path() function.

It's a "good path" because it's only given by Path(), which corrects the paths, and because the string component is
private to the package and only requested when absolutely necessary, like to communicate with Go's official functions
that require a string.
 */
type GPath struct {
	// s is the string that represents the path.
	s string
}

/*
PathFILESDIRS combines a path from the given subpaths of type string or GPath (ONLY), always ending a directory path
with the OS path separator.

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
	var ends_in_separator bool = false
	for i, sub_path := range sub_paths {
		val_str, ok := sub_path.(string)
		if ok {
			sub_paths_str = append(sub_paths_str, val_str)
			if i == len(sub_paths)-1 && (strings.HasSuffix(val_str, "/") || strings.HasSuffix(val_str, "\\")) {
				ends_in_separator = true
			}

			continue
		}

		val_GPath, ok := sub_path.(GPath)
		if ok {
			sub_paths_str = append(sub_paths_str, val_GPath.s)
			if i == len(sub_paths)-1 && (strings.HasSuffix(val_GPath.s, "/") || strings.HasSuffix(val_GPath.s, "\\")) {
				ends_in_separator = true
			}

			continue
		}

		// If it's not a string or GPath, it's an error.
		PanicGENERAL("pathFILESDIRS() received an invalid type of parameter. " + getVariableInfoGENERAL(sub_path))
	}

	// The call to Join() is on purpose - it correctly joins *and cleans* the final path string.
	var gPath GPath = GPath{filepath.Join(sub_paths_str...)}

	// Here can be os.PathSeparator and it will always be correct because the path is already cleaned.
	var path_separator string = string(os.PathSeparator)
	// Check if the path represents a directory and if it does, make sure the path separator is at the end.
	if gPath.Exists() {
		// Check if it's a directory through OS stats.
		if gPath.IsDir() && !strings.HasSuffix(gPath.s, path_separator) {
			gPath.s += path_separator
		}
	} else {
		// As last resort, check if it's a directory through the last character (project convention).
		if ends_in_separator && !strings.HasSuffix(gPath.s, path_separator) {
			gPath.s += path_separator
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
	return gPath.s
}

/*
ReadFile reads the contents of a file.

Note: all line breaks are replaced by "\n" for internal use, just like Python does.

-----------------------------------------------------------

– Returns:
  - the contents of the file or nil if an error occurs (including if the path describes a directory)
 */
func (gPath GPath) ReadFile() *string {
	if gPath.IsDir() {
		return nil
	}

	data, err := os.ReadFile(gPath.s)
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

Note: all line breaks are replaced by the OS line break(s). So for Windows, "\r" and "\n" are replaced by "\r\n" and for
any other, "\r\n" and "\r" are replaced by "\n".

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
WriteFile writes the contents of a file, creating it and any directories if necessary.

-----------------------------------------------------------

– Params:
  - content – the contents to write

– Returns:
  - nil if the file was written successfully, an error otherwise (including if the path describes a directory)
 */
func (gPath GPath) WriteFile(content []byte) error {
	if gPath.IsDir() || !gPath.CreateDir() {
		return nil
	}

	var err = os.WriteFile(gPath.s, content, 0o777)
	_ = os.Chmod(gPath.s, 0o777)

	return err
}

/*
IsDir checks if a path is a directory or a file, no matter if it exists and we have permissions to see it or not.

-----------------------------------------------------------

– Returns:
  - true if the path describes a directory, false if it describes a file
*/
func (gPath GPath) IsDir() bool {
	file_info, err := os.Stat(gPath.s)
	if nil == err {
		return file_info.IsDir()
	}

	return strings.HasSuffix(gPath.GPathToStringConversion(), string(os.PathSeparator))
}

/*
Exists checks if a path exists.

-----------------------------------------------------------

– Returns:
  - true if the path exists (meaning the program also has permissions to *see* the file), false otherwise
*/
func (gPath GPath) Exists() bool {
	_, err := os.Stat(gPath.s)

	return nil == err
}

/*
CreateDir creates a path and any necessary subdirectories in case they don't exist already, excluding the file if the
path represents a file.

-----------------------------------------------------------

– Returns:
  - true if the path was created successfully, false otherwise
 */
func (gPath GPath) CreateDir() bool {
	var path_list []string = strings.Split(gPath.s, string(os.PathSeparator))
	if !gPath.IsDir() {
		// If the path is a file, remove the file part of the path from the list.
		path_list = path_list[:len(path_list) - 1]
	}

	var current_path GPath = GPath{}
	if strings.HasPrefix(gPath.s, string(os.PathSeparator)) {
		current_path.s = string(os.PathSeparator)
	}
	for _, sub_path := range path_list {
		if "" == sub_path {
			continue
		}

		// Keep adding the subpaths until we reach the file part of the path.
		current_path.s += sub_path + string(os.PathSeparator)

		if !current_path.Exists() {
			if nil == os.Mkdir(current_path.s, 0o777) {
				_ = os.Chmod(current_path.s, 0o777)
			} else {
				return false
			}
		}
	}

	return true
}
