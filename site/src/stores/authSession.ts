import type { AuthSession } from "aws-amplify/auth"
import { defineStore } from "pinia"

type AuthSessionState = {
    authSession?: AuthSession
}

export const useAuthSession = defineStore("authSession", {
    state: (): AuthSessionState => ({ authSession: undefined }),
})
