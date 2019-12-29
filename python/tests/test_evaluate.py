import milisp


def test_evaluate_number():
    r = milisp.evaluate(None, milisp.parse('1'))
    assert r == 1


def test_evaluate_string():
    r = milisp.evaluate(None, milisp.parse('"ok"'))
    assert r == 'ok'


def test_evaluate_none():
    r = milisp.evaluate(None, milisp.parse('()'))
    assert r is None


def test_evaluate_var():
    r = milisp.evaluate({'x': 1}, milisp.parse('x'))
    assert r == 1


def test_evaluate_simple_expression():
    expr = milisp.parse('(+ 1 2)')
    env = {'+': lambda e, a: milisp.evaluate(e, a[0]) + milisp.evaluate(e, a[1])}
    r = milisp.evaluate(env, expr)
    assert r == 3


def test_evaluate_simple_with_strings():
    expr = milisp.parse('(concat x "World")')
    env = {
        'concat': lambda e, a: milisp.evaluate(e, a[0]) + milisp.evaluate(e, a[1]),
        'x': 'Hello',
    }
    r = milisp.evaluate(env, expr)
    assert r == 'HelloWorld'


def test_evaluate_nontrivial():
    expr = milisp.parse('((op "+") x 2)')
    env = {
        '+': lambda e, a: milisp.evaluate(e, a[0]) + milisp.evaluate(e, a[1]),
        'op': lambda e, a: e[milisp.evaluate(e, a[0])],
        'x': 1,
    }
    r = milisp.evaluate(env, expr)
    assert r == 3
