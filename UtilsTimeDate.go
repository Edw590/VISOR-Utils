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
	"time"
)

//////////////////////////////////////////////////////

//var UDateTime _DateTime_s
type _DateTime_s struct {
	/*
		GetDateTimeStr gets the current time and date in the format DATE_TIME_FORMAT.

		-----------------------------------------------------------

		– Returns:
		  - the current time and date in the default format
	*/
	GetDateTimeStr func() string
	/*
		GetTimeStr gets the current time in the format TIME_FORMAT.

		-----------------------------------------------------------

		– Returns:
		  - the current time in the default format
	*/
	GetTimeStr func() string
	/*
		GetDateStr gets the current date in the format DATE_FORMAT.

		-----------------------------------------------------------

		– Returns:
		  - the current date in the default format
	*/
	GetDateStr func() string
}
//////////////////////////////////////////////////////

const TIME_FORMAT string = "15:04:05"
const DATE_FORMAT string = "2006-01-02"
const DATE_TIME_FORMAT string = DATE_FORMAT + " -- " + TIME_FORMAT + " (MST)"

func GetDateTimeStrTIMEDATE() string {
	return time.Now().Format(DATE_TIME_FORMAT)
}

func GetDateStrTIMEDATE() string {
	return time.Now().Format(DATE_FORMAT)
}

func GetTimeStrTIMEDATE() string {
	return time.Now().Format(TIME_FORMAT)
}
