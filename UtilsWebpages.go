package Utils

import (
	"io"
	"net/http"
)

//////////////////////////////////////////////////////

var UWebpages _Webpages_s
type _Webpages_s struct {
	/*
		GetPageHtml gets the HTML of a page.

		-----------------------------------------------------------

		> Params:
		  - url – the URL of the page

		> Returns:
		  - the HTML of the page or nil if an error occurs
	*/
	GetPageHtml func(url string) *string
}
//////////////////////////////////////////////////////

func getPageHtmlTIMEDATE(url string) *string {
	resp, err := http.Get(url)
	if nil == err {
		body, err := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if resp.StatusCode <= 299 && nil == err {
			var ret string = string(body)

			return &ret
		}
	}

	return nil
}
