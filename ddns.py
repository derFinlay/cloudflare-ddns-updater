import requests
import json

def load_config():
    with open("./config.json", "r") as configFile:
        test_config()
        #read config file contents
        content = configFile.read()
        #convert content to json object
        config = json.loads(content)
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

def getCurrentIpAddressCF() -> str:
    #fetch cloudflare API for ip
    url = "https://cloudflare.com/cdn-cgi/trace"
    response = requests.get(url)
    data = response.text
    #only keep the line containing the ip and split it by "=" and take the second part of the split, which is the ip
    ip = [line for line in data.split("\n") if line.startswith("ip=")][0].split("=")[1]
    return ip

def update_record(record_id: str, ip: str) -> None:
    api_key = config["token"]
    requests.patch('https://api.cloudflare.com/client/v4/zones/' + config["zone"] + '/dns_records/' + record_id,
                       headers={"Authorization": 'Bearer ' + api_key}, json={"content": ip})

def test_config():
    if config.get("token") is None or config.get("token") == "":
        print("Invalid value for API Token in config.json")
        exit(1)
        
    if config.get("zone") is None or config.get("zone") == "":
        print("Invalid value for zone in config.json")
        exit(1)

    if config.get("skipUpdate") is None or config.get("skipUpdate") == "":
        print("Invalid value for skipUpdate in config.json")
        exit(1)

    if config.get("records") is None or config.get("records")[0] is None:
        print("Invalid value for records in config.json")
        exit(1)

def main():
    ip = getCurrentIpAddressCF()
    records = get_record_list()
    for record in records:
        update_record(record["id"], ip)

main()