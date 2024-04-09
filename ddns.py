import requests
import json

def test_config(c):
    if c.get("token") is None or c.get("token") == "":
        print("Invalid value for API Token in config.json")
        exit(1)
        
    if c.get("zone") is None or c.get("zone") == "":
        print("Invalid value for zone in c.json")
        exit(1)

    if c.get("skipUpdate") is None or c.get("skipUpdate") == "":
        print("Invalid value for skipUpdate in c.json")
        exit(1)

def load_config():
    with open("./config.json", "r") as configFile:
        #read config file contents
        content = configFile.read()
        #convert content to json object
        config = json.loads(content)
        test_config(config)
        return config

#load configuration file
config = load_config()

def get_record_list() -> list:
    data = requests.get('https://api.cloudflare.com/client/v4/zones/' + config["zone"] + '/dns_records/',
                       headers={"Authorization": 'Bearer ' + config["token"]}).json()
    a_records = []
    for record in data["result"]:
        if record["type"] != "A":
            continue
        a_records.append(record)

    update_records = []
    for record in a_records:
        if record.get("comment") is not None:
            if config["skipUpdate"] in record["comment"]:
                continue
        update_records.append(record)

    return update_records

def getCurrentIpAddress() -> str:
    url = "https://api.ipify.org"
    response = requests.get(url)
    ip = response.text
    return ip

def update_record(record_id: str, ip: str) -> None:
    api_key = config["token"]
    requests.patch('https://api.cloudflare.com/client/v4/zones/' + config["zone"] + '/dns_records/' + record_id,
                       headers={"Authorization": 'Bearer ' + api_key}, json={"content": ip})

def main():
    ip = getCurrentIpAddress()
    records = get_record_list()
    for record in records:
        update_record(record["id"], ip)

main()
