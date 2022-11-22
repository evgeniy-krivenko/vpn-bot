package usecases

import _ "embed"

const (
	MainMenuText         = "Создать или проверить подключения можно через основное меню    \n👇"
	NoConnectionsText    = "У вас нет подключений. Создайте новое.  \nСписок доступных серверов:"
	FirstConnectionText  = "Создайте подключение.  \nСписок доступных серверов:"
	ConnectionsText      = "Ваши подключения:\n"
	ConnectCreated       = "Вы создали подключение, сверху 👆 его конфигурация, скопируйте и вставьте ее в приложение Outline. \nДля использования вам необходимо активировать его. Активация действует 12 часов. Потом придется повторить активацию."
	ActivateBtn          = "Активировать"
	ConnectionInfo       = "Страна: %s\nСтатус: %s\n"
	ConnectionTimeLeft   = "Отключение через: %s\n"
	ConnectionExists     = "Подключение для этого сервера уже существует, активируйте его или выберите другой"
	ShowConnectionConfig = "Кликните, чтобы скопировать подключение и вставьте его в Outline, оно удалится через 30 секунд  \n👇"
	AdvertisingMock      = "Тут может быть ваша реклама.  \n\nВы можете поддержать нас через криптокошельки.  \nДождитесь активации и нажмите \"Подтвердить\""
	ActivateDoneBtn      = "Подтвердить ✅"
)

//go:embed md/start-text.md
var StartText string

var ConnectionStatusEmoji = map[bool]string{true: "🟢", false: "🟡"}
var ConnectionStatusText = map[bool]string{true: "Активно", false: "Отключено"}

var LocationFullName = map[string]string{
	"NL": "The Netherlands 🇳🇱",
	"FI": "Finland 🇫🇮",
}
