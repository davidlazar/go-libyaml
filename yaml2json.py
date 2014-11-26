#!/usr/bin/env python
# create test cases for the yaml parser
import json
import yaml
import sys
import os

yamlpath = sys.argv[1]
jsonpath = os.path.splitext(yamlpath)[0] + '.json'

print(jsonpath)

with open(yamlpath, 'r') as yp:
    # BaseLoader because our parser does not yet support
    # constructors or resolvers
    data = yaml.load(yp, Loader=yaml.BaseLoader)

with open(jsonpath, 'w') as jp:
    json.dump(data, jp, sort_keys=True, indent=2)

