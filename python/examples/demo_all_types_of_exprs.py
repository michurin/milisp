from milisp import evaluate as E, parse as P


def p(env, text, exp):
    res = E(env, P(text))
    print(f'{text} -> {res}')
    assert res == exp  # slightly naive eq check


def op_get_operation_by_name(env, args):
    return env[E(env, args[0])]


def op_sum(env, args):
    return sum(E(env, a) for a in args)


def main():
    p(None, '1', 1)
    p(None, '"ok"', 'ok')
    p(None, '()', None)
    p({'x': Ellipsis}, 'x', Ellipsis)
    p({'+': op_sum}, '(+ 1 2)', 3)
    p({
        'op': op_get_operation_by_name,
        '+': op_sum,
        'x': 1.,
        'y': 2.,
    }, '((op "+") x y)', 3)


if __name__ == '__main__':
    main()
