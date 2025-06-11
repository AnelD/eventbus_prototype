import asr_service
from websocket_client import WebSocketClient
from logger_helper import LoggerHelper
import asyncio

log = LoggerHelper(__name__).get_logger()

async def main():
    asr = asr_service.ASRService()
    queue = asyncio.Queue()
    client = WebSocketClient("ws://localhost:8080/ws", queue)
    await client.connect()
    await client.send_message({'type': 'subscribe', 'topic': 'file.upload'})
    res = asr.transcribe(await queue.get())
    await client.send_message({'type': 'publish', 'topic': 'audio.transcript', 'data': res})
    print(res)


if __name__ == "__main__":
    asyncio.run(main())
