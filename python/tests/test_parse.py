import pytest

import milisp


@pytest.mark.parametrize('p', ('', '(a', '(a))', ')'))
def test_invalid_progs(p):
    with pytest.raises(milisp.LispError):
        list(milisp.parse(p))
