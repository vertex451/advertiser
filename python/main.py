import asyncio
from flask import Flask, request, jsonify
from telethon import TelegramClient

api_id = 27886399
api_hash = 'e9c6d0cf04d88aae4bf528fcdc4d8270'

app = Flask(__name__)


@app.route('/subscribers', methods=['POST'])
def subscribers():
    data = request.json
    if data and 'channel_handle' in data:
        channel_handle = data['channel_handle']
        # Run the asynchronous function to get channel subscribers count
        subscribers_count = asyncio.run(get_channel_subscribers(channel_handle))
        return jsonify({'subscribers_count': subscribers_count})
    else:
        return jsonify({'error': 'Channel handle not provided'})


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


# @app.route('/subscribers', methods=['POST'])
# def subscribers():
#     data = request.json
#     if data and 'channel_handle' in data:
#         channel_handle = data['channel_handle']
#         # Run the asynchronous function to get channel subscribers count
#         subscribers_count = asyncio.run(get_channel_subscribers(channel_handle))
#         return jsonify({'subscribers_count': subscribers_count})
#     else:
#         return jsonify({'error': 'Channel handle not provided'})

@app.route('/post_view_count', methods=['POST'])
def post_view_count():
    data = request.json
    if data and 'channel_handle' in data and 'message_id' in data:
        res = asyncio.run(get_post_view_count(data['channel_handle'], data['message_id']))
        # Wait for the task to complete
        return jsonify({'view_count': res})
    else:
        return jsonify({'error': 'Data is not provided'})


# async def get_channel_subscribers(channel_handle):
#     async with TelegramClient('session_name2', api_id, api_hash) as client:
#         try:
#             # Get information about the channel
#             channel = await client.get_entity(channel_handle)
#             # Get the participant list of the channel
#             participants = await client.get_participants(channel, aggressive=True)
#             # Get the total count of participants
#             participants_count = participants.total
#             return participants_count
#         except Exception as e:
#             return f"An error occurred: {e}"

async def get_post_view_count(channel_handle, message_id):
    async with TelegramClient('session_name2', api_id, api_hash) as client:
        try:
            message = await client.get_messages(channel_handle, ids=int(message_id))
            if message.views:
                return message.views
            else:
                return 'No view count available for this message.'
        except Exception as e:
            return f"An error occurred: {e}"


if __name__ == '__main__':
    app.run(debug=True)
