# Core Radiation Detector GQ GMC-800 Geiger Counter

A simple Go-based tool for interfacing with the **Core Radiation Detector GQ GMC-800 Geiger Counter** via USB serial connection.  
It reads real-time radiation data from the device and exposes it via a local HTTP server in JSON format.

## Features

- Connects to GQ GMC-800 over USB (serial interface)
- Parses CPM (counts per minute) values from device output
- Launches a local HTTP server on port `8091`
- Serves current radiation data in JSON format:
  ```json
  {"cpm": 20}
  ```

## Platform

Tested and verified on **macOS only**.  
Other platforms (Linux, Windows) are untested and may require adjustments (e.g., serial port naming).

## Protocol Reference

This tool is based on the official serial communication protocol from GQ Electronics:  
**GQ-RFC1201.txt** — copied directly from the [original specification](https://www.gqelectronicsllc.com/download/GQ-RFC1201.txt).

## Requirements

- Go 1.18 or later
- GQ GMC-800 device connected via USB

## Notes
- The tool uses go.bug.st/serial for serial communication.
- This is a minimal utility intended for diagnostics, logging, or data forwarding.

## License

MIT License

## Author

[Aleksei Rytikov](https://github.com/chlp)