import pytest
import milisp

def test_empty():
    with pytest.raises(milisp.LispError):
        list(milisp.parse(''))

def test_empty_expr():
    with pytest.raises(milisp.LispError):
        list(milisp.parse('()'))

def test_a():
    t = list(milisp.parse('(a)'))
    assert t == ['a']

def test_noclosed():
    with pytest.raises(milisp.LispError):
        list(milisp.parse('(a'))

def test_overclosed():
    with pytest.raises(milisp.LispError):
        list(milisp.parse('(a))'))

def test_invalidclosed():
    with pytest.raises(milisp.LispError):
        list(milisp.parse(')'))

def test_invalid_first_token():
    with pytest.raises(milisp.LispError):
        list(milisp.parse('(1)'))
