import requests
import json
import asyncio
import websockets
from essential_generators import DocumentGenerator
import signal
import sys

base_url = "http://localhost:8000/api"  # replace with your server URL

def getCookies(cookie_jar, domain):
    cookie_dict = cookie_jar.get_dict(domain=domain)
    found = ['%s=%s' % (name, value) for (name, value) in cookie_dict.items()]
    return ';'.join(found)

gen = DocumentGenerator()

# Create the first user
url = base_url+"/register"
headers = {'Content-Type': 'application/json'}
fullname1 = gen.name()
firstname1, lastname1 = fullname1[:fullname1.find(" ")], fullname1[fullname1.find(" ")+1:]
user1 = firstname1 + lastname1 + str(gen.small_int())
data = {"username":user1,"password":"password1", "phone":gen.phone().replace("-",""), "firstname": firstname1, "lastname": lastname1}
print(data)
response = requests.post(url, headers=headers, data=json.dumps(data))
print(response.text)
token1 = json.loads(response.text)["Token"]
id1 = int(json.loads(response.text)["ID"])

# Create the second user
fullname2 = gen.name()
firstname2, lastname2 = fullname2[:fullname2.find(" ")], fullname2[fullname2.find(" ")+1:]
user2 = firstname2 + lastname2 + str(gen.small_int())
data = {"username":user2,"password":"password2", "phone":gen.phone().replace("-",""), "firstname": firstname2, "lastname": lastname2}
print(data)
response = requests.post(url, headers=headers, data=json.dumps(data))
print(response.text)
token2 = json.loads(response.text)["Token"]
id2 = int(json.loads(response.text)["ID"])

# Create a chat between the two users
url = base_url+"/chats"
data = {"people":[id1, id2], "token":token1}  # replace with the actual user IDs
response = requests.post(url, headers=headers, data=json.dumps(data))
print(response.text)
chat_id = int(response.text)

def signal_handler(sig, frame):
    print("DELETE request to user1 with auth")
    url = base_url+f"/users/{id1}"
    headers = {'Content-Type': 'application/json'}
    data = {"token":token1}
    response = requests.delete(url, headers=headers, data=json.dumps(data))
    print(response.text)
    
    print("DELETE request to user2 with auth")
    url = base_url+f"/users/{id2}"
    data = {"token":token2}
    response = requests.delete(url, headers=headers, data=json.dumps(data))
    print(response.text)
    sys.exit(0)

signal.signal(signal.SIGINT, signal_handler)

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
        print(f"{user1} received message '{response_dict['Content']}' from {response_dict['Sender']}")
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
