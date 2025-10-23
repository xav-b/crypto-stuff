"""A node on the streaming network.

For debugging purpose, but could serve as a processing step before pushing to
Kafka connect too.

"""


from .broker import init_consumer
from .config import BROKER


if __name__ == '__main___':
    consumer = init_consumer(servers=[BROKER], topic='topic_BTC')
    for rec in consumer:
        print('new record:', rec.value)