import milisp as mi


def prog_op(env, args):
    for a in args:
        r = mi.evaluate(env, a)
    return r


def mul_op(env, args):
    return mi.evaluate(env, args[0]) * mi.evaluate(env, args[1])


def set_op(env, args):
    env[mi.evaluate(env, args[0])] = mi.evaluate(env, args[1])


def loop_op(env, args):
    var_name = mi.evaluate(env, args[0])
    first = mi.evaluate(env, args[1])
    last = mi.evaluate(env, args[2])
    for i in range(int(first), int(last)+1):
        env[var_name] = float(i)
        mi.evaluate(env, args[3])


def main():
    text = """
    (prog                     # execute all following expressions and return result of last
        (set "x" 1)           # x = 1
        (loop "i" 1 N         # for i = 1; i <= N; i++
            (set "x" (* x i)) # x = x * i
        )
        x                     # return x
    )
    """
    env = {
        'prog': prog_op,
        'set': set_op,
        'loop': loop_op,
        '*': mul_op,
        'N': 5.,
    }
    expr = mi.parse(text)
    res = mi.evaluate(env, expr)
    assert res == 120.
    print(res)


if __name__ == '__main__':
    main()
