<script setup lang="ts" generic="TData">
import { Button } from "@/components/ui/button"
import {
    DropdownMenu,
    DropdownMenuCheckboxItem,
    DropdownMenuContent,
    DropdownMenuLabel,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"

import type { Table } from "@tanstack/vue-table"
import { SlidersHorizontal } from "lucide-vue-next"
import { computed } from "vue"

const props = defineProps<{ table: Table<TData> }>()

const columns = computed(() =>
    props.table.getAllColumns().filter((column) => typeof column.accessorFn !== "undefined" && column.getCanHide())
)
</script>

<template>
    <DropdownMenu>
        <DropdownMenuTrigger as-child>
            <Button variant="outline" size="sm" class="hidden h-8 ml-auto lg:flex">
                <SlidersHorizontal class="w-4 h-4 mr-2" />
                View
            </Button>
        </DropdownMenuTrigger>

        <DropdownMenuContent align="end" class="w-[150px]">
            <DropdownMenuLabel>Toggle columns</DropdownMenuLabel>

            <DropdownMenuSeparator />

            <DropdownMenuCheckboxItem
                v-for="column in columns"
                :key="column.id"
                class="capitalize"
                :checked="column.getIsVisible()"
                @update:checked="(value) => column.toggleVisibility(!!value)"
            >
                {{ column.id }}
            </DropdownMenuCheckboxItem>
        </DropdownMenuContent>
    </DropdownMenu>
</template>
