#!/usr/bin/env python

import json
import sys
from urlparse import parse_qs, urlparse

import pafy

# encoding=utf8
reload(sys)
sys.setdefaultencoding('utf8')

url = sys.argv[1]
parsed_url = urlparse(url)
filename = parse_qs(parsed_url.query)["v"][0]
video = pafy.new(url)
best = video.getbest()
best.download(quiet=True, filepath=filename)
v_dict = {
    "filename": filename.decode('utf-8'),
    "url": url,
    "title": video.title.decode('utf-8'),
    "thumb": video.thumb,
    "length": video.length,
    "author": video.author.decode('utf-8'),
    "id": video.videoid
}
print json.dumps(v_dict)
