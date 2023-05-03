package templatemethod

import "testing"

func TestHTTPDownloader(t *testing.T) {
	var downloader Downloader = NewHTTPDownloader()

	downloader.Download("http://example.com/abc.zip")
	// Output:
	// prepare downloading
	// download http://example.com/abc.zip via http
	// http save
	// finish downloading
}

func TestFTPDownloader(t *testing.T) {
	var downloader Downloader = NewFTPDownloader()

	downloader.Download("ftp://example.com/abc.zip")
	// Output:
	// prepare downloading
	// download ftp://example.com/abc.zip via ftp
	// default save
	// finish downloading
}
