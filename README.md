# Youtube Podcast App


This project is an attempt to build a system that downloads youtube videos
and create an rss feed for podcast clients. 

Podcast clients can download contents and this allows you to consume the content while you are offline.

`youtube-podcast` app lets you watch youtube videos offline.

Check out [this post](http://ilkinulas.github.io/life/2016/02/12/dont-make-commuting-ruin-your-life.html) if you like podcasts.

### Requirements

 * *[pafy](https://github.com/mps-youtube/pafy)* 
 
 `pafy` is a python library that is used to download content and metadata from youtube.
 
 ```bash
 pip install pafy
 ```
 
 * *[youtube-dl](https://rg3.github.io/youtube-dl/)*
 
 `youtube-dl` is a command line tool for downloading videos from youtube. `pafy` requires `youtube-dl`.
 
 ```bash
 pip install youtube-dl
 ```

 * *AWS Account*
 
 This release stores the downloaded youtube videos in an AWS S3 bucket. Podcast clients
 will download the videos from S3. 


## Build & Run
```bash
$ go build
$ ./youtube-podcast --config config.toml
```

## Configuration

Configuration file [config.toml](./config.toml) is self explanatory.


## How does it work?

`youtube-podcast` app has three parts:
 1. Url Queue
 2. Video Downloader
 3. Rss Generator
 
 These three parts are bundled in the same binary.
 
 
### 1. Url Queue
When `save` endpoint receives a URL it persists it. The persistent queue is
implemented by using sqlite3 database.

Sample `save` request:  

> http://host:port/save?url=https://www.youtube.com/watch?v=UdiqXGCzMUo

### 2. Video Downloader

When `youtube-podcast` app starts it launches a goroutine (video downloader). 
Video downloader consumes queued urls one by one. It downloads the url using the 
python video downloader. Downloaded video is stored in the working directory of 
`youtube-podcast`. After a successful download file is uploaded to S3 and removed from 
file system.

Video downloader persists video metadata and public S3 link of the video. This data is then used by
the Rss Generator to generate a podcast feed.

Persistence layer is implemented by using sqlite3. 

### Rss Generator

The `rss` endpoint generates podcast rss feeds. The feed can be customized by updating the
`Podcast` section of the config file.

> http://host:port/rss



#### Development environment
I use [localstack](https://github.com/localstack/localstack) while developing `youtube-podcast` app.

To start a local S3 provider execute the following command.
```bash
TMPDIR=/private$TMPDIR SERVICES=s3 docker-compose -f docker/localstack-docker-compose.yml up -d
```

Then create the S3 bucket for storing downloaded videos with the below command:

```bash
aws --endpoint-url=http://localhost:4572 s3api create-bucket --bucket ilkinulas-youtube-podcast
```
 
