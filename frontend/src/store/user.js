import { defineStore } from 'pinia'

export const useUserStore = defineStore('user', {
    state: () => (
        {
            tg: undefined,
            tgInited: false,
        }
    ),
    actions: {
        init() {
            this.tg = window.Telegram?.WebApp
            this.tgInited = this.tg.initData != ""

            if (this.tgInited) {
                console.log("Telegram mini app inited")
                this.tg.ready()
                this.tg.expand()
            }
        }
    },
    getters: {
        username: (state) => state.tg?.initDataUnsafe.user?.username ?? "",
        id: (state) => state.tg?.initDataUnsafe.user?.id ?? -1,
        first_name: (state) => state.tg?.initDataUnsafe.user?.first_name,
        last_name: (state) => state.tg?.initDataUnsafe.user?.last_name,
        photo_url: (state) => state.tg?.initDataUnsafe.user?.photo_url,
    }
})