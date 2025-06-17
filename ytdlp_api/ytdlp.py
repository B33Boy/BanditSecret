import os
import subprocess

from flask import Flask, jsonify, request

app = Flask(__name__)

DOWNLOAD_DIR = "/tmp/downloads"
os.makedirs(DOWNLOAD_DIR, exist_ok=True)


@app.route('/health', methods=['GET'])
def health_check():
    """Health check endpoint."""
    return jsonify({"status": "OK"}), 200


@app.route('/get_metadata', methods=['GET'])
def get_metadata():
    url = request.args.get('url')
    if not url:
        return jsonify({'error': 'Valid Youtube video url is required'}), 400

    try:
        cmd = ['yt-dlp',
               '--get-id',
               '--get-title',
               '--no-warnings',
               '--skip-download',
               url]

        res = subprocess.check_output(cmd, text=True).strip().split('\n')
        video_title = res[0] if len(res) > 0 else "N/A"
        video_id = res[1] if len(res) > 1 else "N/A"

        return jsonify({'id': video_id, 'title': video_title})

    except subprocess.CalledProcessError as e:
        return jsonify({'error': str(e)}), 500
    except Exception as e:
        return jsonify({'error': f"An unexpected error occurred: {str(e)}"}), 500


@app.route('/get_captions', methods=['POST'])
def get_captions():

    data = request.json

    url = data['url']
    output_dir = data['output_dir']

    try:
        cmd = ['yt-dlp',
               '--write-subs',
               '--write-auto-subs',
               '--no-warnings',
               '--sub-langs', 'en',
               '--skip-download',
               url,
               '-o',
               output_dir]

        res = subprocess.check_output(cmd, text=True).strip().split('\n')
        return jsonify({'response': res})

    except subprocess.CalledProcessError as e:
        return jsonify({'error': str(e)}), 500
    except Exception as e:
        return jsonify({'error': f"An unexpected error occurred: {str(e)}"}), 500
