import { defineStore } from 'pinia'
import axios from 'axios'

export const useUserStore = defineStore('user', {
    state: () => (
        {
            tg: undefined,
            tgInited: false,
            subscriptions: [
                { name: "test 1", connections: [{name: "connection 1"}] },
                { name: "test 2", connections: [{name: "connection 2"}] },
            ],
            clinets: []
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

                const clientResponse = await axios.get(clients, {
                    headers: {
                        Authorization: `Bearer ${this.tg.initData}`
                    }
                })

                this.clinets = clientResponse.data
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