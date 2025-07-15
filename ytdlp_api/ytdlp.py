import logging
import os
import subprocess
from urllib.parse import parse_qs, urlparse

from flask import Flask, jsonify, request
from google.cloud import storage

# ========================= Setup =========================
app = Flask(__name__)

# Local folder to download captions files
DOWNLOAD_DIR = "/tmp/downloads"

# Instantiate GCS client
CAPTIONS_BUCKET = os.getenv('CAPTIONS_BUCKET')
CAPTIONS_FOLDER = os.getenv('CAPTIONS_FOLDER')

client = storage.Client()
bucket = client.get_bucket(CAPTIONS_BUCKET)

logging.basicConfig(
    level=logging.INFO,
    format='[%(asctime)s] [%(levelname)s] %(message)s',
)
logger = logging.getLogger(__name__)


# ========================= Custom Types & Errors =========================


class YtdlpFetchError(Exception):
    """Custom exception for ytdlp fetch errors"""


# ========================= Helper Functions =========================

def fetch_metadata(url: str) -> tuple[str, str]:
    """Fetch video_id and video_title from a youtube url

    Args:
        url (str): valid youtube url

    Raises:
        YtdlpFetchError: error in calling yt-dlp
        YtdlpFetchError: error in output from yt-dlp is not in the form [video_id, video_title]

    Returns:
        tuple[str, str]: video_id, and video_title
    """

    cmd = ['yt-dlp',
           '--get-id',
           '--get-title',
           '--no-warnings',
           '--skip-download',
           url]

    try:
        res = subprocess.check_output(cmd, stderr=subprocess.STDOUT, text=True) \
            .strip().split('\n')

        if len(res) != 2:
            raise YtdlpFetchError('Unexpected yt-dlp output format')
        video_title = res[0]
        video_id = res[1]

        return video_id, video_title

    except subprocess.CalledProcessError as e:
        raise YtdlpFetchError(f'yt-dlp failed: {e.output.strip()}')


def download_captions(url: str, output_dir: str) -> str:
    """Download video captions to a given folder

    Args:
        url (str): valid youtube url
        output_dir (str): directory to download to

    Raises:
        YtdlpFetchError: error in calling yt-dlp

    Returns:
        str: local path of the caption file
    """
    video_id = extract_video_id(url)
    os.makedirs(output_dir, exist_ok=True)

    cmd = ['yt-dlp',
           '--write-subs',
           '--write-auto-subs',
           '--no-warnings',
           '--sub-langs', 'en',
           '--skip-download',
           '-o', f'{output_dir}/%(id)s.%(ext)s',
           url]

    try:
        res = subprocess.check_output(cmd, text=True).strip().split('\n')
        logger.info(res)

        caption_path = os.path.join(output_dir, f"{video_id}.en.vtt")

        if not os.path.exists(caption_path):
            raise YtdlpFetchError(f"Caption file not found at {caption_path}")

        return caption_path

    except subprocess.CalledProcessError as e:
        raise YtdlpFetchError(f'yt-dlp failed: {e.output.strip()}')


def extract_video_id(url: str) -> str:
    """Extract the video ID from a YouTube URL."""
    parsed = urlparse(url)

    # Handle standard YouTube URLs
    if parsed.hostname in ("www.youtube.com", "youtube.com"):
        return parse_qs(parsed.query).get("v", [None])[0]

    # Handle short URLs like youtu.be/<id>
    if parsed.hostname == "youtu.be":
        return parsed.path.lstrip("/")

    raise ValueError("Unsupported YouTube URL format")


# ========================= Endpoints =========================


@app.route('/health', methods=['GET'])
def health_check():
    """Health check endpoint."""
    return jsonify({"status": "OK"}), 200


@app.route('/v1/metadata', methods=['GET'])
def get_metadata():

    url = request.args.get('url')

    if not url:
        return jsonify({'error': 'Valid Youtube video url is required'}), 400

    try:
        video_id, video_title = fetch_metadata(url)
        return jsonify({'id': video_id, 'title': video_title})

    except YtdlpFetchError as e:
        return jsonify({'error': str(e)}), 500
    except Exception as e:
        return jsonify({'error': f"An unexpected error occurred: {str(e)}"}), 500


@app.route('/v1/captions', methods=['POST'])
def get_captions():

    data = request.json

    url = data['url']
    output_dir = data['output_dir']

    try:
        caption_path = download_captions(url, output_dir)

    except YtdlpFetchError as e:
        return jsonify({'error': str(e)}), 500
    except Exception as e:
        return jsonify({'error': f"An unexpected error occurred: {str(e)}"}), 500

    return jsonify({"message": f"Successfully downloaded captions for {url} as {caption_path}"}), 200


@app.route('/v2/captions', methods=['POST'])
def upload_captions_to_gcs():

    # 1. Get url of yt video
    url = request.args.get('url')
    if not url:
        return jsonify({'error': 'Valid Youtube video url is required'}), 400

    try:
        # 2. Download captions using ytdlp
        caption_path = download_captions(url, DOWNLOAD_DIR)
        file_name = os.path.basename(caption_path)

        # 3. Upload to GCS
        blob = bucket.blob(f"{CAPTIONS_FOLDER}/{file_name}")
        blob.upload_from_filename(caption_path)

        # 4. Clean up locally
        if os.path.exists(caption_path):
            os.remove(caption_path)

        return jsonify({
            "message": f"Successfully downloaded and uploaded captions for {url} to GCS",
            "file": file_name
        }), 200

    except YtdlpFetchError as e:
        return jsonify({'error': str(e)}), 500
    except Exception as e:
        return jsonify({'error': f"An unexpected error occurred: {str(e)}"}), 500


if __name__ == "__main__":

    port = int(os.environ.get("PORT"))
    app.run("0.0.0.0", port)
