#!/usr/bin/env python

import json
import sys

import pafy

url = sys.argv[1]
video = pafy.new(url)
best = video.getbest()
filename = best.download(quiet=True)
v_dict = {
    "filename": filename,
    "url": url,
    "title": video.title,
    "thumb": video.thumb,
    "length": video.length,
    "author": video.author,
    "id": video.videoid
}
print json.dumps(v_dict)
