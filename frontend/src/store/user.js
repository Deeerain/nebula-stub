import { defineStore } from 'pinia'
import axios from 'axios'

export const useUserStore = defineStore('user', {
    state: () => ({
        isLoading: false,
        tg: undefined,
        tgInited: false,
        token: '',
        client: null,
        links: [],
    }),
    actions: {
        async init() {
            this.tg = window.Telegram?.WebApp
            const hasInitData = Boolean(this.tg?.initData)

            if (!hasInitData) {
                return
            }

            this.tg.ready()
            this.tg.expand()

            await this.authenticate()
            this.tgInited = true
        },
        async authenticate() {
            const { data } = await axios.post('auth', { initData: this.tg.initData })
            this.token = data.token
            axios.defaults.headers.common.Authorization = `Bearer ${this.token}`
        },
        async loadUserData() {
            this.isLoading = true
            try {
                const { data } = await axios.get('client')
                this.client = data
            } catch (err) {
                if (err.response?.status === 404) {
                    this.client = null
                } else {
                    console.error('failed to get user data:', err)
                }
            } finally {
                this.isLoading = false
            }
        },
        async loadSubscriptions() {
            if (!this.client?.sub_id) {
                this.links = []
                return
            }

            try {
                const { data } = await axios.get(`subscriptions/${this.client.sub_id}`)
                this.links = Array.isArray(data) ? data : []
            } catch (err) {
                console.error('failed to get subscriptions:', err)
            }
        },
        async connect(app) {
            const { data } = await axios.get(`connect/${app}`)

            if (data?.deeplink) {
                window.location.href = data.deeplink
            }
        },
    },
    getters: {
        username: (state) => state.tg?.initDataUnsafe?.user?.username ?? '',
        id: (state) => state.tg?.initDataUnsafe?.user?.id ?? -1,
        first_name: (state) => state.tg?.initDataUnsafe?.user?.first_name ?? '',
        last_name: (state) => state.tg?.initDataUnsafe?.user?.last_name ?? '',
        photo_url: (state) => state.tg?.initDataUnsafe?.user?.photo_url ?? '',
        get_links: (state) =>
            state.links.map((link) => {
                try {
                    const url = new URL(link)
                    const name = decodeURIComponent(url.hash.replace(/^#/, ''))
                    return `${url.protocol} ${name || url.hostname}`
                } catch {
                    return link
                }
            }),
    },
})
