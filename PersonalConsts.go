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
	"strings"
)

var PersonalConsts_GL PersonalConsts = PersonalConsts{}

// personalConstsEOG is the internal struct with the format of the PersonalConsts_EOG.json file.
type personalConstsEOG struct {
	VISOR_DIR string

	VISOR_EMAIL_ADDR string
	VISOR_EMAIL_PW string

	USER_EMAIL_ADDR string
}

// PersonalConsts is a struct containing the constants that are personal to the user.
type PersonalConsts struct {
	// _VISOR_DIR is the full path to the main directory of VISOR.
	_VISOR_DIR GPath

	// _VISOR_EMAIL_ADDR is VISOR's email address
	_VISOR_EMAIL_ADDR string
	// _VISOR_EMAIL_PW is VISOR's email password
	_VISOR_EMAIL_PW string

	// USER_EMAIL_ADDR is the email address of the user, used for all email communication
	USER_EMAIL_ADDR string
}

/*
init is the function that initializes the global variables of the PersonalConsts struct.
 */
func (personalConsts *PersonalConsts) init() error {
	const PERSONAL_CONSTS_FILE string = "PersonalConsts_EOG.json"

	bytes, err := os.ReadFile(PERSONAL_CONSTS_FILE)
	if err != nil {
		cwd, err := os.Getwd()
		if err != nil {
			cwd = "[ERROR]"
		}
		return errors.New("No " + PERSONAL_CONSTS_FILE + " file found in the current working directory: \"" + cwd + "\"! Aborting...")
	}

	var struct_file_format personalConstsEOG
	if !FromJsonGENERAL(bytes, &struct_file_format) {
		return errors.New("file " + PERSONAL_CONSTS_FILE + " file is corrupted! Aborting...")
	}

	// Set the global variables

	// Ending _VISOR_DIR with a slash to be sure it's there (as it should - it's a directory)
	personalConsts._VISOR_DIR = PathFILESDIRS("", struct_file_format.VISOR_DIR + "/")

	personalConsts._VISOR_EMAIL_ADDR = struct_file_format.VISOR_EMAIL_ADDR
	personalConsts._VISOR_EMAIL_PW = struct_file_format.VISOR_EMAIL_PW

	personalConsts.USER_EMAIL_ADDR = struct_file_format.USER_EMAIL_ADDR

	if !strings.Contains(personalConsts._VISOR_EMAIL_ADDR, "@") || personalConsts._VISOR_EMAIL_PW == "" ||
				!strings.Contains(personalConsts.USER_EMAIL_ADDR, "@") {
		return errors.New("Some fields in " + PERSONAL_CONSTS_FILE + " are empty! Aborting...")
	}

	var visor_path GPath = personalConsts._VISOR_DIR
	if !visor_path.Exists() {
		return errors.New("The VISOR directory \"" + visor_path.GPathToStringConversion() + "\" does not exist! Aborting...")
	}

	return nil
}
