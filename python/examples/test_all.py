"""Just shallow check that examples are runable. Asserts in examples are naive and it is not qualitative test cases. Do not inclue this "tests" to codecov testset"""

import glob
from os.path import basename, dirname, isfile, join

import pytest


@pytest.mark.parametrize('f', (__import__(basename(f)[:-3]).main for f in glob.glob(join(dirname(__file__), '*.py')) if isfile(f) and not f.endswith('/test_all.py')))
def test(f):
    f()
