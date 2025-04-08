package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
)

func SetEvents(newEvents []Event) {
	events = newEvents
}

func GetEvents() []Event {
	return events
}

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

func SaveEvents() {
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

			if now.Format("15:04") == "10:15" && notificationChannelID != "" {
				s.ChannelMessageSend(notificationChannelID, "Hell yeah, it's 10:15 already")
			}

			loadEvents()
			updated := false
			for i, event := range events {
				eventTime := event.DateTime.In(loc)

				// Evento vencido e ainda não notificado
				if !event.Notified && now.After(eventTime) && now.Format("2006-01-02 15:04") >= eventTime.Format("2006-01-02 15:04") {
					age := now.Year() - event.DateTime.Year()
					msg := fmt.Sprintf("🔔 Evento '%s' começou!\n📅 %s\n📝 %s", event.Name, eventTime.Format("02/01/2006 15:04"), event.Description)

					if strings.HasPrefix(strings.ToLower(event.Name), "aniversário de") {
						msg += fmt.Sprintf("\n🎉 %s está fazendo %d anos! 🎂", event.Name[14:], age)

						// Reagenda o próximo aniversário no mesmo horário do próximo ano
						nextYear := now.Year() + 1
						nextAnniversary := time.Date(nextYear, eventTime.Month(), eventTime.Day(), eventTime.Hour(), eventTime.Minute(), 0, 0, eventTime.Location())
						events[i].DateTime = nextAnniversary
						events[i].Notified = false
					} else {
						events[i].Notified = true
					}

					s.ChannelMessageSend(notificationChannelID, msg)
					updated = true
				}
			}
			if updated {
				SaveEvents()
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
	case strings.HasPrefix(content, "!meneximage"):
		sendRandomMenexImage(s, m.ChannelID)
	case strings.HasPrefix(content, "!removeevent"):
		args := strings.TrimPrefix(content, "!removeevent")
		nameToRemove := strings.TrimSpace(args)
		if nameToRemove == "" {
			s.ChannelMessageSend(m.ChannelID, "Uso correto: !removeevent [nome]")
			return
		}
		loadEvents()
		initialLen := len(events)
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
		SaveEvents()
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Evento '%s' removido com sucesso!", nameToRemove))

	case strings.HasPrefix(content, "!birthday"):
		args := strings.TrimSpace(strings.TrimPrefix(content, "!birthday"))
		parts := strings.Fields(args)

		if len(parts) < 3 {
			s.ChannelMessageSend(m.ChannelID, "Uso correto: !birthday nome DD/MM/AAAA descrição")
			return
		}

		// Find the date field
		dateIndex := -1
		for i, p := range parts {
			if len(p) == 10 && strings.Count(p, "/") == 2 {
				dateIndex = i
				break
			}
		}

		if dateIndex == -1 {
			s.ChannelMessageSend(m.ChannelID, "Formato de data inválido. Use DD/MM/AAAA")
			return
		}

		name := strings.Join(parts[:dateIndex], " ")
		dateStr := parts[dateIndex]
		desc := strings.Join(parts[dateIndex+1:], " ")

		birthDate, err := time.Parse("02/01/2006", dateStr)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Formato de data inválido. Use DD/MM/AAAA")
			return
		}

		if len(name) > 25 || len(desc) > 100 {
			s.ChannelMessageSend(m.ChannelID, "Nome deve ter até 25 caracteres e descrição até 100.")
			return
		}

		currentYear := time.Now().Year()
		anniversary := time.Date(currentYear, birthDate.Month(), birthDate.Day(), 10, 0, 0, 0, time.Local)
		event := Event{
			ID:          uuid.NewString(),
			Name:        fmt.Sprintf("Aniversário de %s", name),
			DateTime:    anniversary,
			Description: desc,
			Notified:    false,
		}
		events = append(events, event)
		SaveEvents()
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Aniversário de %s registrado para %s!", name, anniversary.Format("02/01")))

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
		args := strings.TrimPrefix(content, "!addevent")
		fields := strings.FieldsFunc(args, func(r rune) bool { return r == '"' })
		if len(fields) < 4 {
			s.ChannelMessageSend(m.ChannelID, "Uso correto: !addevent \"nome\" DD/MM/AAAA HH:MM \"descrição\"")
			return
		}
		name := strings.TrimSpace(fields[0])
		dateStr := strings.TrimSpace(fields[1])
		timeStr := strings.TrimSpace(fields[2])
		desc := strings.TrimSpace(fields[3])
		datetime, err := time.Parse("02/01/2006 15:04", dateStr+" "+timeStr)
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
		SaveEvents()
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Evento '%s' adicionado com sucesso!", name))

	case strings.HasPrefix(content, "!motorola"):
		phone := motorolaPhones[time.Now().UnixNano()%int64(len(motorolaPhones))]
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Estou fazendo testes no %s", phone))

	case strings.HasPrefix(content, "!menexhelp") || strings.HasPrefix(content, "!menex help"):
		help := `Comandos disponíveis:
!setchannel - Define o canal atual como o canal de notificações
!events - Lista todos os eventos pendentes
!addevent "nome" DD/MM/AAAA HH:MM "descrição" - Adiciona um novo evento
!motorola - Mostra um celular aleatório da Motorola
!menexhelp ou !menex help - Mostra todos os comandos
!removeevent "nome" - Remove um evento pelo nome
!birthday nome DD/MM/AAAA descrição - Registra seu aniversário para notificação anual
!meneximage - Envia uma imagem aleatória do Menex`
		s.ChannelMessageSend(m.ChannelID, help)
	}
}

func sendRandomMenexImage(s *discordgo.Session, channelID string) {
	var rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	imageDir := "images/menex"

	files, err := os.ReadDir(imageDir)
	if err != nil {
		log.Println("Error reading image directory:", err)
		return
	}

	if len(files) == 0 {
		s.ChannelMessageSend(channelID, "No Menex images found 😢")
		return
	}

	randomFile := files[rng.Intn(len(files))].Name()
	fullPath := filepath.Join(imageDir, randomFile)

	file, err := os.Open(fullPath)
	if err != nil {
		log.Println("Could not open random image:", err)
		s.ChannelMessageSend(channelID, "Failed to load Menex 😢")
		return
	}
	defer file.Close()

	_, err = s.ChannelFileSend(channelID, randomFile, file)
	if err != nil {
		log.Println("Error sending image:", err)
	}
}
