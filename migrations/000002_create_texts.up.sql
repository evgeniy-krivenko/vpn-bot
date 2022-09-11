CREATE TABLE IF NOT EXISTS texts
(
    id         serial    NOT NULL,
    text       text,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp
);

INSERT INTO texts (text)
VALUES ('Добро пожаловать в Particius VPN. Мы предлагаем бесплатный VPN в обмен на просмотр рекламы. Не пишем логи и не продаем ваши данные'),
       ('Мы предоставляем вам доступ на сутки, далее вам нужно будет зайти в бот и обновить доступ'),
       ('Для использования вам нужно будет скачать клиент Outline для вашего устройства: https://getoutline.org/ru/get-started/ и после этого запросить конфигурацию');