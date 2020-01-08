import milisp as mi


def prog_op(env, args):
    for a in args:
        r = mi.evaluate(env, a)
    return r


def mul_op(env, args):
    return mi.evaluate(env, args[0]) * mi.evaluate(env, args[1])


def def_op(env, args):
    env[mi.evaluate(env, args[0])] = (mi.evaluate(env, args[1]), args[2])


def call_op(env, args):
    argname, funcbody = env[mi.evaluate(env, args[0])]
    localenv = env.copy()  # shallow copy if enough for our purpose
    localenv[argname] = mi.evaluate(env, args[1])
    return mi.evaluate(localenv, funcbody)


def if_gt_one_op(env, args):
    if mi.evaluate(env, args[0]) > 1.:
        return mi.evaluate(env, args[1])
    return mi.evaluate(env, args[2])


def plus_op(env, args):
    return mi.evaluate(env, args[0]) + mi.evaluate(env, args[1])


def main():
    text = """
    (prog
        (def "F" "x" (if_gt_one   # if x > 1 then F(x-1) else 1
            x
            (* x (call "F" (+ x -1)))
            1
        ))
        (call "F" N)
    )
    """
    env = {
        'prog': prog_op,
        '*': mul_op,
        '+': plus_op,
        'def': def_op,
        'call': call_op,
        'if_gt_one': if_gt_one_op,
        'N': 5.,
    }
    expr = mi.parse(text)
    res = mi.evaluate(env, expr)
    assert res == 120.
    print(res)


if __name__ == '__main__':
    main()
