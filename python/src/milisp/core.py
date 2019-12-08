import re

WHITESPACE = re.compile(r'[ \f\t\r\n]*')
PSEUDOTOKEN = re.compile(r'(?:#[^\n]*(?:\n|$)|\(|\)|"[^"\\]*(?:\\.[^"\\]*)*"|[^ "()\\#\f\t\r\n]+)')
ESCAPED = re.compile(r'\\.')

Symbol = type('Symbol', (str,), {})
BEGIN = type('Begin', (object,), {'__repr__': lambda x: 'Begin'})()
END = type('End', (object,), {'__repr__': lambda x: 'End'})()

class LispError(Exception): pass

def tokenize(text):
    pos = 0
    while True:
        pos = WHITESPACE.match(text, pos).end()
        if pos >= len(text):
            break
        m = PSEUDOTOKEN.match(text, pos)
        if m is None:
            raise LispError('Symbols ...%s' % text[pos:])
        t = m.group()
        pos = m.end()
        f = t[0]
        r = None
        if f == '(':
            r = BEGIN
        elif f == ')':
            r = END
        elif f == '"':
            r = ESCAPED.sub(lambda x: x.group()[1], t[1:-1])
        elif f != '#':
            try:
                r = float(t)
            except ValueError:
                r = Symbol(t)
        if r is not None:
            yield r

BREAK = object()

def parse(text):
    tz = iter(tokenize(text))
    ast = None
    try:
        ast = recursive_descent(tz)
        next(tz) # it has to raise StopIteration
    except StopIteration:
        if ast is None:
            raise LispError('Unexpected end of file')
        if ast is BREAK:
            raise LispError('Extra ")"')
        return ast
    raise LispError('File too long')

def recursive_descent(tz):
    ch = next(tz)
    if ch is BEGIN:
        c = []
        while True:
            x = recursive_descent(tz)
            if x is BREAK:
                return c
            c.append(x)
    elif ch is END:
        return BREAK
    else:
        return ch

def evaluate(env, ast):
    if type(ast) is list:
        v = env[ast[0]]
        return v(env, ast)
    if type(ast) is Symbol:
        return env[ast]
    return ast # type(ast) in (float, str)
