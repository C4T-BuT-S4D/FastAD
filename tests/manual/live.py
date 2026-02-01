import asyncio
import logging
import signal

from centrifuge import (
    Client,
    ClientEventHandler,
    ConnectedContext,
    ConnectingContext,
    DisconnectedContext,
    ErrorContext,
    JoinContext,
    LeaveContext,
    PublicationContext,
    SubscribedContext,
    SubscribingContext,
    SubscriptionErrorContext,
    UnsubscribedContext,
    SubscriptionEventHandler,
    ServerSubscribedContext,
    ServerSubscribingContext,
    ServerUnsubscribedContext,
    ServerPublicationContext,
    ServerJoinContext,
    ServerLeaveContext,
)

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s - %(name)s - %(levelname)s - %(message)s",
)
# Configure centrifuge-python logger.
cf_logger = logging.getLogger("centrifuge")
cf_logger.setLevel(logging.DEBUG)


class ClientEventLoggerHandler(ClientEventHandler):
    """Check out comments of ClientEventHandler methods to see when they are called."""

    async def on_connecting(self, ctx: ConnectingContext) -> None:
        logging.info("connecting: %s", ctx)

    async def on_connected(self, ctx: ConnectedContext) -> None:
        logging.info("connected: %s", ctx)

    async def on_disconnected(self, ctx: DisconnectedContext) -> None:
        logging.info("disconnected: %s", ctx)

    async def on_error(self, ctx: ErrorContext) -> None:
        logging.error("client error: %s", ctx)

    async def on_subscribed(self, ctx: ServerSubscribedContext) -> None:
        logging.info("subscribed server-side sub: %s", ctx)

    async def on_subscribing(self, ctx: ServerSubscribingContext) -> None:
        logging.info("subscribing server-side sub: %s", ctx)

    async def on_unsubscribed(self, ctx: ServerUnsubscribedContext) -> None:
        logging.info("unsubscribed from server-side sub: %s", ctx)

    async def on_publication(self, ctx: ServerPublicationContext) -> None:
        logging.info("publication from server-side sub: %s", ctx.pub.data)

    async def on_join(self, ctx: ServerJoinContext) -> None:
        logging.info("join in server-side sub: %s", ctx)

    async def on_leave(self, ctx: ServerLeaveContext) -> None:
        logging.info("leave in server-side sub: %s", ctx)


class SubscriptionEventLoggerHandler(SubscriptionEventHandler):
    """Check out comments of SubscriptionEventHandler methods to see when they are called."""

    def __init__(self, channel: str = None):
        self.channel = channel

    async def on_subscribing(self, ctx: SubscribingContext) -> None:
        logging.info("[%s] subscribing: %s", self.channel, ctx)

    async def on_subscribed(self, ctx: SubscribedContext) -> None:
        logging.info("[%s] subscribed: %s", self.channel, ctx)

    async def on_unsubscribed(self, ctx: UnsubscribedContext) -> None:
        logging.info("[%s] unsubscribed: %s", self.channel, ctx)

    async def on_publication(self, ctx: PublicationContext) -> None:
        logging.info("[%s] publication: %s", self.channel, ctx.pub.data)

    async def on_join(self, ctx: JoinContext) -> None:
        logging.info("join: %s", ctx)

    async def on_leave(self, ctx: LeaveContext) -> None:
        logging.info("leave: %s", ctx)

    async def on_error(self, ctx: SubscriptionErrorContext) -> None:
        logging.error("[%s] subscription error: %s", self.channel, ctx)


def run_example():
    client = Client(
        "ws://localhost:8001/centrifuge/websocket",
        events=ClientEventLoggerHandler(),
        use_protobuf=False,
    )

    attacks_sub = client.new_subscription(
        "attacks",
        events=SubscriptionEventLoggerHandler("attacks"),
    )

    scoreboard_sub = client.new_subscription(
        "scoreboard",
        events=SubscriptionEventLoggerHandler("scoreboard"),
    )

    async def run():
        await client.connect()
        await attacks_sub.subscribe()
        await scoreboard_sub.subscribe()

        logging.info("all done, client connection is still alive, press Ctrl+C to exit")

    asyncio.ensure_future(run())
    loop = asyncio.get_event_loop()

    async def shutdown(received_signal):
        logging.info("received exit signal %s...", received_signal.name)
        await client.disconnect()

        tasks = [t for t in asyncio.all_tasks() if t is not asyncio.current_task()]
        for task in tasks:
            task.cancel()

        logging.info("Cancelling outstanding tasks")
        await asyncio.gather(*tasks, return_exceptions=True)
        loop.stop()

    signals = (signal.SIGTERM, signal.SIGINT)
    for s in signals:
        loop.add_signal_handler(
            s, lambda received_signal=s: asyncio.create_task(shutdown(received_signal))
        )

    try:
        loop.run_forever()
    finally:
        loop.close()
        logging.info("successfully completed service shutdown")


if __name__ == "__main__":
    run_example()
