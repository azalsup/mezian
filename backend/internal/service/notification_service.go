// Package service contient la logique métier de l'application.
package service

import (
	"fmt"
	"log"
	"mezian/internal/config"
)

// NotificationChannel représente le canal d'envoi d'un message.
type NotificationChannel string

const (
	ChannelSMS      NotificationChannel = "sms"
	ChannelWhatsApp NotificationChannel = "whatsapp"
)

// NotificationService définit l'interface d'envoi de notifications.
type NotificationService interface {
	// SendOTP envoie un code OTP à un numéro de phone via le canal spécifié.
	SendOTP(phone, code string, channel NotificationChannel) error
	// SendSMS envoie un SMS brut.
	SendSMS(phone, message string) error
}

// --- Mock (development) ---

// MockNotificationService journalise les messages dans la console.
type MockNotificationService struct{}

// NewMockNotificationService creates un service de notification mock.
func NewMockNotificationService() *MockNotificationService {
	return &MockNotificationService{}
}

func (m *MockNotificationService) SendOTP(phone, code string, channel NotificationChannel) error {
	log.Printf("[MOCK NOTIF] OTP %s → %s (canal: %s)", code, phone, channel)
	return nil
}

func (m *MockNotificationService) SendSMS(phone, message string) error {
	log.Printf("[MOCK NOTIF] SMS → %s: %s", phone, message)
	return nil
}

// --- Twilio SMS ---

// TwilioSMSService envoie des SMS via l'API Twilio.
type TwilioSMSService struct {
	accountSID string
	authToken  string
	fromNumber string
}

// NewTwilioSMSService creates un service Twilio.
func NewTwilioSMSService(cfg config.TwilioConfig) *TwilioSMSService {
	return &TwilioSMSService{
		accountSID: cfg.AccountSID,
		authToken:  cfg.AuthToken,
		fromNumber: cfg.FromNumber,
	}
}

func (t *TwilioSMSService) SendOTP(phone, code string, channel NotificationChannel) error {
	body := fmt.Sprintf("Votre code Mezian est: %s (valable 10 min). Ne le partagez pas.", code)
	return t.SendSMS(phone, body)
}

func (t *TwilioSMSService) SendSMS(phone, message string) error {
	// TODO: implémenter l'appel HTTP à l'API Twilio
	// POST https://api.twilio.com/2010-04-01/Accounts/{AccountSid}/Messages.json
	// avec Basic Auth (AccountSID:AuthToken) et les champs To, From, Body
	log.Printf("[TWILIO STUB] SMS → %s: %s", phone, message)
	return nil
}

// --- WhatsApp ---

// WhatsAppService envoie des messages via l'API WhatsApp Business (Meta).
type WhatsAppService struct {
	apiURL        string
	apiToken      string
	phoneNumberID string
}

// NewWhatsAppService creates un service WhatsApp.
func NewWhatsAppService(cfg config.WhatsAppConfig) *WhatsAppService {
	return &WhatsAppService{
		apiURL:        cfg.APIURL,
		apiToken:      cfg.APIToken,
		phoneNumberID: cfg.PhoneNumberID,
	}
}

func (w *WhatsAppService) SendOTP(phone, code string, channel NotificationChannel) error {
	body := fmt.Sprintf("Votre code Mezian est: *%s* (valable 10 min). Ne le partagez jamais.", code)
	return w.SendSMS(phone, body)
}

func (w *WhatsAppService) SendSMS(phone, message string) error {
	// TODO: implémenter l'appel HTTP à l'API Meta WhatsApp Business
	// POST https://graph.facebook.com/v19.0/{phone_number_id}/messages
	log.Printf("[WHATSAPP STUB] Message → %s: %s", phone, message)
	return nil
}

// NewNotificationService creates le service de notification selon la configuration.
func NewNotificationService(cfg *config.Config) NotificationService {
	switch cfg.Notification.Provider {
	case "twilio":
		return NewTwilioSMSService(cfg.Notification.Twilio)
	case "whatsapp":
		return NewWhatsAppService(cfg.Notification.WhatsApp)
	default:
		return NewMockNotificationService()
	}
}
