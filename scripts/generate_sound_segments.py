#!/usr/bin/env python3
import argparse
import json
import re
import shutil
import subprocess
import sys
from pathlib import Path


SILENCE_START_RE = re.compile(r"silence_start:\s*([0-9.]+)")
SILENCE_END_RE = re.compile(r"silence_end:\s*([0-9.]+)")
DURATION_RE = re.compile(r"Duration:\s*(\d+):(\d+):(\d+(?:\.\d+)?)")


def read_items(path):
    return [line.strip() for line in Path(path).read_text(encoding="utf-8").splitlines() if line.strip()]


def parse_duration(output):
    match = DURATION_RE.search(output)
    if not match:
        return None
    hours, minutes, seconds = match.groups()
    return int(hours) * 3600 + int(minutes) * 60 + float(seconds)


def run_silencedetect(audio, noise, min_silence):
    if shutil.which("ffmpeg") is None:
        raise SystemExit("Chyba: ffmpeg není dostupný v PATH. Nainstaluj ffmpeg a spusť skript znovu.")

    cmd = [
        "ffmpeg",
        "-hide_banner",
        "-i",
        str(audio),
        "-af",
        f"silencedetect=noise={noise}:d={min_silence}",
        "-f",
        "null",
        "-",
    ]
    result = subprocess.run(cmd, capture_output=True, text=True)
    output = result.stderr + "\n" + result.stdout
    if result.returncode != 0:
        raise SystemExit(f"Chyba: ffmpeg silencedetect selhal.\n\n{output}")

    starts = [float(value) for value in SILENCE_START_RE.findall(output)]
    ends = [float(value) for value in SILENCE_END_RE.findall(output)]
    return starts, ends, parse_duration(output)


def build_speech_blocks(silence_starts, silence_ends, duration):
    blocks = []
    cursor = 0.0

    for index, silence_start in enumerate(silence_starts):
        if silence_start > cursor:
            blocks.append({"start": cursor, "end": silence_start})
        if index < len(silence_ends):
            cursor = silence_ends[index]
        else:
            cursor = duration if duration is not None else silence_start

    if duration is not None and duration > cursor:
        blocks.append({"start": cursor, "end": duration})

    return blocks


def fail_count_mismatch(cs_items, es_items, silence_starts, silence_ends, speech_blocks, segment_count):
    message = f"""Chyba: počet řečových bloků neodpovídá počtu textových položek.

CZ položek: {len(cs_items)}
ES položek: {len(es_items)}
silence starts: {len(silence_starts)}
silence ends: {len(silence_ends)}
speech blocks: {len(speech_blocks)}
výsledných segmentů: {segment_count}

Doporučení: zkus upravit --noise nebo --min-silence."""
    raise SystemExit(message)


def round_time(value):
    return round(value, 3)


def main():
    parser = argparse.ArgumentParser(description="Generate timed vocabulary segments from an MP3 and CZ/ES text files.")
    parser.add_argument("--audio", required=True)
    parser.add_argument("--cs", required=True)
    parser.add_argument("--es", required=True)
    parser.add_argument("--file", required=True)
    parser.add_argument("--title", required=True)
    parser.add_argument("--lesson", required=True)
    parser.add_argument("--output", required=True)
    parser.add_argument("--noise", default="-35dB")
    parser.add_argument("--min-silence", default="0.35")
    parser.add_argument("--group-size", type=int, default=2)
    args = parser.parse_args()

    if args.group_size < 1:
        raise SystemExit("Chyba: --group-size musí být alespoň 1.")

    cs_items = read_items(args.cs)
    es_items = read_items(args.es)
    if len(cs_items) != len(es_items):
        raise SystemExit(f"Chyba: počty položek nesedí. CZ: {len(cs_items)}, ES: {len(es_items)}.")

    silence_starts, silence_ends, duration = run_silencedetect(args.audio, args.noise, args.min_silence)
    speech_blocks = build_speech_blocks(silence_starts, silence_ends, duration)

    if len(speech_blocks) % args.group_size != 0:
        fail_count_mismatch(cs_items, es_items, silence_starts, silence_ends, speech_blocks, len(speech_blocks) // args.group_size)

    segment_count = len(speech_blocks) // args.group_size
    if segment_count != len(cs_items):
        fail_count_mismatch(cs_items, es_items, silence_starts, silence_ends, speech_blocks, segment_count)

    segments = []
    for index, (spanish, czech) in enumerate(zip(es_items, cs_items)):
        group = speech_blocks[index * args.group_size : (index + 1) * args.group_size]
        segments.append(
            {
                "start": round_time(group[0]["start"]),
                "end": round_time(group[-1]["end"]),
                "spanish": spanish,
                "czech": czech,
            }
        )

    data = {
        "file": args.file,
        "title": args.title,
        "lesson": args.lesson,
        "segments": segments,
    }

    output_path = Path(args.output)
    output_path.parent.mkdir(parents=True, exist_ok=True)
    output_path.write_text(json.dumps(data, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")

    print(f"CZ položek: {len(cs_items)}")
    print(f"ES položek: {len(es_items)}")
    print(f"silence starts: {len(silence_starts)}")
    print(f"silence ends: {len(silence_ends)}")
    print(f"speech blocks: {len(speech_blocks)}")
    print(f"výsledných segmentů: {len(segments)}")
    print(f"output: {output_path}")


if __name__ == "__main__":
    main()
