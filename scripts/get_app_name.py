#!/bin/env python
import json
data = json.load(open('src/bundler.json'))
print(data['app_name'])
