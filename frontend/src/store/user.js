import { defineStore } from 'pinia'
import axios from 'axios'

export const useUserStore = defineStore('user', {
    state: () => (
        {
            isLoading: false,
            tg: undefined,
            tgInited: false,
            subscriptions: [
                { name: "test 1", connections: [{name: "connection 1"}] },
                { name: "test 2", connections: [{name: "connection 2"}] },
            ],
            clinets: [],
            client: undefined,
            links: [],
        }
    ),
    actions: {
        async init() {
            this.tg = window.Telegram?.WebApp
            this.tgInited = this.tg.initData != ""

            if (this.tgInited) {
                console.log("Telegram mini app inited")
                this.tg.ready()
                this.tg.expand()
            }
        },
        async loadUserData() {
            try {
                const response = await axios.get("/client", {
                    headers: {
                        "Authorization": `Bearer ${this.tg.initData}`,
                    }
                })
                const client = response.data[0]?.client
                this.client = client
            } catch(err) {
                console.error("failed to get user data: ", err)
            }
        },
        async loadSubscriptions() {
            try {
                const response = await axios.get(`/subscriptions/${this.client.subId}`, {
                    headers: {"Authorization": `Bearer ${this.tg.initData}`},
                })

                console.log("loaded links", response.data)
                this.links = response.data

            } catch(err) {
                console.error(err)
            }
        }
    },
    getters: {
        username: (state) => state.tg?.initDataUnsafe.user?.username ?? "",
        id: (state) => state.tg?.initDataUnsafe.user?.id ?? -1,
        first_name: (state) => state.tg?.initDataUnsafe.user?.first_name,
        last_name: (state) => state.tg?.initDataUnsafe.user?.last_name,
        photo_url: (state) => state.tg?.initDataUnsafe.user?.photo_url,
        get_links: (state) => {
            const links = state.links.map((link) => {
                const url = URL.parse(link)
                const name = decodeURIComponent(url.hash)
                const protocol = url.protocol

                return `${protocol} ${name}`
            })
            return links
        }
    }
})