import type { FetchUserAttributesOutput } from "aws-amplify/auth"
import { defineStore } from "pinia"

type UserAttributesState = {
    userAttributes?: FetchUserAttributesOutput
}

export const useUserAttributes = defineStore("userAttributes", {
    state: (): UserAttributesState => ({ userAttributes: undefined }),
})
