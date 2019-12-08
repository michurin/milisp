import milisp as mi

# Implement math variant of operations: vector, and, in

def math_vector_op(env, expr):
    return list(float(mi.evaluate(env, e)) for e in expr[1:])

def math_and_op(env, expr):
    return any(bool(mi.evaluate(env, e)) for e in expr[1:]) # lazy in Python3

def math_in_op(env, expr):
    pattern = mi.evaluate(env, expr[1])
    lst = mi.evaluate(env, expr[2])
    return pattern in lst

math_ops = {
    'vector': math_vector_op,
    'and': math_and_op,
    'in': math_in_op,
}

# Implement SQL variant of the same operations

def sql_vector_op(env, expr):
    return 'select\n' + ',\n'.join('  ' + mi.evaluate(env, e) for e in expr[1:]) + '\nfrom\n  hive.events as log;'

def sql_and_op(env, expr):
    return ' and '.join(mi.evaluate(env, e) for e in expr[1:])

def sql_in_op(env, expr):
    return '(' + mi.evaluate(env, expr[1]) + ' in (' + ', '.join(map(repr, mi.evaluate(env, expr[2]))) + '))'

sql_ops = {
    'vector': sql_vector_op,
    'and': sql_and_op,
    'in': sql_in_op,
}

# Constants

consts = {
    'RU': ['+7'],
    'UK': ['+44'],
    'IL': ['+972'],
    'MSK': ['095', '495'],
    'LDN': ['020'],
    'TLV': ['3'],
}

# Text

text = '''
(vector
    (and (in phoneCountryCode UK) (in phoneAreaCode LDN))
    (and (in phoneCountryCode IL) (in phoneAreaCode TLV))
    (and (in phoneCountryCode RU) (in phoneAreaCode MSK))
)
'''

def test_calc():
    prog = mi.parse(text)
    data = {
        'phoneCountryCode': '+972',
        'phoneAreaCode': '3',
    }
    res = mi.evaluate({**math_ops, **consts, **data}, prog)
    assert res == [0.0, 1.0, 0.0]

def test_sql():
    prog = mi.parse(text)
    data = {
        'phoneCountryCode': 'log.phone_country',
        'phoneAreaCode': 'log.phone_area',
    }
    res = mi.evaluate({**sql_ops, **consts, **data}, prog)
    assert res == '''\
select
  (log.phone_country in ('+44')) and (log.phone_area in ('020')),
  (log.phone_country in ('+972')) and (log.phone_area in ('3')),
  (log.phone_country in ('+7')) and (log.phone_area in ('095', '495'))
from
  hive.events as log;'''
