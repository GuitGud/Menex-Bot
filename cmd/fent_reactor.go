package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
)

type Event struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	DateTime    time.Time `json:"dateTime"`
	Description string    `json:"description"`
	Notified    bool      `json:"notified"`
}

var (
	eventsFilePath          = "events.json"
	events                  []Event
	notificationChannelID   string
	notificationTickerStart = false
	motorolaPhones          = []string{
		"Moto G100",
		"Moto G200",
		"Moto G Power",
		"Moto E40",
		"Moto Edge 30",
		"Moto Razr 5G",
		"Moto G82",
		"Moto G73",
		"Moto Edge 40 Pro",
		"Moto G13",
		"...ERROR...FENT REACTOR DEFECTIVE",
		"Moto G Stylus",
		"Moto G Play",
		"Moto G Fast",
		"Moto G Stylus 5G",
	}
)

func loadEvents() {
	file, err := os.ReadFile(eventsFilePath)
	if err != nil {
		log.Println("Não foi possível carregar events.json, iniciando vazio.")
		return
	}
	_ = json.Unmarshal(file, &events)
}

func saveEvents() {
	file, _ := json.MarshalIndent(events, "", "  ")
	_ = os.WriteFile(eventsFilePath, file, 0644)
}

func startTicker(s *discordgo.Session) {
	if notificationTickerStart {
		return
	}
	notificationTickerStart = true

	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for {
			<-ticker.C
			loc, _ := time.LoadLocation("America/Sao_Paulo")
			now := time.Now().In(loc)

			// Mensagem especial às 10:15
			if now.Format("15:04") == "10:15" && notificationChannelID != "" {
				s.ChannelMessageSend(notificationChannelID, "Hell yeah, it's 10:15 already")
			}

			// Verificação de eventos
			loadEvents()
			updated := false
			for i, event := range events {
				eventTime := event.DateTime.In(loc)
				if !event.Notified && now.After(eventTime) {
					msg := fmt.Sprintf("🔔 Evento '%s' começou!\n📅 %s\n📝 %s",
						event.Name,
						eventTime.Format("02/01/2006 15:04"),
						event.Description,
					)
					s.ChannelMessageSend(notificationChannelID, msg)
					events[i].Notified = true
					updated = true
				}
			}
			if updated {
				saveEvents()
			}
		}
	}()
}

func HandleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	content := strings.TrimSpace(m.Content)
	if !strings.HasPrefix(content, "!") {
		return
	}

	switch {
	case strings.HasPrefix(content, "!removeevent"):
		args := strings.SplitN(content, " ", 2)
		if len(args) < 2 {
			s.ChannelMessageSend(m.ChannelID, "Uso correto: !removeevent [nome]")
			return
		}
		nameToRemove := strings.TrimSpace(args[1])
		loadEvents()
		initialLen := len(events)

		// Filtra todos exceto o com o nome exato
		newEvents := make([]Event, 0)
		for _, e := range events {
			if e.Name != nameToRemove {
				newEvents = append(newEvents, e)
			}
		}

		if len(newEvents) == initialLen {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Nenhum evento com nome '%s' foi encontrado.", nameToRemove))
			return
		}

		events = newEvents
		saveEvents()
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Evento '%s' removido com sucesso!", nameToRemove))

	case strings.HasPrefix(content, "!setchannel"):
		notificationChannelID = m.ChannelID
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Canal de notificações definido: <#%s>", notificationChannelID))
		startTicker(s)

	case strings.HasPrefix(content, "!events"):
		loadEvents()
		if len(events) == 0 {
			s.ChannelMessageSend(m.ChannelID, "Nenhum evento encontrado.")
			return
		}
		var builder strings.Builder
		builder.WriteString("Lista de eventos:\n")
		for _, event := range events {
			builder.WriteString(fmt.Sprintf("- %s em %s: %s\n", event.Name, event.DateTime.Format("02/01/2006 15:04"), event.Description))
		}
		s.ChannelMessageSend(m.ChannelID, builder.String())

	case strings.HasPrefix(content, "!addevent"):
		args := strings.SplitN(content, " ", 5)
		if len(args) < 5 {
			s.ChannelMessageSend(m.ChannelID, "Uso correto: !addevent [nome] [DD/MM/AAAA] [HH:MM] [descrição]")
			return
		}

		name := args[1]
		dateStr := args[2]
		timeStr := args[3]
		desc := args[4]
		loc, _ := time.LoadLocation("America/Sao_Paulo")
		datetime, err := time.ParseInLocation("02/01/2006 15:04", dateStr+" "+timeStr, loc)

		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Formato de data/hora inválido. Use DD/MM/AAAA HH:MM")
			return
		}

		event := Event{
			ID:          uuid.NewString(),
			Name:        name,
			DateTime:    datetime,
			Description: desc,
			Notified:    false,
		}
		events = append(events, event)
		saveEvents()
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Evento '%s' adicionado com sucesso!", name))

	case strings.HasPrefix(content, "!motorola"):
		phone := motorolaPhones[time.Now().UnixNano()%int64(len(motorolaPhones))]
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Um celular Motorola aleatório: %s", phone))

	case strings.HasPrefix(content, "!menexhelp") || strings.HasPrefix(content, "!menex help"):
		help := `Comandos disponíveis:
!setchannel - Define o canal atual como o canal de notificações
!events - Lista todos os eventos pendentes
!addevent [nome] [DD/MM/AAAA] [HH:MM] [descrição] - Adiciona um novo evento
!motorola - Mostra um celular aleatório da Motorola
!menexhelp ou !menex help - Mostra todos os comandos
!removeevent [nome] - Remove um evento pelo nome`
		s.ChannelMessageSend(m.ChannelID, help)
	}
}
