import argparse
import json
import os

import webvtt


def parse_captions(vttpath: str, jsonpath: str):
    print("[extract_captions.py] Reading from ", vttpath)

    captions: list[dict] = []

    # video_id = vttpath.split("/")[-1].replace('.vtt', '')

    video_id = jsonpath.split("/")[-1].replace('.en.json', "")
    # output_dir: str = 'tmp/captions_parsed/'
    # json_file: str = f"{output_dir}{video_id}.json"
    # video_id = video_id.replace('.en', '')

    for cap in webvtt.read(vttpath):
        cap_entry = {
            "video_id": video_id,
            "start": cap.start,
            "end": cap.end,
            "text": cap.text
        }
        captions.append(cap_entry)

    os.makedirs(os.path.dirname(jsonpath), exist_ok=True)

    with open(jsonpath, "w") as f:
        json.dump(captions, f, indent=4)


if __name__ == "__main__":
    argparser = argparse.ArgumentParser(
        description="Convert vtt captions to JSON")

    argparser.add_argument("vttpath", type=str)
    argparser.add_argument("jsonpath", type=str)

    args = argparser.parse_args()

    parse_captions(args.vttpath, args.jsonpath)
