#!/usr/bin/python
import argparse
import requests
import json

# Exports contract ABI in JSON

parser = argparse.ArgumentParser()
parser.add_argument('chain', type=str, help='chain')
parser.add_argument('addr', type=str, help='Contract address')
parser.add_argument('-o', '--output', type=str,
                    help="Path to the output JSON file", required=True)


def endpoint(chain):
    if chain == "avax":
        return "https://api.snowtrace.io/api?module=contract&action=getabi&address="
    elif chain == "polygon":
        return "https://api.polygonscan.com/api?apikey=11PJRTF1A7HNW2BRG35274WKP9TEEFBWBX&module=contract&action=getabi&address="
    else:
        return "https://api.etherscan.io/api?module=contract&action=getabi&address="


def __main__():

    args = parser.parse_args()

    url = '%s%s' % (endpoint(args.chain), args.addr)
    print(url)
    response = requests.get(url)
    response_json = response.json()
    abi_json = json.loads(response_json['result'])
    result = json.dumps({"abi": abi_json}, indent=4, sort_keys=True)

    open(args.output, 'w').write(result)


if __name__ == '__main__':
    __main__()
