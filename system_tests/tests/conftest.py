import pytest


@pytest.fixture
def base_url():
    import os
    return os.environ['URL']


@pytest.fixture
def correct_token():
    return 'correct-token-1'


@pytest.fixture
def words():
    return [
        'Lorem', 'ipsum', 'dolor', 'sit', 'amet', 'consectetur',
        'adipiscing', 'elit', 'Fusce', 'eget', 'facilisis', 'massa',
        'a', 'fermentum', 'est'
    ]


@pytest.fixture
def tweet(words):
    import random
    def wrapped(length=None, tags_count=1):
        if length is None:
            length = random.randint(0, 20)
        tags = random.choices(words, k=tags_count)
        text = (
            ['#' + t for t in tags]
            + random.choices(words, k=length)
        )
        random.shuffle(text)

        return (' '.join(text), tags)
    return wrapped
