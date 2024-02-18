from telethon import TelegramClient

# Remember to use your own values from my.telegram.org!
api_id = 27886399
api_hash = 'e9c6d0cf04d88aae4bf528fcdc4d8270'
channel_id = -1002111285950

client = TelegramClient('anon', api_id, api_hash)


async def main():

    # Getting information about yourself
    me = await client.get_me()

    # "me" is a user object. You can pretty-print
    # any Telegram object with the "stringify" method:
    print(me.stringify())

    # When you print something, you see a representation of it.
    # You can access all attributes of Telegram objects with
    # the dot operator. For example, to get the username:
    # username = me.username
    # print(username)
    # print(me.phone)


    # # # You can, of course, use markdown in your messages:
    # message = await client.send_message(
    #     -1002111285950,
    #     'This message has **bold**, `code`, __italics__ and '
    #     'a [nice website](https://example.com)!',
    #     link_preview=False
    # )

    # msg = await client.get_messages(channel_id, ids=14)
    # print("msg", msg)
    # print("msg.views", msg.views)
    #
    # You can print the message history of any chat:
    async for message in client.iter_messages(channel_id):
        print(message)


with client:
    client.loop.run_until_complete(main())