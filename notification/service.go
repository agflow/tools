package notification

// Service is an interface of notification.Service
type Service interface {
	Send(string, string) error
}
