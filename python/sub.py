from telethon.sync import TelegramClient

api_id = 27886399
api_hash = 'e9c6d0cf04d88aae4bf528fcdc4d8270'

channel_id = -1002093237940

with TelegramClient('session_name', api_id, api_hash) as client:
    try:
        # Get information about the channel
        channel = client.get_entity(channel_id)
        # Get the subscriber count
        subscribers_count = client.get_participants(channel, aggressive=True).total
        print(f"Subscribers count of the channel: {subscribers_count}")
    except Exception as e:
        print(f"An error occurred: {e}")
