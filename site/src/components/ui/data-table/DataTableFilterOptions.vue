<script setup lang="ts" generic="TData">
import { Input } from "@/components/ui/input"
import { Select, SelectContent, SelectGroup, SelectLabel, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"

import type { Table } from "@tanstack/vue-table"
import * as changeCase from "change-case"
import { Filter } from "lucide-vue-next"
import { computed } from "vue"

const props = defineProps<{ table: Table<TData> }>()
const filterByColumnId = defineModel<string>({ required: true })

const columns = computed(() => props.table.getAllColumns().filter((column) => column.getCanFilter()))
</script>

<template>
    <div class="grid grid-cols-2">
        <Select v-model="filterByColumnId">
            <SelectTrigger class="rounded-r-none focus:z-10">
                <span class="flex items-center gap-2">
                    <Filter class="w-4 h-4 opacity-50 shrink-0" />
                    <SelectValue />
                </span>
            </SelectTrigger>

            <SelectContent>
                <SelectGroup>
                    <SelectLabel>Columns</SelectLabel>

                    <SelectItem v-for="column in columns" :key="column.id" :value="column.id || ''">
                        {{ changeCase.capitalCase(column.id || "") }}
                    </SelectItem>
                </SelectGroup>
            </SelectContent>
        </Select>

        <Input
            class="max-w-sm rounded-l-none border-l-0"
            :placeholder="`Filter by ${changeCase.noCase(filterByColumnId)}...`"
            :model-value="table.getColumn(filterByColumnId)?.getFilterValue() as string"
            @update:model-value="table.getColumn(filterByColumnId)?.setFilterValue($event)"
        />
    </div>
</template>
