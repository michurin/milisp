from milisp import evaluate as E, parse as P


def p(env, text):
    res = E(env, P(text))
    print(f'{text} -> {res}')


def op_get_operation_by_name(env, args):
    return env[E(env, args[0])]


def op_sum(env, args):
    return sum(E(env, a) for a in args)


def main():
    p(None, '1')
    p(None, '"ok"')
    p(None, '()')
    p({'x': Ellipsis}, 'x')
    p({'+': op_sum}, '(+ 1 2)')
    p({
        'op': op_get_operation_by_name,
        '+': op_sum,
        'x': 1.,
        'y': 2.,
    }, '((op "+") x y)')


if __name__ == '__main__':
    main()
