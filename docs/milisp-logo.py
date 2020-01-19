#!/usr/bin/python3

import base64
import re
import urllib.request

def replace_font(m):
    url = m.group(1)
    with urllib.request.urlopen(url) as f:
        response = f.read()
    data = base64.b64encode(response)
    return 'url("data:font/woff2;base64,' + data.decode('U8') + '")'

def replace_import(m):
    url = m.group(1)
    with urllib.request.urlopen(url) as f:
        response = f.read()
    body = response.decode('U8')
    return re.sub(r'url\(([^)]+)\)', replace_font, body)

def main(in_file, out_file):
    with open(in_file) as r, open(out_file, 'w') as w:
        orig = r.read()
        res = re.sub(r'@import\s+url\("([^"]+)"\);?', replace_import, orig)
        w.write(res)

if __name__ == '__main__':
    main('milisp-logo.svg', 'milisp-logo-solid.svg')
