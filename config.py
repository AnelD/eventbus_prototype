"""Configuration file for setting application parameters."""

import logging

# WebSocket settings
WEBSOCKET_URI = "ws://localhost:8080"

# Audio directory for intermediate audio files
AUDIO_OUT_DIR = "./data/out"

# ASR (Automatic Speech Recognition) settings
# Default: openai/whisper-large-v3-turbo
ASR_MODEL_NAME = "openai/whisper-large-v3-turbo"
ASR_LANGUAGE = "german"

# Logging settings available: DEBUG, INFO, WARNING, ERROR
LOG_LEVEL = logging.DEBUG

# Log file name and location
LOG_FILE = "logs/app.log"
