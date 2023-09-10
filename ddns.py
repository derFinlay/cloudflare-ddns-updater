import requests
import json

def loadConfig():
    with open("./config.json", "r") as configFile:
        #read config file contents
        content = configFile.read()
        #convert content to json object
        config = json.loads(content)
        return config

#load configuration file
config = loadConfig()

def getCurrentIpAddressCF():
    #fetch cloudflare API for ip
    url = "https://cloudflare.com/cdn-cgi/trace"
    response = requests.get(url)
    data = response.text
    #only keep the line containing the ip and split it by "=" and take the second part of the split, which is the ip
    ip = [line for line in data.split("\n") if line.startswith("ip=")][0].split("=")[1]
    return ip


def getCurrentIpAddress():
    ip = requests.get('https://api.ipify.org/').text
    return ip


def updateRecordIP(ip, zoneId, recordId):
    auth_token = config["token"]
    headers = {"Authorization": f"Bearer {auth_token}"}
    requests.patch('https://api.cloudflare.com/client/v4/zones/' + zoneId + '/dns_records/' + recordId,
                    headers=headers, json={"content": ip})


def getRecordIp(zoneId, recordId):
    res = requests.get('https://api.cloudflare.com/client/v4/zones/' + zoneId + '/dns_records/' + recordId,
                       headers={"Authorization": 'Bearer Yo_2MDoOoMUYHIkVMiqBMfMpInSX1_bbOZzaTu5T'})
    return res.json()["result"]["content"]


def main():
    ip = getCurrentIpAddressCF()
    if ip != getRecordIp(config["zone"], config["records"][0]):
        for record in config["records"]:
            updateRecordIP(ip, config["zone"], record)
        print('Updated!')
    else:
        print('Not updated!')

main()