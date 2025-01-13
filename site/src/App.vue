<script setup lang="ts">
import Toaster from "@/components/ui/toast/Toaster.vue"

import { useColorMode } from "@vueuse/core"
import { Amplify } from "aws-amplify"
import { cognitoUserPoolsTokenProvider } from "aws-amplify/auth/cognito"
import { CookieStorage } from "aws-amplify/utils"
import { RouterView } from "vue-router"

Amplify.configure({
    Auth: {
        Cognito: {
            userPoolId: import.meta.env.VITE_USER_POOL_ID,
            userPoolClientId: import.meta.env.VITE_USER_POOL_CLIENT_ID,
        },
    },
})

cognitoUserPoolsTokenProvider.setKeyValueStorage(new CookieStorage())

useColorMode().value = "auto"
</script>

<template>
    <RouterView />
    <Toaster />
</template>
