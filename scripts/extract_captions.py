import argparse
import json
import os

import webvtt


def parse_captions(filepath: str):
    print("Reading from ", filepath)

    captions: list[dict] = []

    for cap in webvtt.read(filepath):

        cap_entry = {
            "start": cap.start,
            "end": cap.end,
            "text": cap.text
        }
        captions.append(cap_entry)

    output_dir: str = "tmp/captions_parsed/"
    json_file: str = f"{output_dir}{filepath.split('/')[-1].replace('.vtt', '.json')}"

    os.makedirs(os.path.dirname(json_file), exist_ok=True)

    with open(json_file, "w") as f:
        json.dump(captions, f, indent=4)


if __name__ == "__main__":
    argparser = argparse.ArgumentParser(
        description="Convert vtt captions to JSON")
    argparser.add_argument("filepath", type=str)
    args = argparser.parse_args()

    parse_captions(args.filepath)
