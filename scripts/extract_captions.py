import argparse
import json
import os

import webvtt


def parse_captions(filepath: str):
    print("Reading from ", filepath)

    captions: list[dict] = []

    video_id = filepath.split("/")[-1].replace('.vtt', '')
    output_dir: str = 'tmp/captions_parsed/'
    json_file: str = f"{output_dir}{video_id}.json"

    video_id = video_id.replace('.en', '')

    for cap in webvtt.read(filepath):
        cap_entry = {
            "video_id": video_id,
            "start": cap.start,
            "end": cap.end,
            "text": cap.text
        }
        captions.append(cap_entry)

    os.makedirs(os.path.dirname(json_file), exist_ok=True)

    with open(json_file, "w") as f:
        json.dump(captions, f, indent=4)


if __name__ == "__main__":
    argparser = argparse.ArgumentParser(
        description="Convert vtt captions to JSON")
    argparser.add_argument("filepath", type=str)
    args = argparser.parse_args()

    parse_captions(args.filepath)
