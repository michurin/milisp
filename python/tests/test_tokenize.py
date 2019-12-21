import pytest
import milisp.core as core


def test_empty():
    t = list(core.tokenize(''))
    assert t == []


def test_tokens():
    t = list(core.tokenize(r'(xy 1. "ab\c\"")'))
    assert t == [core.BEGIN, 'xy', 1., 'abc"', core.END]
    assert [type(x) for x in t] == [
        type(core.BEGIN),
        core.Symbol,
        float,
        str,
        type(core.END),
    ]


def test_empty_comment():
    t = list(core.tokenize('#\nA'))
    assert t == ['A']


def test_comment():
    t = list(core.tokenize('#x\nA'))
    assert t == ['A']


def test_invalid_token():
    with pytest.raises(core.LispError):
        list(core.tokenize('A\\'))
