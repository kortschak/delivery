// Copyright ©2018 Dan Kortschak. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The delivery program reads a Posts.json file downloaded from Google
// Takeout (https://takeout.google.com/settings/takeout) for Google+
// Communities. delivery will download each of the posts listed in the
// Posts.json file and archive them and the Posts.json file into a zip
// archive. This makes up for the absence on any capacity in Takeout to
// give more than the metadata for download.
package main

import (
	"archive/zip"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"
)

type gPlusPost struct {
	CommunityName string `json:"communityName"`
	Post          []struct {
		URL          string `json:"url"`
		CreationTime string `json:"creationTime"`
		UpdateTime   string `json:"updateTime"`
	} `json:"post"`
}

func main() {
	in := flag.String("i", "", "infile JSON for posts (required)")
	out := flag.String("o", "", "outfile for posts (required)")
	flag.Parse()
	if *in == "" || *out == "" {
		flag.Usage()
		os.Exit(2)
	}

	b, err := ioutil.ReadFile(*in)
	if err != nil {
		log.Fatalf("failed to read posts metadata: %v", err)
	}
	var posts gPlusPost
	err = json.Unmarshal(b, &posts)
	if err != nil {
		log.Fatalf("failed to unmarshal posts: %v", err)
	}

	if filepath.Ext(*out) != ".zip" {
		*out += ".zip"
	}
	f, err := os.Create(*out)
	if err != nil {
		log.Fatalf("failed to create %q: %v", *out, err)
	}
	z := zip.NewWriter(f)
	defer z.Close()

	var last time.Time
	for _, p := range posts.Post {
		ct, err := strconv.ParseInt(p.CreationTime, 10, 64)
		if err != nil {
			log.Fatalf("failed to parse ctime: %v", err)
		}
		ctime := time.Unix(ct/1e3, (ct%1e3)*1e6).UTC()

		ut, err := strconv.ParseInt(p.UpdateTime, 10, 64)
		if err != nil {
			log.Fatalf("failed to parse utime: %v", err)
		}
		utime := time.Unix(ut/1e3, (ut%1e3)*1e6).UTC()
		if utime.After(last) {
			last = utime
		}

		u, err := url.Parse(p.URL)
		if err != nil {
			log.Fatalf("failed to parse url: %v", err)
		}

		base := fmt.Sprintf("%s_%s", ctime.Format(time.RFC3339), path.Base(u.Path))
		fmt.Println(base)

		resp, err := http.Get(p.URL)
		if err != nil {
			log.Fatalf("failed to GET %v: %v", p.URL, err)
		}

		fh := &zip.FileHeader{
			Name:     base + ".html",
			Method:   zip.Deflate,
			Modified: ctime,
		}
		w, err := z.CreateHeader(fh)
		if err != nil {
			log.Fatalf("failed to create %q: %v", fh.Name, err)
		}
		_, err = io.Copy(w, resp.Body)
		if err != nil {
			log.Fatalf("failed to copy %q: %v", fh.Name, err)
		}
		resp.Body.Close()
	}

	fh := &zip.FileHeader{
		Name:     *in,
		Method:   zip.Deflate,
		Modified: last,
	}
	w, err := z.CreateHeader(fh)
	if err != nil {
		log.Fatalf("failed to create %q: %v", *in, err)
	}
	_, err = w.Write(b)
	if err != nil {
		log.Fatalf("failed to write %q: %v", *in, err)
	}

}
