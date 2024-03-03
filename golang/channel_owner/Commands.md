start - start by choosing your role
all_topics - List all available topics
my_channels - List my channels where bot is present


# TODO
1. Move state to db to scale pods
2. Add report bug button
3. Fail if topics is nof in the list during Ad creation. Do the same for channels.


# Potential issues:
1. We sticks not to userID but to chatID and in case of new bot(new chat id) we need migration.
2. If bot i s banned, we can create new bot, we need to add it to each channel again. Store channel_owners contact