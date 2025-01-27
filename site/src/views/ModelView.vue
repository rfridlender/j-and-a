<script setup lang="ts">
import { DataTable } from "@/components/ui/data-table"
import { useToast } from "@/components/ui/toast"
import { modelDefinitions } from "@/models"
import { useAuthSession } from "@/stores/authSession"

import axios, { AxiosError } from "axios"
import * as changeCase from "change-case"
import { computed, onMounted, ref } from "vue"
import { useRouter } from "vue-router"

const router = useRouter()
const { toast } = useToast()
const { authSession } = useAuthSession()

const modelType = computed(() => router.currentRoute.value.params.modelType.toString())
const data = ref([])
const columns = computed(() => modelDefinitions[modelType.value].columns)

onMounted(async () => {
    try {
        const { data } = await axios({
            method: "GET",
            url: `${import.meta.env.VITE_API_ENDPOINT}/${changeCase.pascalCase(router.currentRoute.value.params.modelType.toString())}`,
            headers: { Authorization: authSession?.tokens?.idToken?.toString() },
        })

        // console.log(data)

        data.value = data
    } catch (error) {
        console.error(error)

        if (error instanceof AxiosError) {
            toast({
                title: error.name,
                description: error.response?.data.message,
                variant: "destructive",
            })
        }
    }
})
</script>

<template>
    <main class="w-full h-full flex flex-col items-start gap-4 p-4">
        <div class="w-full">
            <DataTable :columns="columns" v-model:data="data" />
        </div>
    </main>
</template>
