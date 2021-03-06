"""Just shallow check that examples are runable. Asserts in examples are naive and it is not qualitative test cases. Do not inclue this "tests" to codecov testset"""

# false positive for f
# pylint: disable=undefined-variable

import glob
from os.path import basename, dirname, isfile, join

import pytest


@pytest.mark.parametrize('f', (pytest.param(__import__(basename(f)[:-3]).main, id=f) for f in glob.glob(join(dirname(__file__), '*.py')) if isfile(f) and not f.endswith('/test_all.py')))
def test(f):
    f()
