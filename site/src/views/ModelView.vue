<script setup lang="ts">
import { AutoForm } from "@/components/ui/auto-form"
import { Button } from "@/components/ui/button"
import { DataTable } from "@/components/ui/data-table"
import { Sheet, SheetContent, SheetDescription, SheetFooter, SheetHeader, SheetTitle } from "@/components/ui/sheet"
import { useToast } from "@/components/ui/toast"
import { definitions, type ModelType, type ModelTypes } from "@/models"
import { useAuthSession } from "@/stores/authSession"
import { toTypedSchema } from "@vee-validate/zod"

import axios, { AxiosError } from "axios"
import * as changeCase from "change-case"
import { Eraser, LoaderCircle, Plus, Save } from "lucide-vue-next"
import { v7 as uuidv7 } from "uuid"
import { useForm } from "vee-validate"
import { h } from "vue"
import { computed, ref, watch } from "vue"
import { useRouter } from "vue-router"

const { authSession } = useAuthSession()
const router = useRouter()
const { toast } = useToast()

const definition = computed(() => definitions[router.currentRoute.value.params.modelType.toString()])

const formContext = useForm({
    validationSchema: toTypedSchema(definition.value.schema),
})

const initialValues = ref()
watch(initialValues, (v) => v && formContext.setValues(v))
const isFormOpen = computed({
    get: () => !!initialValues.value,
    set: (v) => (initialValues.value = v ? initialValues.value : undefined),
})

const data = ref<ModelTypes[ModelType][]>([])
const isDataLoading = ref(false)
const refetchData = ref(false)
watch(
    refetchData,
    async () => {
        try {
            isDataLoading.value = true

            const { data: modelData } = await axios<ModelTypes[ModelType][]>({
                method: "GET",
                url: `${import.meta.env.VITE_API_ENDPOINT}/${definition.value.sortType}`,
                headers: { Authorization: authSession?.tokens?.idToken?.toString() },
            })

            data.value = modelData
        } catch (error) {
            console.error(error)

            if (error instanceof AxiosError) {
                toast({
                    title: error.name,
                    description: error.response?.data.message,
                    variant: "destructive",
                })
            }
        } finally {
            isDataLoading.value = false
        }
    },
    { immediate: true }
)

async function onDelete(originalRow: ModelTypes[ModelType]) {
    const url = new URL(
        `${import.meta.env.VITE_API_ENDPOINT}/${definition.value.partitionType}/${originalRow[definition.value.partitionIdKey] || ""}/${definition.value.sortType}`
    )
    if (definition.value.sortIdKey) {
        url.pathname = `/${originalRow[definition.value.sortIdKey]}`
    }

    await axios({
        method: "DELETE",
        url: url.toString(),
        headers: { Authorization: authSession?.tokens?.idToken?.toString() },
    })

    refetchData.value = !refetchData.value
}

async function onDuplicate(originalRow: ModelTypes[ModelType]) {
    originalRow.personId = uuidv7()
    initialValues.value = originalRow
}

async function onEdit(originalRow: ModelTypes[ModelType]) {
    initialValues.value = originalRow
}

async function onRestore(originalRow: ModelTypes[ModelType]) {
    const url = new URL(
        `${import.meta.env.VITE_API_ENDPOINT}/${definition.value.partitionType}/${originalRow[definition.value.partitionIdKey] || ""}/${definition.value.sortType}`
    )
    if (definition.value.sortIdKey) {
        url.pathname = `/${originalRow[definition.value.sortIdKey]}`
    }

    await axios({
        method: "PUT",
        url: url.toString(),
        headers: { Authorization: authSession?.tokens?.idToken?.toString() },
        data: originalRow,
    })

    refetchData.value = !refetchData.value
}

async function onSubmit(values: ModelTypes[ModelType]) {
    const url = new URL(
        `${import.meta.env.VITE_API_ENDPOINT}/${definition.value.partitionType}/${values[definition.value.partitionIdKey] || ""}/${definition.value.sortType}`
    )
    if (definition.value.sortIdKey) {
        url.pathname = `/${values[definition.value.sortIdKey]}`
    }

    await axios({
        method: "PUT",
        url: url.toString(),
        headers: { Authorization: authSession?.tokens?.idToken?.toString() },
        data: values,
    })

    refetchData.value = !refetchData.value
    initialValues.value = undefined
}
</script>

<template>
    <main class="w-full h-full flex flex-col items-start gap-4 p-4">
        <div class="w-full flex justify-end">
            <Button @click="initialValues = definition.getInitialValues()">
                <Plus />
                Create {{ changeCase.noCase(definition.sortType) }}
            </Button>
        </div>

        <DataTable
            :columns="definition.columns"
            v-model:data="data"
            :is-data-loading="isDataLoading"
            @delete="onDelete"
            @duplicate="onDuplicate"
            @edit="onEdit"
            @restore="onRestore"
        />
    </main>

    <Sheet v-model:open="isFormOpen">
        <SheetContent @interact-outside="(event) => event.preventDefault()">
            <SheetHeader>
                <SheetTitle>Edit {{ changeCase.noCase(definition.sortType) }}</SheetTitle>
                <SheetDescription>Make changes here. Click save when you're done.</SheetDescription>
            </SheetHeader>

            <AutoForm
                class="w-full space-y-6 py-4"
                id="form"
                :schema="definition.schema"
                :form="formContext"
                :field-config="{ personId: { component: h('div', { class: 'hidden' }) } }"
                @submit="onSubmit"
            />

            <SheetFooter class="gap-2">
                <Button type="button" variant="secondary" @click="formContext.resetForm()">
                    <Eraser />
                    Clear
                </Button>

                <Button type="submit" form="form" :disabled="formContext.isSubmitting.value">
                    <Save v-if="!formContext.isSubmitting.value" />
                    <LoaderCircle v-else class="animate-spin" />
                    Save
                </Button>
            </SheetFooter>
        </SheetContent>
    </Sheet>
</template>
