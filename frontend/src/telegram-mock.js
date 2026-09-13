export function initTelegramMock() {
    if (!window.Telegram) window.Telegram = {}
    const mockWebApp = {
        initData: "abcdefg",
        initDataUnsafe: {
            user: {
                id: 1,
                first_name: "Ivan",
                last_name: "Ivanov",
                username: "ivan_ivanov",
                photo_url: "https://ui-avatars.com/api/?name=Ivan+Ivanov&size=320"
            },
        },
        ready() {},
        expand() {},
    }

    Object.defineProperty(window.Telegram, 'WebApp', {
        value: mockWebApp,
        writable: false,
        configurable: true
    })

    console.warn("Telegram WebApp Mock inited")
}