Basic EventBus Prototype

requirements
go version 1.24
python version 3.11 or 3.12

# Installation
pip install -r requirements.txt

# Usage
Start the EventBus server:
```
go run main.go
```
Start the EventBus client:
```
python client.py
```


Open the test.html file in a web browser.
You can use it to test the EventBus functionality.
if you subscribe to file.upload, you will see when a file is moved to watched_folder.
if you subscribe to audio.transcript, you will see when a transcription is generated.

For a simple test just move the person-test.flac file into the watched_folder.
Or you can manually publish a message via the browser on the topic file.upload
with the full file path to an audio file you want transcribed.
