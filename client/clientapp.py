# Copyright 2020 Adap GmbH. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
# ==============================================================================
"""Flower client app."""


import time
from logging import INFO

from flwr.common import GRPC_MAX_MESSAGE_LENGTH
from flwr.common.logger import log

from flwr.client.client import Client
from connection import secure_grpc_connection, insecure_grpc_connection
from flwr.client.grpc_client.message_handler import handle
from flwr.client.numpy_client import NumPyClient, NumPyClientWrapper
import grpc


def start_client(
    server_address: str,
    client: Client,
    keyPrefix: str,
    grpc_max_message_length: int = GRPC_MAX_MESSAGE_LENGTH,
) -> None:
    """Start a Flower Client which connects to a gRPC server.

    Arguments:
        server_address: str. The IPv6 address of the server. If the Flower
            server runs on the same machine on port 8080, then `server_address`
            would be `"[::]:8080"`.
        client: flwr.client.Client. An implementation of the abstract base
            class `flwr.client.Client`.
        grpc_max_message_length: int (default: 536_870_912, this equals 512MB).
            The maximum length of gRPC messages that can be exchanged with the
            Flower server. The default should be sufficient for most models.
            Users who train very large models might need to increase this
            value. Note that the Flower server needs to be started with the
            same value (see `flwr.server.start_server`), otherwise it will not
            know about the increased limit and block larger messages.

    Returns:
        None.
    """
    while True:
        sleep_duration: int = 0
        with open(keyPrefix+'.crt', 'rb') as f:
            trusted_certs = f.read()
        credentials = grpc.ssl_channel_credentials(root_certificates=trusted_certs)
        with secure_grpc_connection(
            credentials, server_address, max_message_length=grpc_max_message_length
        ) as conn:
            receive, send = conn
            log(INFO, "Opened (insecure) gRPC connection")

            while True:
                server_message = receive()
                client_message, sleep_duration, keep_going = handle(
                    client, server_message
                )
                send(client_message)
                if not keep_going:
                    break
        if sleep_duration == 0:
            log(INFO, "Disconnect and shut down")
            break
        # Sleep and reconnect afterwards
        log(
            INFO,
            "Disconnect, then re-establish connection after %s second(s)",
            sleep_duration,
        )
        time.sleep(sleep_duration)


def start_numpy_client(
    server_address: str,
    client: NumPyClient,
    keyPrefix: str,
    grpc_max_message_length: int = GRPC_MAX_MESSAGE_LENGTH,
) -> None:
    """Start a Flower NumPyClient which connects to a gRPC server.

    Arguments:
        server_address: str. The IPv6 address of the server. If the Flower
            server runs on the same machine on port 8080, then `server_address`
            would be `"[::]:8080"`.
        client: flwr.client.NumPyClient. An implementation of the abstract base
            class `flwr.client.NumPyClient`.
        grpc_max_message_length: int (default: 536_870_912, this equals 512MB).
            The maximum length of gRPC messages that can be exchanged with the
            Flower server. The default should be sufficient for most models.
            Users who train very large models might need to increase this
            value. Note that the Flower server needs to be started with the
            same value (see `flwr.server.start_server`), otherwise it will not
            know about the increased limit and block larger messages.

    Returns:
        None.
    """

    # Wrap the NumPyClient
    flower_client = NumPyClientWrapper(client)

    # Start
    start_client(
        server_address=server_address,
        client=flower_client,
        keyPrefix=keyPrefix,
        grpc_max_message_length=grpc_max_message_length,
    )