import requests
import re


def test_tweet_creation_without_tags(base_url, tweet, correct_token):
    """
    Ensure that we can't create tweet without tag
    """
    text, _ = tweet(tags_count=0)
    resp = requests.post(
        base_url + '/tweets',
        json={'text': text},
        headers={'X-Auth-Token': correct_token}
    )
    assert resp.status_code == 400
    assert resp.json()[0]['code'] == 'tagless_tweets_not_implemented'


def test_tweet_creation_with_tags(base_url, tweet, correct_token):
    """
    Ensure that we can create tweet with tags and check that all published
    tweets are accessible through /tags/{tag}/tweets
    """
    for tag_count in range(1, 5):
        text, tags = tweet(tags_count=tag_count)
        resp = requests.post(
            base_url + '/tweets',
            json={'text': text},
            headers={'X-Auth-Token': correct_token}
        )
        location = resp.headers['location']

        assert resp.status_code == 201
        assert re.match('/tweets/[0-9]+', location) is not None, \
            'invalid Location header: ' + location

        tweet_id = location.split('/')[-1]

        for tag in tags:
            resp = requests.get(
                base_url + '/tags/{}/tweets'.format(tag)
            )

            assert resp.status_code == 200
            tweets = resp.json()['tweets']
            assert any([tweet['id'] == tweet_id for tweet in tweets]), \
                'tweet_id {} not found in: {}'.format(tweet_id, tweets)
