import { defineStore } from "pinia"

type AuthSessionState = {
    authSession?: unknown
}

export const useAuthSession = defineStore("authSession", {
    state: (): AuthSessionState => ({ authSession: undefined }),
})
