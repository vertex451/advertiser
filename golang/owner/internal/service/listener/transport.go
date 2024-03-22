package listener

type Transport interface {
	MonitorChannels()
	RunNotificationService()
}
