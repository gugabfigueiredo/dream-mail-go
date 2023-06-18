#!/usr/bin/env python3
import requests

url = 'http://localhost:8080/dmail/send'
mail = {
    'from': {'addr': 'gugabfigueiredo@gmail.com'},
    'to': [{'addr': 'gugabfigueiredo@gmail.com'}],
    'subject': 'Hello, World!',
    'text': 'Hello, World!',
    'html': '<strong>Hello, World!</strong>'
}

x = requests.post(url, json=mail)

print(x.text)
