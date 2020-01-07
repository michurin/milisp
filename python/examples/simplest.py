from milisp import evaluate as E, parse as P


def main():
    def operation_plus(env, args):
        return sum(E(env, a) for a in args)
    root_expression = P('(+ 1 2)')
    env = {'+': operation_plus}
    result = E(env, root_expression)
    assert result == 3
    print(result)


if __name__ == '__main__':
    main()
