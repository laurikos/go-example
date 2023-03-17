import requests

def get_joke():
    url = "https://official-joke-api.appspot.com/random_joke"
    response = requests.get(url)
    data = response.json()
    return data

joke_data = get_joke()
print(f"Here's a {joke_data['type']} joke:")
print(joke_data['setup'])
print(joke_data['punchline'])
