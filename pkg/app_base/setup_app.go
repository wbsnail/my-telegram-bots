package app_base

func SetupBaseApp(options BaseOptions) BaseApp {
	return BaseApp{
		Version:             options.Version,
		Addr:                options.Addr,
		Name:                options.Name,
		TelegramAdminChatID: options.TelegramAdminChatID,
	}
}
