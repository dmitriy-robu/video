db = db.getSiblingDB('api');
db.createCollection("run_time")

db.run_time.insert({
    started_at: new Date()
})

db.pages.insertMany([
    {
        "delete": 0,
        "link": "/keyboard",
        "template": "<p>Добро пожаловать, <strong class=\"font-bold\">${USER}!</strong></p>",
        "title": "Начальная клавиатура",
        "uuid": "ab0f0e7c-43b0-481c-84a2-812695a822ac"
    },
    {
        "delete": 0,
        "link": "/start",
        "template": "<p><strong class=\"font-bold\">Добро пожаловать в бот автопродаж Риддика </strong>🤝</p><p></p><p>• Наш магазин предлагает только органические (марихуана, гашиш) товары, которые как было доказано медициной, не подвергают риску ваше здоровье или жизнь.</p><p></p><p>• Наш оператор проконсультирует и ответит на все ваши вопросы <a target=\"_blank\" rel=\"noopener noreferrer nofollow\" class=\"inline-flex items-center gap-x-1 text-blue-600 decoration-2 hover:underline font-medium dark:text-white\" href=\"https://t.me/RDDCK420\">@RDDCK420</a></p><p></p><p>• Список обменников и краткая инструкция по оплате доступна если нажать на команду <a target=\"_blank\" rel=\"noopener noreferrer nofollow\" class=\"inline-flex items-center gap-x-1 text-blue-600 decoration-2 hover:underline font-medium dark:text-white\" href=\"/obmen\">/obmen</a></p><p></p><p>• Важно: перед поднятием клада, обязательно ознакомьтесь с политикой магазина по разрешению ненаходов или несоответствие веса нажимая тут <a target=\"_blank\" rel=\"noopener noreferrer nofollow\" class=\"inline-flex items-center gap-x-1 text-blue-600 decoration-2 hover:underline font-medium dark:text-white\" href=\"/rules\">/rules</a></p><p></p><p>Далее выбирайте нужный вам город:</p>",
        "title": "Приветсвенное сообщение",
        "uuid": "7ef0b9c7-2056-431d-8eb5-6cbe3f0569bb"
    },
    {
        "delete": 0,
        "link": "/city",
        "template": "<p>Вы выбрали <strong class=\"font-bold\">${CITY}</strong>!<br>Выберите интересующий вас товар:</p>",
        "title": "Город выбран",
        "uuid": "6f5721c4-ace9-4f5b-9fc7-30fd485d1d6c"
    },
    {
        "delete": 0,
        "link": "/product",
        "template": "<p>Вы выбрали товар: <strong class=\"font-bold\">${PRODUCT}</strong></p><p>Город: <strong class=\"font-bold\">${CITY}</strong></p><p>Цена: <strong class=\"font-bold\">${PRICE}$</strong></p><p></p><p>Далее, выберите нужный вам район:</p>",
        "title": "Товар выбран",
        "uuid": "09466452-6f3f-41a9-9b0a-4bc21260d919"
    },
    {
        "delete": 0,
        "link": "/product-description",
        "template": "<p>Название: <strong class=\"font-bold\">${PRODUCT}</strong><br><br><strong class=\"font-bold\">Описание</strong>: ${DESCRIPTION}</p>",
        "title": "Описание товара",
        "uuid": "700e4138-daa1-48f5-8f44-da501b28d85d"
    },
    {
        "delete": 0,
        "link": "/orders",
        "template": "<p>Ваши заказы:<br>Общее кол-во: <strong class=\"font-bold\">${TOTAL_COUNT}</strong></p>",
        "title": "Заказы пользователя",
        "uuid": "88a8fdd8-1425-4e3f-a83a-b5d2e377c666"
    },
    {
        "delete": 0,
        "link": "/orders-empty",
        "template": "<p>У вас нет заказов на данный момент!<br>Пожалуйста, используйте команду для <a target=\"_blank\" rel=\"noopener noreferrer nofollow\" class=\"inline-flex items-center gap-x-1 text-blue-600 decoration-2 hover:underline font-medium dark:text-white\" href=\"/start\"><strong class=\"font-bold\">/start</strong></a> создания заказа!</p>",
        "title": "У пользователя нет заказов",
        "uuid": "09565847-b6e1-47aa-a179-a352b0d410e6"
    },
    {
        "delete": 0,
        "link": "/orders-collected",
        "template": "<p>Все нужные данные были успешно собраны:<br>Город:<strong class=\"font-bold\"> ${CITY}</strong><br>Район: <strong class=\"font-bold\">${REGION}</strong><br>Товар: <strong class=\"font-bold\">${PRODUCT}</strong></p><p></p><p>Заказ может быть оплачен до: <strong class=\"font-bold\">${RESERVED_TO}</strong></p>",
        "title": "Заказ собран",
        "uuid": "03ad4208-18f0-4517-88bd-922764d30211"
    },
    {
        "delete": 0,
        "link": "/orders-view",
        "template": "<p>Ваш заказ от <strong class=\"font-bold\">${DATE}:</strong></p><p>Состояние: <strong class=\"font-bold\">${STATUS}</strong><br>Город:<strong class=\"font-bold\"> ${CITY}</strong><br>Район: <strong class=\"font-bold\">${REGION}</strong><br>Товар: <strong class=\"font-bold\">${PRODUCT}</strong><br>Цена: <strong class=\"font-bold\">${PRICE}</strong></p><p></p><p>Описание:</p><p>${ADDRESS_DESCRIPTION}</p>",
        "title": "Информация о заказе",
        "uuid": "b9ddf981-5102-4b11-9912-f47d836a5388"
    },
    {
        "delete": 0,
        "link": "/balance",
        "template": "<p>Ваш текущий баланс: <strong class=\"font-bold\">${BALANCE}$</strong></p><p><br>Вы можете внести средства используя ваш постоянный LTC адрес:<br><code class=\"text-red-600 dark:text-red-400 p-1 rounded-md font-mono\">${LTC_ADDRESS}</code></p>",
        "title": "Баланс пользователя",
        "uuid": "4575abf5-e9ac-48c4-9cbf-9e5d58f29b80"
    },
    {
        "delete": 0,
        "link": "/balance-income",
        "template": "<p>Мы нашли новую транзакцию в вашем кошельке:<br><strong class=\"font-bold\">Адрес</strong>: <code class=\"text-red-600 dark:text-red-400 p-1 rounded-md font-mono\">${LTC_ADDRESS}</code></p><p><strong class=\"font-bold\">TxId</strong>: <a target=\"_blank\" rel=\"noopener noreferrer nofollow\" class=\"inline-flex items-center gap-x-1 text-blue-600 decoration-2 hover:underline font-medium dark:text-white\" href=\"https://sochain.com/tx/LTCTEST/${TXID_LINK}\"><strong class=\"font-bold\">${TXID}</strong></a></p><p><strong class=\"font-bold\">Текущий курс USD</strong>: <strong class=\"font-bold\">${USD_PRICE}</strong></p><p></p><p><strong class=\"font-bold\">Сумма в LTC</strong>: <strong class=\"font-bold\">${LTC_AMOUNT} LTC</strong></p><p><strong class=\"font-bold\">Сумма в USD</strong>:<strong class=\"font-bold\"> $${USD_AMOUNT}</strong></p><p></p><p>Сейчас мы ожидаем ${MIN_CONFIRMATIONS} подтверждений транзакции!</p>",
        "title": "Новая транзакция",
        "uuid": "e982d0e2-8e6c-44e8-a87c-fd0b7483e1e9"
    },
    {
        "delete": 0,
        "link": "/balance-update",
        "template": "<p>Транзакция <strong class=\"font-bold\">${TXID} </strong>была подтверждена <strong class=\"font-bold\">${CONFIRMATIONS} </strong>раз(-а).<br>Ваш текущий баланс: <strong class=\"font-bold\">${BALANCE}</strong></p>",
        "title": "Пополнение баланса",
        "uuid": "7825a6d7-1560-4a61-b360-08ac3dfc1502"
    },
    {
        "delete": 0,
        "link": "/balance-insufficient-funds",
        "template": "<p>На вашем балансе не хватает средств!</p><p></p><p>Текущий баланс: <strong class=\"font-bold\">${BALANCE}</strong></p><p>Сумма заказа: <strong class=\"font-bold\">${AMOUNT}</strong></p>",
        "title": "На балансе не хватает средств",
        "uuid": "d621a8f5-42bf-4289-a556-7e38ff358305"
    },
    {
        "delete": 0,
        "link": "/address-sent",
        "template": "<p>Номер заказа: <strong class=\"font-bold\">${ORDER_ID}</strong></p><p>Товар: <strong class=\"font-bold\">${PRODUCT}</strong></p><p>Город: <strong class=\"font-bold\">${CITY}</strong></p><p>Рацон: <strong class=\"font-bold\">${REGION}</strong></p><p>Цена: <strong class=\"font-bold\">${PRICE}$</strong></p><p>➖➖➖➖➖➖➖➖➖➖</p><p>Description: <strong class=\"font-bold\">${PRODUCT_DESCRIPTION}</strong></p>",
        "title": "Отправка адреса",
        "uuid": "5519c954-a5e4-4b70-9eb0-5000a71dd543"
    },
    {
        "delete": 0,
        "link": "/orders-cancelled",
        "template": "<p>Заказ был отменён по причине истечения срока резервации!</p>",
        "title": "Заказ был отменён",
        "uuid": "4bba6f74-26f3-40d3-922b-5c250d8932fd"
    },
    {
        "delete": 0,
        "link": "/address-empty",
        "template": "<p>На данный момент адресов нет!</p>",
        "title": "Нет адресов",
        "uuid": "abacc8f8-5553-49cf-a99c-f4929d7110f6"
    },
    {
        "delete": 0,
        "link": "/banned",
        "template": "<p>К сожалению, вы были заблокированы в нашей системе по причине превышения возможных попыток отмены!<br><br>Дата снятия блокировки: <strong class=\"font-bold\">${BANNED_UNTIL}</strong></p>",
        "title": "Пользователь забанен",
        "uuid": "7e3e0016-dd7a-49b4-b76b-3483eee09d2e"
    },
])

db.settings.insertMany([
    {
        "uuid": "b800842d-d1bf-46c2-aff7-342f21645b8b",
        "key": "bot.token",
        "value": "6496625047:AAFe96joxCIkwPXCbb_jrEdWH2T4FULJfgY",
    },
    {
        "uuid": "8190b3ed-c533-4c7c-81e6-dbcdaf058b9c",
        "key": "bot.state",
        "value": "1",
    },
    {
        "uuid": "2703280c-93eb-4515-bf8a-6ed61bd404fd",
        "key": "bot.message",
        "value": "Бот временно недоступен!",
    },
    {
        "key": "payment.wallet",
        "uuid": "edd06427-6fac-4dfb-8f43-de2ae05fbbb4",
        "value": "tltc1qlx5plauq6deektl7f8kxmm3afym6p5qcj5nwh7"
    },
    {
        "key": "payment.min_confirmations",
        "uuid": "3f9ca7c5-a851-4a8c-a017-29ab4dcf583c",
        "value": "1"
    },
    {
        "key": "order.timeout",
        "uuid": "86e8c375-0817-43b1-9478-6e051a8bb0c7",
        "value": "60"
    },
    {
        "uuid": "bbc82f5e-e42f-4227-8337-e7dde21d79ce",
        "key": "spam.processing",
        "value": "1",
    },
    {
        "uuid": "1afe4cc8-d755-42f3-ae9b-4b91e33affee",
        "key": "spam.cancel",
        "value": "1",
    },
    {
        "uuid": "23f1c493-2021-4be5-a2c1-c69de33fc1bf",
        "key": "spam.ban",
        "value": "999",
    },
    {
        "uuid": "62fa6372-6fe7-4c50-affe-a760c3d0a648",
        "key": "security.proxy",
        "value": "1",
    },
    {
        "key": "bot.image",
        "uuid": "df79a610-415c-4f88-9cd8-54d8c2d30243",
        "value": "https://telegra.ph/file/ab0ac93952b1da2888d05.jpg"
    }
])