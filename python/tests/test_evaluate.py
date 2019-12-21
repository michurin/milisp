import milisp


def test_evaluate_example_1():
    r = milisp.evaluate({
        '+': lambda e, x: milisp.evaluate(e, x[1]) + milisp.evaluate(e, x[2]),
        'x': 1.0,
    }, milisp.parse('(+ x 2)'))
    assert r == 3.


def test_evaluate_example_2():
    r = milisp.evaluate({
        'concat': lambda e, x: milisp.evaluate(e, x[1]) + milisp.evaluate(e, x[2]),
        'x': "Hello",
    }, milisp.parse('(concat x "World")'))
    assert r == 'HelloWorld'
