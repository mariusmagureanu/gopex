#!/usr/bin/env python3
from pynats import NATSClient

with NATSClient() as client:
    client.connect()
    print("start listening on nats...\n")

    def callback(msg):
        print(f"received a message with subject {msg.subject}: {msg}")

    client.subscribe(subject="sse", callback=callback)
    client.wait(count=100)