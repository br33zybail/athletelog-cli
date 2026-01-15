![Go Version](https://img.shields.io)

# AthleteLog – Multi-Language Training Tracker

A personal workout logging system built by a former athlete pivoting to software engineering in 2026.

**Architecture highlight**: Demonstrates deliberate tool selection for each layer:
- **Go** – fast CLI orchestrator & subprocess glue
- **Rust** – safe, high-performance core calculations (estimated 1RM)
- **Python** – data analysis & static visualization (progress charts)
- **TypeScript** – interactive browser dashboard (table + live trends)

This project shows disciplined, incremental system building — much like a structured training block.

## Features
- Add/view workouts via CLI
- Rust-powered 1RM estimation
- Python-generated progress reports (PNG + summary)
- Interactive TS dashboard (table + Chart.js trends)

## Quick Start (One-Command Flow)

```bash
git clone https://github.com/br33zybail/athletelog-cli.git
cd athletelog-cli
make build              # Go + TS build + npm install
make run                # Run CLI
./athletelog-cli add 2026-01-15 squat 225 5
make dashboard          # Open browser view

## Requires Go 1.23+, Python 3.12+, and Node.js 20+

Setup (Detailed)

Go: go mod tidy && go build
Rust: cd rust-stats/stats && cargo build --release
Python: python3 -m venv .venv && source .venv/bin/activate && pip install -r python-report/requirements.txt
TypeScript: cd ts-dashboard && npm install && npm run build
