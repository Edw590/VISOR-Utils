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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"
	"unsafe"

	"github.com/dchest/jsmin"
	"github.com/ztrue/tracerr"
)

//////////////////////////////////////////////////////

//var UGeneral _General_s
type _General_s struct {
	/*
		RandString generates a random string with uppercase and lowercase letters of the given length.

		-----------------------------------------------------------

		– Params:
		  - letters_num – the length of the string to generate

		– Returns:
		  - the generated string
	*/
	RandString func(letters_num int) string
	/*
		FindAllIndexes finds all the indexes of a substring in a string.

		-----------------------------------------------------------

		– Params:
		  - s – the string to search in
		  - substr – the substring to search for

		– Returns:
		  - the indexes of the substring in the string
	*/
	FindAllIndexes func(s string, substr string) []int
	/*
		GetFullErrorMsg gets the full error message from an error, including its stacktrace.

		-----------------------------------------------------------

		– Params:
		  - err – the error to get the full message from

		– Returns:
		  - the full error message
	*/
	GetFullErrorMsg func(err any) string
	/*
		PanicGENERAL panics with a custom string as the error.

		This function *never* returns.

		-----------------------------------------------------------

		– Params:
		  - err – the string to panic with
	*/
	PanicGENERAL func(err string)
	/*
		ToJson converts the given data to a JSON string and indents it.

		All the needed fields of the struct must be exported like with json.Marshal().

		-----------------------------------------------------------

		– Params:
		  - v – the data to convert to Json. Check the json.Marshal function for more info (used directly here).

		– Returns:
		  - true if the file was written successfully, false otherwise
	*/
	ToJson func(v any) *string
	/*
		FromJson minifies and parses the given JSON data.

		All the needed fields of the struct must be exported like with json.Marshal().

		-----------------------------------------------------------

		– Params:
		  - json_str – the JSON string to parse
		  - parsed_data – a pointer of where to write the parsed data to

		– Returns:
		  - true if the data was parsed correctly, false otherwise
	*/
	FromJson func(json_data []byte, parsed_data any) bool
}
//////////////////////////////////////////////////////

///////////////////////////
// Took from https://stackoverflow.com/a/31832326/8228163, by https://stackoverflow.com/users/1705598/icza.
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)
var src = rand.NewSource(time.Now().UnixNano())

func RandStringGENERAL(letters_num int) string {
	// Original function name: RandStringBytesMaskImprSrcUnsafe
	b := make([]byte, letters_num)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := letters_num-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}

///////////////////////////

func FindAllIndexesGENERAL(s string, substr string) []int {
	var indexes []int = nil
	var chars_processed int = 0
	var s_len int = len(s)
	for i := 0; i < s_len; i++ {
		var idx int = strings.Index(s[chars_processed:], substr)
		if -1 == idx {
			break
		}

		indexes = append(indexes, idx + chars_processed)
		chars_processed += idx + len(substr)
	}

	return indexes
}

func GetFullErrorMsgGENERAL(err any) string {
	return tracerr.SprintSource(tracerr.Wrap(err.(error)), 0)
}

func PanicGENERAL(err string) {
	panic(errors.New(err))
}

func ToJsonGENERAL(v any) *string {
	json_data, err := json.Marshal(v)
	if nil != err {
		return nil
	}

	var dst bytes.Buffer
	if nil == json.Indent(&dst, json_data, "", "\t") {
		json_data = dst.Bytes()
	}

	var json_str string = string(json_data)

	return &json_str
}

func FromJsonGENERAL(json_data []byte, parsed_data any) bool {
	var json_final []byte = nil
	var json_min, err = jsmin.Minify(json_data)
	if nil == err {
		json_final = json_min
	} else {
		// If the minifier fails, try to parse the original JSON (probably won't work, but I'll let Unmarshal() decide).
		json_final = json_data
	}

	return nil == json.Unmarshal(json_final, parsed_data)
}

/*
getVariableInfoGENERAL gets general information about a variable in a string in a default format.

-----------------------------------------------------------

– Params:
  - v – the variable to get the information about

– Returns:
  - the information about the variable
*/
func getVariableInfoGENERAL(v any) string {
	return "Information about it:" +
		"\n- Type of information (%T): " + fmt.Sprintf("%T", v) +
		"\n- Value of information (%+v): " + fmt.Sprintf("%+v", v) +
		"\n- Go representation of the value (%#v): " + fmt.Sprintf("%#v", v)
}
