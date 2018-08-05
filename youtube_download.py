#!/usr/bin/env python

import json
import pafy
import sys

url = sys.argv[1]
download = False
if len(sys.argv) > 2:
    download = sys.argv[2] == "-d"

video = pafy.new(url)
if download:
    best = video.getbest()
    best.download(quiet=False)
else:
    v_dict = {
        "url": url,
        "title": video.title,
        "thumb": video.thumb,
        "length": video.length,
        "author": video.author,
        "id": video.videoid
    }
    print json.dumps(v_dict)
