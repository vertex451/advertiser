import asyncio
from flask import Flask, request, jsonify
from telethon import TelegramClient

api_id = 27886399
api_hash = 'e9c6d0cf04d88aae4bf528fcdc4d8270'

# channel_id = -1002093237940

app = Flask(__name__)


async def get_channel_subscribers(channel_handle):
    async with TelegramClient('session_name2', api_id, api_hash) as client:
        try:
            # Get information about the channel
            channel = await client.get_entity(channel_handle)
            # Get the participant list of the channel
            participants = await client.get_participants(channel, aggressive=True)
            # Get the total count of participants
            participants_count = participants.total
            return participants_count
        except Exception as e:
            return f"An error occurred: {e}"


@app.route('/stats', methods=['POST'])
def stats():
    data = request.json
    if data and 'channel_handle' in data:
        channel_handle = data['channel_handle']
        # Run the asynchronous function to get channel subscribers count
        subscribers_count = asyncio.run(get_channel_subscribers(channel_handle))
        return jsonify({'subscribers_count': subscribers_count})
    else:
        return jsonify({'error': 'Channel handle not provided'})


if __name__ == '__main__':
    app.run(debug=True)
