import json
import logging
import os
import sys
from pathlib import Path

import functions_framework
import webvtt
from cloudevents.http import CloudEvent
from google.api_core.exceptions import GoogleAPIError, NotFound
from google.cloud import storage

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

try:
    gcs_client = storage.Client()
except Exception as e:
    logger.error(f"Failed to initialize Google Cloud Storage client: {e}")
    sys.exit()


def download_from_gcs(bucket_name: str, file_name: str, local_path: Path) -> None:
    """
    Downloads a file from Google Cloud Storage.

    Args:
        bucket_name: The name of the GCS bucket.
        file_name: The name of the file in the bucket.
        local_path: The local path to save the downloaded file.

    Returns:
        The local_path if download is successful.

    Raises:
        google.cloud.exceptions.NotFound: If the bucket or blob does not exist.
        google.api_core.exceptions.GoogleAPIError: For other GCS related errors.
        IOError: If there's an issue writing to local_path.
    """
    try:
        bucket = gcs_client.bucket(bucket_name)
        blob = bucket.blob(file_name)

        # Ensure parent directories exist
        local_path.parent.mkdir(parents=True, exist_ok=True)

        with open(local_path, "wb") as f:
            blob.download_to_file(f)

        logger.info(
            f"Downloaded gs://{bucket_name}/{file_name} to {local_path}")

    except NotFound as e:
        logger.error(
            f"File not found in GCS: gs://{bucket_name}/{file_name}. Error: {e}")
        raise
    except GoogleAPIError as e:
        logger.error(
            f"Google API Error downloading {file_name} from GCS: {e}", exc_info=True)
        raise
    except Exception as e:
        logger.error(
            f"General error downloading {file_name} from GCS: {e}", exc_info=True)
        raise


def convert_vtt_to_json(vtt_local_path: Path) -> Path:
    """
    Converts a VTT file to a JSON file.

    Args:
        vtt_local_path: The local path to the VTT file.

    Returns:
        The local path to the generated JSON file.

    Raises:
        FileNotFoundError: If the VTT file does not exist.
        webvtt.errors.MalformedFileError: If the VTT file is not valid.
        IOError: If there's an issue writing the JSON file.
    """

    video_id = vtt_local_path.stem  # stem gets "file_name" without extension

    json_local_path = Path(f'/tmp/{video_id}.json')

    try:
        # Ensure parent directories exist
        json_local_path.parent.mkdir(parents=True, exist_ok=True)

        captions: list[dict] = []

        for cap in webvtt.read(str(vtt_local_path)):
            cap_entry = {
                "video_id": video_id,
                "start": cap.start,
                "end": cap.end,
                "text": cap.text
            }
            captions.append(cap_entry)

        with open(json_local_path, "w") as f:
            json.dump(captions, f)

        return json_local_path

    except FileNotFoundError as e:
        logger.error(f"VTT file not found at {vtt_local_path}. Error: {e}")
        raise
    except webvtt.errors.MalformedFileError as e:
        logger.error(f"Malformed VTT file: {vtt_local_path}. Error: {e}")
        raise
    except Exception as e:
        logger.error(
            f"Error converting VTT to JSON for {vtt_local_path}: {e}", exc_info=True)
        raise


def upload_to_gcs(local_path: Path, bucket_name: str, dest_blob_name: str) -> None:
    """
    Uploads a file to Google Cloud Storage.

    Args:
        local_path: The local path of the file to upload.
        bucket_name: The name of the GCS bucket.
        dest_blob_name: The destination blob name in GCS.

    Raises:
        google.api_core.exceptions.GoogleAPIError: For GCS related errors.
        FileNotFoundError: If the source_file_path does not exist.
    """
    try:
        bucket = gcs_client.bucket(bucket_name)
        blob = bucket.blob(dest_blob_name)
        blob.upload_from_filename(str(local_path))

        logger.info(
            f"Uploaded {local_path} to gs://{bucket_name}/{dest_blob_name}")
    except FileNotFoundError as e:
        logger.error(
            f"Source file not found for upload: {local_path}. Error: {e}")
        raise
    except Exception as e:
        logger.error(
            f"Error uploading {local_path} to GCS as {dest_blob_name}: {e}", exc_info=True)
        raise


@functions_framework.cloud_event
def vtt_to_json_converter(cloud_event: CloudEvent) -> None:

    data = cloud_event.data
    bucket_name = data.get('bucket')
    file_name = data.get('name')

    if not bucket_name or not file_name:
        logger.error("Error: Missing bucket or file namne in event data")
        return

    json_dest_folder = os.getenv('JSON_FOLDER', 'converted_json/')

    if not json_dest_folder.endswith('/'):
        json_dest_folder += '/'

    logger.info(f'Triggered by file: gs://{bucket_name}/{file_name}')

    vtt_local_path = Path(f'/tmp/{Path(file_name).name}')
    json_local_path = None

    if vtt_local_path.suffix != '.vtt':
        logger.warning(f"Skipping non-VTT file: {vtt_local_path}")
        return

    try:
        # 1. Download VTT from GCS
        download_from_gcs(bucket_name, file_name, vtt_local_path)

        # 2. Convert VTT to JSON
        json_local_path = convert_vtt_to_json(vtt_local_path)

        # 3. Upload JSON back to GCS
        json_blob_name = f'{json_dest_folder}{json_local_path.name}'
        upload_to_gcs(json_local_path, bucket_name, json_blob_name)

    except Exception as e:
        logger.critical(
            f"Unhandled error processing {file_name}: {e}", exc_info=True)

    finally:
        # Clean up local files regardless of success or failure
        for p in [vtt_local_path, json_local_path]:
            if p and p.exists():
                try:
                    p.unlink()  # Delete the file
                    logger.info(f"Cleaned up temporary file: {p}")
                except OSError as e:
                    logger.warning(
                        f"Failed to clean up temporary file {p}: {e}")
