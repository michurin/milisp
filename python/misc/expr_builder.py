# pylint: disable=C0103,W0212

class E:
    '''
    E is "expression" class. It is demo of approach
    that helps to convert any (you are free to extend the list of supported operations)
    Python expression to MiLisp expression.
    It is slightly oversimplified. In real wold you may
    want to split this class to clear "wrapper" (see __init__) and
    "expression" (see _extended_expr) with common base class.
    '''
    def __init__(self, token):
        self._tokens = (token,)
    @classmethod
    def _extended_expr(cls, *tokens):
        e = cls(None)
        e._tokens = tokens  # slightly hackish, for sure
        return e
    def _expr(self, op_name, *args):
        return type(self)._extended_expr(type(self)(op_name), *args)
    def __add__(self, e):
        return self._expr('+', self, e)
    def __sub__(self, e):
        return self._expr('-', self, e)
    def __mul__(self, e):
        return self._expr('*', self, e)
    def __truediv__(self, e):
        return self._expr('/', self, e)
    def __str__(self):
        if len(self._tokens) == 1:
            t, = self._tokens
            if t is None:
                return '()'
            return str(t)
        return '(' + ' '.join(map(str, self._tokens)) + ')'

def main():
    e = E('X') + (E('Y') - E('Z')) * E('Q') / E(2)
    print(e)  # will print "(+ X (/ (* (- Y Z) Q) 2))"
    e = E('X') + (E('Y') - 'Z') * 'Q' / 2  # accidentally it works too
    print(e)  # will show the same result

if __name__ == '__main__':
    main()
