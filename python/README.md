Micro Lisp
==========

The python part of [milisp](https://github.com/michurin/milisp) project.

Installation
------------

```sh
pip install 'git+https://github.com/michurin/milisp#subdirectory=python'
```

Contributing
------------

```sh
pip install pylint
pip install pytest
pip install pytest-cov
```

```sh
PYTHONPATH=src pytest --cov=milisp tests
PYTHONPATH=src pylint --variable-rgx='.*' --argument-rgx='.*' milisp
```

Todo
----

- Decorators for debugging
- Decorators for caching
- Docstrings
