# Menex-Bot
Simple golang bot for discord, main goal is just to have some fun

# Project use
This bot will basically check if event dates that I have registered have been reached, simple.
!setchannel - Define o canal atual como o canal de notificações
!events - Lista todos os eventos pendentes
!addevent [nome] [DD/MM/AAAA] [HH:MM] [descrição] - Adiciona um novo evento
!menex help or !menexhelp - shows all commands
!motorola - tells you a random motorola phone
!removeevent [nome]
!birthday [nome] [DD/MM] [descrição] - cria uma data aniversário
!meneximage - Envia uma imagem aleatória do Menex


# Repository structure
        |
        V
## cmd
The cmd repo will register the bot commands (the fent core)

## main
Basic golang main structure, used for calling upon the function and commands.