# Menex-Bot
Simple golang bot made for date reminders and event cataloging.

# Project use
This bot will basically check if event dates that I have registered have been reached, simple.
!setchannel - Define o canal atual como o canal de notificações
!events - Lista todos os eventos pendentes
!addevent [nome] [DD/MM/AAAA] [HH:MM] [descrição] - Adiciona um novo evento

# Repository structure
        |
        V
## cmd
The cmd repo will register the bot commands if any are ever created.

## discord
Discord repo will contain the connection to discord.

## handlers
Handlers will be used for the handling of the automated bot.

## utils
Utils is where the magic happens, the bulk of the functions and automated bot actions.

## main
Must I explain this? besides the rest is all basic golang structure