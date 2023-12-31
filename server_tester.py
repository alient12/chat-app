import requests
import json
import time

base_url = "http://127.0.0.1:8000/api"

print("POST request to /register")
url = base_url+"/register"
headers = {'Content-Type': 'application/json'}
data = {"username":"user1","password":"password1", "phone":"0123456789", "firstname": "Ali", "lastname": "Ent"}
response = requests.post(url, headers=headers, data=json.dumps(data))
print(response.text)
cookies = response.cookies
print(cookies.get_dict)

print("GET request to /users/id")
id = int(response.text)
url = base_url+f"/users/{id}"
response = requests.get(url, data=json.dumps(data))
print(response.text)

print("PATCH request to /users/id without auth")
url = base_url+f"/users/{id}"
headers = {'Content-Type': 'application/json'}
data = {"username":"user1","password":"password1", "phone":"0123456789", "firstname": "Ali", "lastname": "Entwdwd"}
response = requests.patch(url, headers=headers, data=json.dumps(data))
print(response.text)

print("GET request to /users/id")
url = base_url+f"/users/{id}"
response = requests.get(url, data=json.dumps(data))
print(response.text)

print("PATCH request to /users/id with auth")
url = base_url+f"/users/{id}"
headers = {'Content-Type': 'application/json'}
data = {"username":"user1","password":"password1", "phone":"0123456789", "firstname": "Ali", "lastname": "Entwdwd"}
response = requests.patch(url, headers=headers, data=json.dumps(data), cookies=cookies)
print(response.text)

print("GET request to /users/id")
url = base_url+f"/users/{id}"
response = requests.get(url, data=json.dumps(data))
print(response.text)

print("POST request to /register")
url = base_url+"/register"
headers = {'Content-Type': 'application/json'}
data = {"username":"user2","password":"password2", "phone":"0123456788", "firstname": "Javad", "lastname": "Zein"}
response = requests.post(url, headers=headers, data=json.dumps(data))
print(response.text)
cookies2 = response.cookies
id2 = int(response.text)
print(cookies.get_dict)

print("PATCH request to user1 with auth user2")
url = base_url+f"/users/{id}"
headers = {'Content-Type': 'application/json'}
data = {"username":"user1","password":"password1", "phone":"0123456789", "firstname": "Ali", "lastname": "Entwdwd"}
response = requests.patch(url, headers=headers, data=json.dumps(data), cookies=cookies2)
print(response.text)

print("DELETE request to user2 without auth")
url = base_url+f"/users/{id2}"
response = requests.delete(url)
print(response.text)

print("DELETE request to user2 with auth user1")
url = base_url+f"/users/{id2}"
response = requests.delete(url, cookies=cookies)
print(response.text)

print("DELETE request to user2 with auth")
url = base_url+f"/users/{id2}"
response = requests.delete(url, cookies=cookies2)
print(response.text)

print("DELETE request to user1 with auth")
url = base_url+f"/users/{id}"
response = requests.delete(url, cookies=cookies)
print(response.text)

