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
	"mime/quotedprintable"
	"os"
	"strings"

	PersonalConsts "VISOR_S_L/PERSONAL_FILES_EOG"
)

// EmailInfo is the info needed to send an email through QueueEmail().
type EmailInfo struct {
	// Sender name (can be anything)
	Sender  string
	// Mail_to is the email address to send the email to.
	Mail_to string
	// Subject of the email.
	Subject string
	// Html is the HTML body of the email.
	Html    string
	// Multipart is the list of multipart items to attach to the email aside from the main HTML.
	Multiparts []Multipart
}

// Multipart is an item to attach to an email as described in RFC 1521.
type Multipart struct {
	Content_type              string
	Content_transfer_encoding string
	Content_id                string
	Body                      string
}

const RAND_STR_LEN int = 10

const TO_SEND_REL_FOLDER string = "to_send/"
const _EMAIL_MODELS_FOLDER string = "email_models/"

const _TEMP_EML_FILE string = "msg_temp.eml"

const MODEL_FILE_INFO string = "model_email_info.html"
const MODEL_FILE_RSS string = "model_email_rss.html"
const MODEL_FILE_YT_VIDEO string = "model_email_video_YouTube.html"
const MODEL_FILE_DISKS_SMART string = "model_email_disks_smart.html"
const _MODEL_FILE_MESSAGE_EML string = "model_message.eml"
/*
GetModelFileEMAIL returns the contents of an email model file.

-----------------------------------------------------------

– Params:
  - file_name – the name of the file

– Returns:
  - the contents of the file or nil if an error occurred
*/
func GetModelFileEMAIL(file_name string) *string {
	return getModDirMODULES(NUM_MOD_EmailSender).Add(_EMAIL_MODELS_FOLDER, file_name).ReadFile()
}

/*
QueueEmailEMAIL queues an email to be sent by the UEmail Sender module.

-----CONSTANTS-----
  - MODEL_FILE_INFO – model file for information emails.
  - MODEL_FILE_RSS – model file for RSS feed notification emails.
  - MODEL_FILE_YT_VIDEO – model file for YouTube video notification emails.
  - _MODEL_FILE_MESSAGE_EML – model file for the main message.eml file.
-----CONSTANTS-----

-----------------------------------------------------------

– Params:
  - emailInfo – the email info
  - multiparts – the list of multipart items to attach to the email aside from the main HTML or nil to ignore

– Returns:
  - nil if the email was queued successfully, otherwise an error
*/
func QueueEmailEMAIL(emailInfo EmailInfo) error {
	var message_eml, _ string = prepareEmlEMAIL(emailInfo)

	var file_name string = ""
	var to_send_dir GPath = getModDataDirMODULES(NUM_MOD_EmailSender).Add(TO_SEND_REL_FOLDER)
	for {
		var rand_string string = RandStringGENERAL(RAND_STR_LEN)
		_, err := os.ReadFile(to_send_dir.Add(rand_string + emailInfo.Mail_to + ".eml").
			GPathToStringConversion())
		if nil != err {
			// If the file doesn't exist, choose that name.
			file_name = rand_string + emailInfo.Mail_to + ".eml"

			return getModDataDirMODULES(NUM_MOD_EmailSender).Add(TO_SEND_REL_FOLDER + file_name).
				WriteTextFile(message_eml)
		}
	}
}

/*
SendEmailEMAIL sends an email with the given message and receiver.

***DO NOT USE OUTSIDE THE EMAIL SENDER MODULE***

-----------------------------------------------------------

– Params:
  - message_eml – the complete message to be sent in EML format
  - mail_to – the receiver of the email

– Returns:
  - nil if the email was sent successfully, otherwise an error
*/
func SendEmailEMAIL(message_eml string, mail_to string) error {
	if err := getModTempDirMODULES(NUM_MOD_EmailSender).Add(_TEMP_EML_FILE).WriteTextFile(message_eml); nil != err {
		return err
	}
	_, err := ExecCmdSHELL(getCurlStringEMAIL(mail_to))

	return err
}

/*
ToQuotedPrintableEMAIL converts a string to a quoted printable string.

-----------------------------------------------------------

– Params:
  - str – the string to convert

– Returns:
  - the quoted printable string or nil if an error occurs
*/
func ToQuotedPrintableEMAIL(str string) *string {
	var ac bytes.Buffer
	w := quotedprintable.NewWriter(&ac)
	_, err := w.Write([]byte(str))
	if nil != err {
		return nil
	}
	err = w.Close()
	if nil != err {
		return nil
	}
	ret := ac.String()

	return &ret
}

/*
prepareEmlEMAIL prepares the EML file of the email.

-----------------------------------------------------------

– Params:
  - emailInfo – the email info
  - multiparts – the list of multipart items to attach to the email aside from the main HTML or nil to ignore

– Returns:
  - the email EML file to be sent
*/
func prepareEmlEMAIL(emailInfo EmailInfo) (string, string) {
	var message_eml string = *GetModelFileEMAIL(_MODEL_FILE_MESSAGE_EML)

	emailInfo.Html = strings.ReplaceAll(emailInfo.Html, "|3234_MSG_SUBJECT|", emailInfo.Subject)
	emailInfo.Html = strings.ReplaceAll(emailInfo.Html, "|3234_MSG_SENDER_NAME|", emailInfo.Sender)

	message_eml = strings.ReplaceAll(message_eml, "|3234_MSG_HTML|", *ToQuotedPrintableEMAIL(emailInfo.Html))
	message_eml = strings.ReplaceAll(message_eml, "|3234_MSG_SUBJECT|", emailInfo.Subject)
	message_eml = strings.ReplaceAll(message_eml, "|3234_MSG_SENDER_NAME|", emailInfo.Sender)

	var multiparts_str string = ""
	if nil != emailInfo.Multiparts {
		for _, multipart := range emailInfo.Multiparts {
			multiparts_str += "\n--|3234_MSG_BOUNDARY|\n" +
						"Content-Type: " + multipart.Content_type + "\n" +
						"Content-Transfer-Encoding: " + multipart.Content_transfer_encoding + "\n" +
						"Content-ID: <" + multipart.Content_id + ">\n" +
						"\n" +
						multipart.Body + "\n\n"
		}
	}
	message_eml = strings.ReplaceAll(message_eml, "|3234_MSG_MULTIPARTS|", multiparts_str)

	var msg_boundary string = RandStringGENERAL(25)
	for {
		if !strings.Contains(message_eml, msg_boundary) {
			break
		}
		msg_boundary = RandStringGENERAL(25)
	}
	message_eml = strings.ReplaceAll(message_eml, "|3234_MSG_BOUNDARY|", msg_boundary)


	return message_eml, emailInfo.Mail_to
}

/*
getCurlStringEMAIL gets the cURL string that sends an email with the default message file path and sender and receiver.

-----------------------------------------------------------

– Returns:
  - the string ready to be executed by the system
*/
func getCurlStringEMAIL(mail_to string) string {
	return "curl --location --connect-timeout 4294967295 "+/*--verbose*/" \"smtp://smtp.gmail.com:587\" --user \"" +
		PersonalConsts.VISOR_EMAIL_ADDR + ":" + PersonalConsts.VISOR_EMAIL_PW + "\" --mail-rcpt \"" + mail_to +
		"\" --upload-file \"" +	getModTempDirMODULES(NUM_MOD_EmailSender).Add(_TEMP_EML_FILE).GPathToStringConversion() +
		"\" --ssl"
}
