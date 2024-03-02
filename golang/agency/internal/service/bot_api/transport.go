package bot_api

type Transport interface {
	MonitorChannels()
	RunNotificationService()
}
