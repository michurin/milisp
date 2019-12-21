import milisp as mi


def prog_op(env, expr):
    for e in expr[1:]:
        r = mi.evaluate(env, e)
    return r


def mul_op(env, expr):
    return mi.evaluate(env, expr[1]) * mi.evaluate(env, expr[2])


def test_factorial_loop():
    text = '''
    (prog                     # execute all floowing expressions and return result of last
        (set "x" 1)           # x = 1
        (loop "i" 1 N         # for i = 1; i <= N; i++
            (set "x" (* x i)) # x = x * i
        )
        x                     # return x
    )
    '''

    def set_op(env, expr):
        env[mi.evaluate(env, expr[1])] = mi.evaluate(env, expr[2])

    def loop_op(env, expr):
        var_name = mi.evaluate(env, expr[1])
        first = mi.evaluate(env, expr[2])
        last = mi.evaluate(env, expr[3])
        for i in range(int(first), int(last)+1):
            env[var_name] = float(i)
            mi.evaluate(env, expr[4])
    env = {
        'prog': prog_op,
        'set': set_op,
        'loop': loop_op,
        '*': mul_op,
        'N': 5.,
    }
    prog = mi.parse(text)
    res = mi.evaluate(env, prog)
    assert res == 120.


def test_factorial_recursion():
    text = '''
    (prog
        (def "F" "x" (if_gt_one   # if x > 1 then F(x-1) else 1
            x
            (* x (call "F" (+ x -1)))
            1
        ))
        (call "F" N)
    )
    '''

    def def_op(env, expr):
        env[mi.evaluate(env, expr[1])] = (
            mi.evaluate(env, expr[2]),
            expr[3]
        )

    def call_op(env, expr):
        argname, funcbody = env[mi.evaluate(env, expr[1])]
        localenv = env.copy()  # shallow copy if enough for our purpose
        localenv[argname] = mi.evaluate(env, expr[2])
        return mi.evaluate(localenv, funcbody)

    def if_gt_one_op(env, expr):
        if mi.evaluate(env, expr[1]) > 1.:
            return mi.evaluate(env, expr[2])
        return mi.evaluate(env, expr[3])

    def plus_op(env, expr):
        return mi.evaluate(env, expr[1]) + mi.evaluate(env, expr[2])
    env = {
        'prog': prog_op,
        '*': mul_op,
        '+': plus_op,
        'def': def_op,
        'call': call_op,
        'if_gt_one': if_gt_one_op,
        'N': 5.,
    }
    prog = mi.parse(text)
    res = mi.evaluate(env, prog)
    assert res == 120.
