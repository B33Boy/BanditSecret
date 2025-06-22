[![Go](https://github.com/B33Boy/BanditSecret/actions/workflows/go.yml/badge.svg)](https://github.com/B33Boy/BanditSecret/actions/workflows/go.yml)

# Bandit Secret
A tool that allows you to quickly find your favourite moments in youtube videos

## Features
- Search through captions of YouTube videos
- Automatically downloads English subtitles via yt-dlp
- Converts .vtt captions to structured JSON for easy processing
- Built with Go for speed

## Requirements
- Go 1.18+
- yt-dlp installed and the folder containing the executable is added to PATH (on Windows)
- Python 3 (converting .vtt caption file to JSON)

## Setting up Python environment and installing dependencies (Windows)
```
python -m venv venv
.\venv\Scripts\activate

pip install -r requirements.txt
```

## Running Locally 
Build the binary 
```
go build -o bin/server.exe ./cmd/server
```

```
.\bin\server.exe
```

## Running with Docker
Completely remove network, volume mount, and container
```
docker-compose down -v
```


```
docker-compose up --build
```


## Sending requests
```bash
curl --location '127.0.0.1:6969/v1/captions' \
--header 'Content-Type: text/plain' \
--data 'https://youtu.be/iTOKRWgjOlg'
```

## Searching for a word or phrase
```bash
curl --location '127.0.0.1:6969/v1/search?query=mystery+colony'
```

## License
MIT - use freely, give credit where it's due