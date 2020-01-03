import io

import numpy as np

import milisp as mi

# This 3-in-1 example shows, how to use one common set of lisp expression (Ch.0) to:
# - calculate vector with bare Python (Ch.1)
# - generate SQL request (Ch.2)
# - perfrom bulk calculations using NumPy (Ch.3)

# ---- CHAPTER 0 ---------------------------------------------
# Lisp code. The same for all examples
# There are two parts of Lisp code:
# - CODE_INIT - just for initializatin constants in enviroment
# - CODE_CALCULATE - features expressions itselfs

CODE_INIT = """
(exec
    # country codes
    (set "IL" "+972")
    (set "RU" "+7")
    (set "UK" "+44")
    # area codes
    (set "LDN" "020")
    (set "MSK" "095" "495")
    (set "TLV" "3")
)
"""

CODE_CALCULATE = """
(vector
    (and (in phoneCountryCode UK) (in phoneAreaCode LDN))
    (and (in phoneCountryCode IL) (in phoneAreaCode TLV))
    (and (in phoneCountryCode RU) (in phoneAreaCode MSK))
)
"""

# ---- CHAPTER 1 ---------------------------------------------
# Pure Python implementation. It calculates one vector per run.
# Interesting points:
# - lazy "and" operation
# - we can parse code once and reuse the result


def pure_exec_op(env, args):
    [mi.evaluate(env, a) for a in args]


def pure_set_op(env, args):
    vals = [mi.evaluate(env, a) for a in args]
    env[vals[0]] = vals[1:]


def pure_vector_op(env, args):
    return [float(mi.evaluate(env, a)) for a in args]


def pure_and_op(env, args):
    return any(bool(mi.evaluate(env, a)) for a in args)  # lazy in Python3


def pure_in_op(env, args):
    pattern = mi.evaluate(env, args[0])
    lst = mi.evaluate(env, args[1])
    return pattern in lst  # by the way, "in" could be lazy too in more complex circumstance


def main_calc_one_vector_with_pure_python():
    env = {  # setup operations
        'vector': pure_vector_op,
        'and': pure_and_op,
        'in': pure_in_op,
        'exec': pure_exec_op,
        'set': pure_set_op,
    }
    mi.evaluate(env, mi.parse(CODE_INIT))  # run init code to setup constants
    prog = mi.parse(CODE_CALCULATE)  # prepere features calculator
    env.update({  # one raw data vector
        'phoneCountryCode': '+972',
        'phoneAreaCode': '3',
    })
    res = mi.evaluate(env, prog)  # calculate features
    assert res == [0.0, 1.0, 0.0]
    print(res)

# ---- CHAPTER 2 ---------------------------------------------
# Implement SQL variant of the same operations
# Interesting points:
# - operation as a class with an internal state. It could be used for caching, debugging, monitoring and instrumentation...
# - features are used to pass column names


class SQLVectorOp:

    def __init__(self, table_name, table_alias):
        self.table_name = table_name
        self.table_alias = table_alias

    def __call__(self, env, args):
        return ''.join((
            'select\n',
            ',\n'.join('  ' + mi.evaluate(env, a) for a in args),
            '\nfrom\n  ',
            self.table_name,
            ' as ',
            self.table_alias,
            ';'))


def sql_and_op(env, args):
    return ' and '.join(mi.evaluate(env, a) for a in args)


def sql_in_op(env, args):
    return ''.join((
        '(',
        mi.evaluate(env, args[0]),
        ' in (',
        ', '.join(map(repr, mi.evaluate(env, args[1]))),
        '))'))


def main_prepare_sql_request():
    table_alias = 'log'
    env = {
        'vector': SQLVectorOp('hive."default".events', table_alias),
        'and': sql_and_op,
        'in': sql_in_op,
        'exec': pure_exec_op,
        'set': pure_set_op,
    }
    mi.evaluate(env, mi.parse(CODE_INIT))
    prog = mi.parse(CODE_CALCULATE)
    env.update({
        'phoneCountryCode': f'{table_alias}.phone_country',
        'phoneAreaCode': f'{table_alias}.phone_area',
    })
    res = mi.evaluate(env, prog)
    assert res == """\
select
  (log.phone_country in ('+44')) and (log.phone_area in ('020')),
  (log.phone_country in ('+972')) and (log.phone_area in ('3')),
  (log.phone_country in ('+7')) and (log.phone_area in ('095', '495'))
from
  hive."default".events as log;"""
    print(res)


# ---- CHAPTER 3 ---------------------------------------------
# Use NumPy
# Interesting points:
# - we are free to store np.arrays as variables values
# - we are free to perform NumPy operations without any limitations


def np_set_op(env, args):
    vals = [mi.evaluate(env, a) for a in args]
    env[vals[0]] = np.array(vals[1:])[:, np.newaxis]


def np_vector_op(env, args):
    return np.stack([mi.evaluate(env, a) for a in args]).T.astype(np.float32)


def np_and_op(env, args):
    return np.stack([mi.evaluate(env, a) for a in args]).all(axis=0)


def np_in_op(env, args):
    pattern = mi.evaluate(env, args[0])
    lst = mi.evaluate(env, args[1])
    return (lst == pattern).any(axis=0)


def main_calc_using_numpy():
    text = """\
            CountryCode   AreaCode
                   +972          3
                     +7        095
                    +44        020
                    +44        023
                    +34        976
    """
    data = np.genfromtxt(io.StringIO(text), skip_header=1, dtype='|U4', autostrip=True)
    env = {
        'vector': np_vector_op,
        'and': np_and_op,
        'in': np_in_op,
        'exec': pure_exec_op,
        'set': np_set_op,
    }
    mi.evaluate(env, mi.parse(CODE_INIT))
    env.update(zip(('phoneCountryCode', 'phoneAreaCode'), data.T))
    res = mi.evaluate(env, mi.parse(CODE_CALCULATE))  # process all data in one run using NumPy
    np.testing.assert_equal(res, [
        # LDN / TLV / MSK
        [0., 1., 0.],
        [0., 0., 1.],
        [1., 0., 0.],
        [0., 0., 0.],
        [0., 0., 0.]])
    print(res)

# ---- M A I N -----------------------------------------------


def main():
    for i, f in enumerate((main_calc_one_vector_with_pure_python, main_prepare_sql_request, main_calc_using_numpy), start=1):
        print(f'\n---- Ch.{i} ----------------')
        f()


if __name__ == '__main__':
    main()
