# yt-dlp server

Runs a server (Flask + Gunicorn) locally that uses [`yt-dlp`](https://github.com/yt-dlp/yt-dlp) to fetch captions from youtube videos, and uploads content to Google Cloud Storage (GCS).

---

## Prerequisites

- [Docker](https://www.docker.com/)
- [Terraform](https://developer.hashicorp.com/terraform/downloads)
- A GCP project with billing enabled
- `gcloud` CLI authenticated (`gcloud auth login`)


---


## Setup

### 1. Generate a service account key

Terraform already creates the service account `ytdlp-svc`, we want to go to the GCP console -> IAM -> Service Accounts, and create a key for the service account. Then rename and place the key in the secrets folder: `./secrets/ytdlp-svc-key.json`. 


### 2. Build the Docker image

From the root directory:

```bash
docker build -t ytdlp-api ./ytdlp_api
```

### 3. Run the API container

```bash
docker run \
  -p 5000:5000 \
  -v "$(pwd)/.secrets/ytdlp-svc-key.json:/secrets/key.json:ro" \
  -e GOOGLE_APPLICATION_CREDENTIALS=/secrets/key.json \
  -e CAPTIONS_BUCKET=banditsecret-captions \
  -e CAPTIONS_FOLDER=raw_vtt \
  -e YTDLP_PORT=5000 \
  ytdlp-api
```

Note: The terraform output contains `gcs_bucket_name` which contains the captions bucket.

---

## Endpoints

The `v1` endpoints are for the local docker-compose approach to ytdlp where the captions files are downloaded locally to a shared volume where another service picks up the captions files.


The `v2` endpoints are for the Google Cloud approach to ytdlp where the files are uploaded to a bucket.


### `GET /v1/metadata?url=<video_url>`

Fetches the video ID and title.

**Example:**

```bash
curl "http://localhost:5000/v1/metadata?url=https://www.youtube.com/watch?v=dQw4w9WgXcQ"
```


### `POST /v1/captions`

Downloads and places captions in a local folder.

**Request JSON:**

```json
{
  "url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
  "output_dir": "/tmp"
}
```

**Example:**

```bash
curl -X POST http://localhost:5000/v1/captions \
  -H "Content-Type: application/json" \
  -d '{"url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ", "output_dir": "/tmp"}'
```

### `POST /v2/captions`

Downloads and uploads captions to GCS.

**Example:**

```bash
curl -X POST http://localhost:5000/v2/captions?url=https://www.youtube.com/watch?v=dQw4w9WgXcQ
```

---

## Security Notes

- The service account key is written to `.secrets/ytdlp_svc_key.json` and picked up by Docker (at runtime)
- Never commit this file!

---
