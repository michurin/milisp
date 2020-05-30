from milisp import evaluate as E, parse as P


def main():
    # expression
    lisp_code = '(+ 1 2)'
    # setup env: define '+' operatoin
    def operation_plus(env, args):
        return sum(E(env, a) for a in args)
    env = {'+': operation_plus}
    # compile expression
    root_expression = P(lisp_code)
    # execute compiled expression
    result = E(env, root_expression)
    # check
    assert result == 3
    print(result)


if __name__ == '__main__':
    main()
