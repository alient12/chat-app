import requests
import json
import asyncio
import websockets
from essential_generators import DocumentGenerator

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
cookies1 = response.cookies
id1 = int(json.loads(response.text)["ID"])

# Create the second user
fullname2 = gen.name()
firstname2, lastname2 = fullname2[:fullname2.find(" ")], fullname2[fullname2.find(" ")+1:]
user2 = firstname2 + lastname2 + str(gen.small_int())
data = {"username":user2,"password":"password2", "phone":gen.phone().replace("-",""), "firstname": firstname2, "lastname": lastname2}
print(data)
response = requests.post(url, headers=headers, data=json.dumps(data))
print(response.text)
cookies2 = response.cookies
id2 = int(json.loads(response.text)["ID"])

# Create a chat between the two users
url = base_url+"/chats"
data = {"people":[id1, id2]}  # replace with the actual user IDs
response = requests.post(url, headers=headers, data=json.dumps(data), cookies=cookies1)
print(response.text)
chat_id = int(response.text)


url = base_url+f'/files/upload/{chat_id}'  # replace with your URL
file_path = 'Dance.mp4'  # replace with your file path

with open(file_path, 'rb') as f:
    files = {'file': f}
    response = requests.post(url, files=files, cookies=cookies1)
    print(response.text)