import re
from os.path import join

from setuptools import setup

MOD_NAME = 'milisp'
SRC_DIR = 'src'


def find_version():
    version_file = open(join(SRC_DIR, MOD_NAME, '__init__.py'), 'r').read()
    version_match = re.search(r"""^__version__ = ['"]([^'"]+)['"]""", version_file, re.M)
    if version_match:
        return version_match.group(1)
    raise RuntimeError('Unable to find version string.')


setup(
    name=MOD_NAME,
    version=find_version(),
    description='Micro Lisp',
    long_description='Python implementation of simple lisp-like highly extendable language engine',
    author='Alexey Michurin',
    author_email='a.michurin@gmail.com',
    packages=[MOD_NAME],
    package_dir={'': SRC_DIR},
    url='https://github.com/michurin/milisp',
    platforms=['any'],
    license='MIT License',
)
