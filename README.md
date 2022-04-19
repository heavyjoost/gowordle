# gowordle
Wordle plugin that I slapped together for [go-chat-bot](https://github.com/go-chat-bot/bot), so you can play wordle on IRC.

By default it will load the file "words" which is a list of words separated by newlines (e.g. a copy of /usr/share/dict/words). To change what file it loads, set the environment variable called GOWORDLE_WORDS_PATH. It'll filter out 5-letter words, though you can change this in the code.
