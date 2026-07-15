import sys, json, urllib.request

key = sys.argv[1]
url = f"https://generativelanguage.googleapis.com/v1beta/models?key={key}"
all_models = []
while url:
    req = urllib.request.Request(url)
    with urllib.request.urlopen(req) as response:
        data = json.loads(response.read().decode())
        all_models.extend([m["name"] for m in data.get("models", [])])
        pt = data.get("nextPageToken")
        if pt:
            url = f"https://generativelanguage.googleapis.com/v1beta/models?key={key}&pageToken={pt}"
        else:
            url = None

print("\n".join(m for m in all_models if "flash" in m.lower()))
