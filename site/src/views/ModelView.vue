<script setup lang="ts">
import { DataTable } from "@/components/ui/data-table"
import { useToast } from "@/components/ui/toast"
import { useAuthSession } from "@/stores/authSession"
import { columns } from "@/models"

import axios, { AxiosError } from "axios"
import { computed, onMounted, ref } from "vue"
import { useRouter } from "vue-router"

const router = useRouter()
const { toast } = useToast()
const { authSession } = useAuthSession()

const data = ref([])
const modelType = computed(() =>
    typeof router.currentRoute.value.params.modelType === "string"
        ? router.currentRoute.value.params.modelType
        : router.currentRoute.value.params.modelType[0]
)

onMounted(async () => {
    try {
        const { data } = await axios({
            method: "GET",
            url: `${import.meta.env.VITE_API_ENDPOINT}/${router.currentRoute.value.params.modelType}`,
            headers: { Authorization: authSession?.tokens?.idToken?.toString() },
        })

        console.log(data)

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
            <DataTable :columns="columns[modelType]" v-model:data="data" />
        </div>
    </main>
</template>
