import requests
import json
import time

base_url = "http://127.0.0.1:8000/api"


########################### register and update users ###########################
print("########################### register and update users ###########################")
print("POST request to /register")
url = base_url+"/register"
headers = {'Content-Type': 'application/json'}
data = {"username":"user1","password":"password1", "phone":"0123456789", "firstname": "Ali", "lastname": "Ent"}
response = requests.post(url, headers=headers, data=json.dumps(data))
print(response.text)
token = json.loads(response.text)["Token"]

print("GET request to /users/id")
id = int(json.loads(response.text)["ID"])
url = base_url+f"/users/{id}"
response = requests.get(url, headers=headers, data=json.dumps(data))
print(response.text)

print("PATCH request to /users/id without auth")
url = base_url+f"/users/{id}"
headers = {'Content-Type': 'application/json'}
data = {"username":"user1","password":"password1", "phone":"0123456789", "firstname": "Ali", "lastname": "Entwdwd"}
response = requests.patch(url, headers=headers, data=json.dumps(data))
print(response.text)

print("GET request to /users/id")
url = base_url+f"/users/{id}"
response = requests.get(url, headers=headers, data=json.dumps(data))
print(response.text)

print("PATCH request to /users/id with auth")
url = base_url+f"/users/{id}"
headers = {'Content-Type': 'application/json'}
data = {"username":"user1","password":"password1", "phone":"0123456789", "firstname": "Ali", "lastname": "Entwdwd", "token":token}
response = requests.patch(url, headers=headers, data=json.dumps(data))
print(response.text)

print("GET request to /users/id")
url = base_url+f"/users/{id}"
response = requests.get(url, headers=headers, data=json.dumps(data))
print(response.text)

print("POST request to /register")
url = base_url+"/register"
headers = {'Content-Type': 'application/json'}
data = {"username":"user2","password":"password2", "phone":"0123456788", "firstname": "Javad", "lastname": "Zein"}
response = requests.post(url, headers=headers, data=json.dumps(data))
print(response.text)
token2 = json.loads(response.text)["Token"]
id2 = int(json.loads(response.text)["ID"])

print("PATCH request to user1 with auth user2")
url = base_url+f"/users/{id}"
headers = {'Content-Type': 'application/json'}
data = {"username":"user1","password":"password1", "phone":"0123456789", "firstname": "Ali", "lastname": "Entwdwd", "token":token2}
response = requests.patch(url, headers=headers, data=json.dumps(data))
print(response.text)

print("GET request to /users?keyword=user")
url = base_url+f'/users?keyword=user'
response = requests.get(url, headers=headers, data=json.dumps(data))
print(response.text)


########################### Chat tests ###########################
print("########################### Chat tests ###########################")
print("POST request to /chats with user1")
url = base_url+"/chats"
headers = {'Content-Type': 'application/json'}
data = {"people":[id, id2], "token":token}
response = requests.post(url, headers=headers, data=json.dumps(data))
print(response.text)

print("GET request to /chats")
url = base_url+f"/chats"
params = {"token":token}
response = requests.get(url, params=params)
print(response.text)

print("POST request to /chats with user2")
url = base_url+"/chats"
headers = {'Content-Type': 'application/json'}
data = {"people":[id2, id], "token":token2}
response = requests.post(url, headers=headers, data=json.dumps(data))
print(response.text)

print("GET request to /chats")
url = base_url+f"/chats"
params = {"token":token2}
response = requests.get(url, params=params)
print(response.text)

########################### Contact tests ###########################
print("########################### Contact tests ###########################")
print("POST request to create contact with user1")
url = base_url+f"/users/{id}/contacts"
headers = {'Content-Type': 'application/json'}
data = {"id":id2, "name":"jzein","token":token}
response = requests.post(url, headers=headers, data=json.dumps(data))
print(response.text)

print("GET request to get user1 contacts")
url = base_url+f"/users/{id}/contacts"
params = {"token":token}
response = requests.get(url, params=params)
print(response.text)

print("GET request to get user1 contacts with user2")
url = base_url+f"/users/{id}/contacts"
params = {"token":token2}
response = requests.get(url, params=params)
print(response.text)

print("DELETE request contact user1")
url = base_url+f"/users/{id}/contacts/{id2}"
data = {"token":token}
response = requests.delete(url, headers=headers, data=json.dumps(data))
print(response.text)

print("GET request to get user1 contacts")
url = base_url+f"/users/{id}/contacts"
params = {"token":token}
response = requests.get(url, params=params)
print(response.text)


########################### Delete users ###########################
print("########################### Delete users ###########################")
print("DELETE request to user2 without auth")
url = base_url+f"/users/{id2}"
response = requests.delete(url)
print(response.text)

print("DELETE request to user2 with auth user1")
url = base_url+f"/users/{id2}"
data = {"token":token}
response = requests.delete(url, headers=headers, data=json.dumps(data))
print(response.text)

print("DELETE request to user2 with auth")
url = base_url+f"/users/{id2}"
data = {"token":token2}
response = requests.delete(url, headers=headers, data=json.dumps(data))
print(response.text)

print("DELETE request to user1 with auth")
url = base_url+f"/users/{id}"
data = {"token":token}
response = requests.delete(url, headers=headers, data=json.dumps(data))
print(response.text)

