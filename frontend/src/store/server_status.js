import axios from 'axios'
import { defineStore } from 'pinia'

export const useServerStatus = defineStore('server', {
    state: () => {
        return {
            serverStatuses: [
                {
                    name: "Service 1",
                    status: "Up 8 days",
                },
                {
                    name: "Service 2",
                    status: "Up 8 days",
                },
            ]
        }
    },
    actions: {
        async loadStatuses() {
            const resp = await axios.get('status')
            this.serverStatuses = resp.data
        }
    }
})