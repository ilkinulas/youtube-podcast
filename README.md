# Youtube Podcast App

## ---> Work in progress <--- 


This is an attempt to build a system that downloads youtube videos
and create and rss feed to for podcast clients. 

I will be able to watch youtube videos offline without much friction. 

Http server:

* saves youtube video urls
* generates rss feed for podcast client apps.


Video downloader:
* works like a cron job.
* uses a python library to download youtube videos. (https://github.com/mps-youtube/pafy)
* videos are stored in a cloud storage (AWS S3)


