import time

import torch
from transformers.pipelines import pipeline
from transformers.pipelines.base import Pipeline


from audio_helper import AudioHelper
from logger_helper import LoggerHelper

import config

log = LoggerHelper(__name__).get_logger()

class ASRService:
    """Automatic Speech Recognition (ASR) service using a Hugging Face pipeline."""

    def __init__(self) -> None:
        """Initialize the ASR service with configuration from config.py."""

        self.__device = torch.device("cuda" if torch.cuda.is_available() else "cpu")
        self.__language = config.ASR_LANGUAGE
        self.__model_name = config.ASR_MODEL_NAME
        self.__audio_helper = AudioHelper()
        self.__transcriber = self.__load_model()

    def transcribe(self, file: str) -> str:
        """Transcribe an audio file to text.

        Args:
            file (str): Path to the input audio file.

        Returns:
            str: Transcribed text.

        Raises:
            TranscriptionError: If the file is empty or an error occurs during transcription.
        """
        if self.__audio_helper.is_file_empty(file):
            log.error(f"file {file} is empty or contains only silence")
            # raise TranscriptionError(f"file {file} is empty or contains only silence")

        outfile = self.__audio_helper.convert_audio_to_wav(file)
        log.info(f"Transcribing {outfile}...")
        t0 = time.time()

        try:
            result = self.__transcriber(
                outfile, generate_kwargs={"language": self.__language}
            )
        except Exception as e:
            log.exception(f"Error while transcribing: {e}")
            # raise TranscriptionError(f"Error while transcribing file: {outfile}")

        t1 = time.time()
        log.info(f"Transcription completed in {t1 - t0:.2f} seconds.")
        return result["text"]

    def __load_model(self) -> Pipeline:
        """Load the ASR model.

        Returns:
            Pipeline: Loaded Hugging Face pipeline for transcription.
        """
        model_kwargs = {
            "device_map": "auto",
            "torch_dtype": (
                # if left on auto sets float16 for cpu which results in very slow transcriptions
                torch.float16
                if self.__device.type == "cuda"
                else torch.float32
            ),
        }
        log.info(
            f"Loading Whisper: {self.__model_name} with model kwargs: {model_kwargs} on device: {self.__device}"
        )
        t0 = time.time()
        model = pipeline(
            task="automatic-speech-recognition",
            model=self.__model_name,
            # Makes chunks of audio with length x
            chunk_length_s=30,
            model_kwargs=model_kwargs,
            # Explicitly loading the model onto a device also slows down both on cpu and gpu ???
        )
        t1 = time.time()
        log.info(f"Whisper model loaded in {t1 - t0:.2f} seconds.")
        return model
