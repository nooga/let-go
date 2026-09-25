#!/usr/bin/env python3
"""Run a balanced, interleaved A/B process benchmark.

Every round measures both commands. Odd rounds run A then B; even rounds run B
then A. The JSON includes every timing, the execution order, and paired B/A
ratios so the comparison does not depend on two unrelated sample minima.
"""

import argparse
import json
import statistics
import subprocess
import time
from pathlib import Path


def measure(command: list[str]) -> float:
    started = time.perf_counter()
    completed = subprocess.run(
        command,
        stdin=subprocess.DEVNULL,
        stdout=subprocess.DEVNULL,
        stderr=subprocess.DEVNULL,
        check=False,
    )
    elapsed = time.perf_counter() - started
    if completed.returncode != 0:
        raise SystemExit(
            f"command failed with exit {completed.returncode}: {command!r}"
        )
    return elapsed


def summary(name: str, times: list[float]) -> dict:
    return {
        "command": name,
        "name": name,
        "mean": statistics.fmean(times),
        "stddev": statistics.stdev(times) if len(times) > 1 else 0.0,
        "median": statistics.median(times),
        "min": min(times),
        "max": max(times),
        "times": times,
        "exit_codes": [0] * len(times),
    }


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--name-a", required=True)
    parser.add_argument("--name-b", required=True)
    parser.add_argument("--binary-a", required=True)
    parser.add_argument("--binary-b", required=True)
    parser.add_argument("--workload", required=True)
    parser.add_argument("--warmup", type=int, default=3)
    parser.add_argument("--runs", type=int, default=15)
    parser.add_argument("--output", required=True)
    args = parser.parse_args()
    if args.warmup < 0 or args.runs < 2:
        parser.error("--warmup must be non-negative and --runs must be at least 2")

    commands = {
        "A": [args.binary_a, args.workload],
        "B": [args.binary_b, args.workload],
    }
    for round_index in range(args.warmup):
        order = ("A", "B") if round_index % 2 == 0 else ("B", "A")
        for label in order:
            measure(commands[label])

    times = {"A": [], "B": []}
    orders = []
    paired_ratios = []
    for round_index in range(args.runs):
        order = ("A", "B") if round_index % 2 == 0 else ("B", "A")
        orders.append("".join(order))
        pair = {}
        for label in order:
            pair[label] = measure(commands[label])
            times[label].append(pair[label])
        paired_ratios.append(pair["B"] / pair["A"])

    data = {
        "schema": "let-go-interleaved-ab-v1",
        "warmup_rounds": args.warmup,
        "measured_rounds": args.runs,
        "round_orders": orders,
        "results": [
            summary(args.name_a, times["A"]),
            summary(args.name_b, times["B"]),
        ],
        "paired_b_over_a": {
            "mean": statistics.fmean(paired_ratios),
            "stddev": statistics.stdev(paired_ratios),
            "median": statistics.median(paired_ratios),
            "min": min(paired_ratios),
            "max": max(paired_ratios),
            "ratios": paired_ratios,
        },
    }
    Path(args.output).write_text(json.dumps(data, indent=2) + "\n")

    a, b = data["results"]
    paired = data["paired_b_over_a"]
    print(f"{a['name']}: median {a['median']:.4f}s, mean {a['mean']:.4f}s")
    print(f"{b['name']}: median {b['median']:.4f}s, mean {b['mean']:.4f}s")
    print(
        "paired B/A: "
        f"median {paired['median']:.4f}x, mean {paired['mean']:.4f}x, "
        f"stddev {paired['stddev']:.4f}"
    )


if __name__ == "__main__":
    main()
