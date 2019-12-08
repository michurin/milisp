import pytest
import milisp

def test_empty():
    with pytest.raises(milisp.LispError):
        list(milisp.parse(''))

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
