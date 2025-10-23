def initProducer():
    """Kafka helpers."""


import datetime as dt
import json
from typing import List

from kafka import KafkaProducer, KafkaConsumer


def init_producer(servers):
    # init an instance of KafkaProducer
    print('{} - Initializing Kafka producer at {}'.format(dt.datetime.utcnow()))
    return KafkaProducer(
      bootstrap_servers=servers,
      value_serializer=lambda v: json.dumps(v, default=str).encode('utf-8')
    )


def init_consumer(servers: List[str], topic:str, timeout: int = 1000):
    return KafkaConsumer(
        topic,
        bootstrap_servers=servers,
        group_id=None,
        auto_offset_reset='earliest',
        enable_auto_commit=False,
        consumer_timeout_ms=timeout,
        value_deserializer=lambda m: json.loads(m.decode('utf-8'))
    )