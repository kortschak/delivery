# delivery

`delivery` provides additional convenience to Takeout for saving Google+ Communities data.

The program will download and archive HTML pages for each community post for a Google Plus community when given a Posts.json file downloaded using Google Takeout.

To use `delivery`, you will need either download the Google+ Communities metadata in JSON format using [Takeout](https://takeout.google.com/settings/takeout) and extract the `Posts.json` file or obtain the relevant file from the owner of the community if you do not have access to Takeout for the community. Then you can run `delivery`:

```
$ delivery -i Posts.json -o posts
```

This will generate an archive, `posts.zip`, containing the HTML pages for each of the posts, and the `Posts.json` file as a manifest.

Note that the HTML files are unmarked UTF-8 and may be incorrectly rendered by Firefox when read from as local files since Firefox does not currently default to UTF-8 for these. Either add `<meta charset="utf-8">` to the files, or set `intl.charset.fallback.utf8_for_file` to `true` in `about:config` to fix this.