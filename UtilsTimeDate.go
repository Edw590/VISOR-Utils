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

const TIME_FORMAT string = "15:04:05"
const DATE_FORMAT string = "2006-01-02"
const DATE_TIME_FORMAT string = DATE_FORMAT + " -- " + TIME_FORMAT + " (MST)"

/*
GetDateTimeStrTIMEDATE gets the current time and date in the format DATE_TIME_FORMAT.

-----------------------------------------------------------

– Returns:
  - the current time and date in the default format
*/
func GetDateTimeStrTIMEDATE(millis int64) string {
	return getTimeDateInFormat(millis, DATE_TIME_FORMAT)
}

/*
GetDateStrTIMEDATE gets the current date in the format DATE_FORMAT.

-----------------------------------------------------------

– Returns:
  - the current time in the default format
*/
func GetDateStrTIMEDATE(millis int64) string {
	return getTimeDateInFormat(millis, DATE_FORMAT)
}

/*
GetTimeStrTIMEDATE gets the current time in the format TIME_FORMAT.

-----------------------------------------------------------

– Returns:
  - the current date in the default format
*/
func GetTimeStrTIMEDATE(millis int64) string {
	return getTimeDateInFormat(millis, TIME_FORMAT)
}

/*
getTimeDateInFormat gets the time and/or date in the given format.

-----------------------------------------------------------

– Params:
  - millis – the time in milliseconds
  - format – the format to use

– Returns:
  - the time and/or date in the given format
 */
func getTimeDateInFormat(millis int64, format string) string {
	if millis != -1 {
		return time.Unix(0, millis*1e6).Format(format)
	} else {
		return time.Now().Format(format)
	}
}
