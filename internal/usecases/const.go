package usecases

const (
	MainMenuText       = "Создать или проверить подключения можно через основное меню    \n👇"
	NoConnectionsText  = "У вас нет подключений\\. Создайте новое\\.\nСписок доступных серверов:"
	ConnectionsText    = "Ваши подключения:\n"
	ConnectCreated     = "Вы создали подключение, сверху 👆 его конфигурация, скопируйте и вставьте ее в приложение Outline\\. \nДля использования вам необходимо активировать его\\. Активация действует 12 часов\\. Потом придется повторить активацию\\."
	ActivateBtn        = "Активировать"
	ConnectionInfo     = "Location: %s\nДата последней активации: %s\n"
	ConnectionTimeLeft = "Часов до окончания активации: %s\n"
)

var ConnectionStatusEmoji = map[bool]string{true: "🟢", false: "⚫️"}

var LocationFullName = map[string]string{
	"NL": "The Netherlands 🇳🇱",
}
