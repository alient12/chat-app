import requests
import json
import asyncio
import websockets
from essential_generators import DocumentGenerator
import signal
import sys

base_url = "http://localhost:8000/api"  # replace with your server URL


gen = DocumentGenerator()

# Ali
token1 = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MTg0MTU1NTQzNCwidXNlcm5hbWUiOiJhbGllbnQxMiIsImV4cCI6MTcwOTMwNjYyNX0.ql2O-ao6n7y1c5dtZPPAjqFhWwmXpFfz2oDE4IPnyBA"
id1 = 1841555434

# Javad
token2 = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MTI0ODIwOTE0LCJ1c2VybmFtZSI6InBoeWpheSIsImV4cCI6MTcwOTMwNzAyMX0.gc0_Hw5nHAg6RGNrAAkLykHZvRbM90L3VeRig7czAbE"
id2 = 124820914

chat_id = 179190051



# WebSocket URL
ws_url = "ws://localhost:8000/api/message"  # replace with your WebSocket URL

async def send_message(user1, user2, ws, delay, initial_delay=0):
    await asyncio.sleep(initial_delay)

    while True:
        # Wait for the specified delay before sending the next message
        await asyncio.sleep(delay)

        # Send a message
        msg = {"chatid": chat_id, "receiver": user2, "content": gen.sentence(), "contenttype": 0}  # replace with your message
        await ws.send(json.dumps(msg))
        
        print(f"I'm {user1} and sent a message")

async def receive_message(user1, ws):
    while True:
        # Receive a message
        response = await ws.recv()
        response_dict = json.loads(response)
        print(f"{user1} received message '{response_dict['content']}' from {response_dict['sender']}")
        print()

async def handle_connection(user1, user2, token, delay, initial_delay=0):
    async with websockets.connect(ws_url+"?token="+token) as ws:
        print(f"I'm {user1} and got connected to websocket.")
        # Run the send_message and receive_message functions concurrently
        await asyncio.gather(
            send_message(user1, user2, ws, delay, initial_delay),
            receive_message(user1, ws)
        )

async def get_chat(chat_id, token, delay, initial_delay=0):
    await asyncio.sleep(initial_delay)

    while True:
        # Wait for the specified delay before sending the next message
        await asyncio.sleep(delay)

        print(f"GET request to /chats/{chat_id}")
        
        url = base_url+f"/chats/{chat_id}"
        params = {"token":token}
        response = requests.get(url, params=params)
                
        print(response.text)

# Run the handle_connection function for both users concurrently with different delays
asyncio.get_event_loop().run_until_complete(asyncio.gather(
    handle_connection(id1, id2, token1, 2, 1),
    handle_connection(id2, id1, token2, 2),
    get_chat(chat_id, token1, 5)
))
